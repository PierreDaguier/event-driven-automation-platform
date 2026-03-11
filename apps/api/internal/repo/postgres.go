package repo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pierre/event-driven-automation-platform/apps/api/internal/models"
)

type PostgresStore struct {
	pool *pgxpool.Pool
}

func NewPostgresStore(ctx context.Context, dbURL string) (*PostgresStore, error) {
	cfg, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, err
	}
	cfg.MaxConns = 10
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return &PostgresStore{pool: pool}, nil
}

func (s *PostgresStore) Close() {
	s.pool.Close()
}

func (s *PostgresStore) Health(ctx context.Context) error {
	return s.pool.Ping(ctx)
}

func (s *PostgresStore) GetWebhookKeyByPublicKey(ctx context.Context, publicKey string) (models.WebhookKey, error) {
	const q = `SELECT id, environment, public_key, secret, created_at FROM webhook_keys WHERE public_key = $1`
	var key models.WebhookKey
	err := s.pool.QueryRow(ctx, q, publicKey).Scan(&key.ID, &key.Environment, &key.PublicKey, &key.Secret, &key.CreatedAt)
	return key, err
}

func (s *PostgresStore) ListWebhookKeys(ctx context.Context) ([]models.WebhookKey, error) {
	const q = `SELECT id, environment, public_key, secret, created_at FROM webhook_keys ORDER BY created_at DESC`
	rows, err := s.pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	keys := []models.WebhookKey{}
	for rows.Next() {
		var k models.WebhookKey
		if err := rows.Scan(&k.ID, &k.Environment, &k.PublicKey, &k.Secret, &k.CreatedAt); err != nil {
			return nil, err
		}
		keys = append(keys, k)
	}
	return keys, rows.Err()
}

func (s *PostgresStore) CreateEvent(ctx context.Context, event models.Event) error {
	const q = `INSERT INTO events (id, source, idempotency_key, event_type, payload, received_at) VALUES ($1,$2,$3,$4,$5,$6)`
	_, err := s.pool.Exec(ctx, q, event.ID, event.Source, event.IdempotencyKey, event.EventType, event.Payload, event.ReceivedAt)
	return err
}

func (s *PostgresStore) ListEnabledWorkflows(ctx context.Context) ([]models.Workflow, error) {
	const q = `SELECT id, name, description, enabled, trigger_event, conditions, actions, webhook_key, created_at, updated_at
	FROM workflows WHERE enabled = true ORDER BY created_at DESC`
	return s.fetchWorkflows(ctx, q)
}

func (s *PostgresStore) ListEnabledWorkflowsByTriggerAndKey(ctx context.Context, trigger string, webhookKey string) ([]models.Workflow, error) {
	const q = `SELECT id, name, description, enabled, trigger_event, conditions, actions, webhook_key, created_at, updated_at
	FROM workflows WHERE enabled = true AND trigger_event = $1 AND webhook_key = $2 ORDER BY created_at DESC`
	return s.fetchWorkflows(ctx, q, trigger, webhookKey)
}

func (s *PostgresStore) ListWorkflows(ctx context.Context) ([]models.Workflow, error) {
	const q = `SELECT id, name, description, enabled, trigger_event, conditions, actions, webhook_key, created_at, updated_at
	FROM workflows ORDER BY created_at DESC`
	return s.fetchWorkflows(ctx, q)
}

func (s *PostgresStore) GetWorkflow(ctx context.Context, workflowID uuid.UUID) (models.Workflow, error) {
	const q = `SELECT id, name, description, enabled, trigger_event, conditions, actions, webhook_key, created_at, updated_at
	FROM workflows WHERE id = $1`
	rows, err := s.fetchWorkflows(ctx, q, workflowID)
	if err != nil {
		return models.Workflow{}, err
	}
	if len(rows) == 0 {
		return models.Workflow{}, errors.New("workflow not found")
	}
	return rows[0], nil
}

func (s *PostgresStore) fetchWorkflows(ctx context.Context, query string, args ...any) ([]models.Workflow, error) {
	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	workflows := []models.Workflow{}
	for rows.Next() {
		var w models.Workflow
		var conditionsRaw, actionsRaw []byte
		if err := rows.Scan(&w.ID, &w.Name, &w.Description, &w.Enabled, &w.Trigger, &conditionsRaw, &actionsRaw, &w.WebhookKey, &w.CreatedAt, &w.UpdatedAt); err != nil {
			return nil, err
		}
		if len(conditionsRaw) > 0 {
			if err := json.Unmarshal(conditionsRaw, &w.Conditions); err != nil {
				return nil, err
			}
		}
		if len(actionsRaw) > 0 {
			if err := json.Unmarshal(actionsRaw, &w.Actions); err != nil {
				return nil, err
			}
		}
		workflows = append(workflows, w)
	}
	return workflows, rows.Err()
}

func (s *PostgresStore) CreateRun(ctx context.Context, run models.Run) error {
	const q = `INSERT INTO runs (id, workflow_id, event_id, action_name, status, retries, latency_ms, error_msg, request_preview, response_preview, created_at, updated_at)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`
	_, err := s.pool.Exec(ctx, q, run.ID, run.WorkflowID, run.EventID, run.ActionName, run.Status, run.Retries, run.LatencyMS, run.Error, run.Request, run.Response, run.CreatedAt, run.UpdatedAt)
	return err
}

func (s *PostgresStore) UpdateRunStatus(ctx context.Context, runID uuid.UUID, status models.RunStatus, retries int) error {
	const q = `UPDATE runs SET status=$2, retries=$3, updated_at=$4 WHERE id=$1`
	_, err := s.pool.Exec(ctx, q, runID, status, retries, time.Now().UTC())
	return err
}

func (s *PostgresStore) FinalizeRun(ctx context.Context, runID uuid.UUID, status models.RunStatus, retries int, latencyMS int, response string, errMsg string, reqPreview []byte) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	const updateRun = `UPDATE runs
	SET status=$2, retries=$3, latency_ms=$4, response_preview=$5, error_msg=$6, request_preview=$7, updated_at=$8
	WHERE id=$1`
	if _, err := tx.Exec(ctx, updateRun, runID, status, retries, latencyMS, response, errMsg, reqPreview, time.Now().UTC()); err != nil {
		return err
	}

	const insertLog = `INSERT INTO execution_logs (run_id, workflow_id, status, latency_ms, request_preview, response_preview, error_msg, created_at)
	SELECT r.id, r.workflow_id, $2, $3, $4, $5, $6, $7 FROM runs r WHERE r.id = $1`
	if _, err := tx.Exec(ctx, insertLog, runID, status, latencyMS, reqPreview, response, errMsg, time.Now().UTC()); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *PostgresStore) GetOverview(ctx context.Context) (models.Overview, error) {
	const q = `SELECT
	COUNT(*) AS total_runs,
	COALESCE(AVG(CASE WHEN status IN ('succeeded','failed') THEN latency_ms END),0) AS avg_latency,
	COALESCE(100.0 * SUM(CASE WHEN status='succeeded' THEN 1 ELSE 0 END) / NULLIF(SUM(CASE WHEN status IN ('succeeded','failed') THEN 1 ELSE 0 END),0),0) AS success_rate,
	SUM(CASE WHEN status='failed' THEN 1 ELSE 0 END) AS failed_runs,
	SUM(CASE WHEN status IN ('pending','running') THEN 1 ELSE 0 END) AS pending_runs
	FROM runs`

	var o models.Overview
	err := s.pool.QueryRow(ctx, q).Scan(&o.TotalRuns, &o.AvgLatencyMS, &o.SuccessRate, &o.FailedRuns, &o.PendingRuns)
	return o, err
}

func (s *PostgresStore) ListExecutionLogs(ctx context.Context, workflowID string, limit int) ([]models.ExecutionLog, error) {
	if limit <= 0 {
		limit = 50
	}
	base := `SELECT id, run_id, workflow_id, status, latency_ms, request_preview, response_preview, error_msg, created_at
	FROM execution_logs`
	args := []any{}
	if workflowID != "" {
		base += ` WHERE workflow_id = $1`
		args = append(args, workflowID)
	}
	base += fmt.Sprintf(" ORDER BY created_at DESC LIMIT %d", limit)

	rows, err := s.pool.Query(ctx, base, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	logs := []models.ExecutionLog{}
	for rows.Next() {
		var l models.ExecutionLog
		if err := rows.Scan(&l.ID, &l.RunID, &l.WorkflowID, &l.Status, &l.LatencyMS, &l.Request, &l.Response, &l.Error, &l.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, rows.Err()
}

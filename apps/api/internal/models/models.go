package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Condition struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
}

type Action struct {
	Name    string            `json:"name"`
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Timeout int               `json:"timeout_ms"`
	Body    map[string]any    `json:"body"`
}

type Workflow struct {
	ID          uuid.UUID   `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Enabled     bool        `json:"enabled"`
	Trigger     string      `json:"trigger"`
	Conditions  []Condition `json:"conditions"`
	Actions     []Action    `json:"actions"`
	WebhookKey  string      `json:"webhook_key"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type Event struct {
	ID             uuid.UUID       `json:"id"`
	Source         string          `json:"source"`
	IdempotencyKey string          `json:"idempotency_key"`
	EventType      string          `json:"event_type"`
	Payload        json.RawMessage `json:"payload"`
	ReceivedAt     time.Time       `json:"received_at"`
}

type RunStatus string

const (
	RunStatusPending   RunStatus = "pending"
	RunStatusRunning   RunStatus = "running"
	RunStatusSucceeded RunStatus = "succeeded"
	RunStatusFailed    RunStatus = "failed"
)

type Run struct {
	ID         uuid.UUID       `json:"id"`
	WorkflowID uuid.UUID       `json:"workflow_id"`
	EventID    uuid.UUID       `json:"event_id"`
	ActionName string          `json:"action_name"`
	Status     RunStatus       `json:"status"`
	Retries    int             `json:"retries"`
	LatencyMS  int             `json:"latency_ms"`
	Error      string          `json:"error,omitempty"`
	Request    json.RawMessage `json:"request_preview"`
	Response   string          `json:"response_preview"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

type ExecutionLog struct {
	ID         int64           `json:"id"`
	RunID      uuid.UUID       `json:"run_id"`
	WorkflowID uuid.UUID       `json:"workflow_id"`
	Status     RunStatus       `json:"status"`
	LatencyMS  int             `json:"latency_ms"`
	Request    json.RawMessage `json:"request_preview"`
	Response   string          `json:"response_preview"`
	Error      string          `json:"error,omitempty"`
	CreatedAt  time.Time       `json:"created_at"`
}

type QueueJob struct {
	RunID      uuid.UUID      `json:"run_id"`
	WorkflowID uuid.UUID      `json:"workflow_id"`
	Action     Action         `json:"action"`
	Payload    map[string]any `json:"payload"`
	Attempt    int            `json:"attempt"`
	MaxRetries int            `json:"max_retries"`
}

type Overview struct {
	TotalRuns    int     `json:"total_runs"`
	SuccessRate  float64 `json:"success_rate"`
	AvgLatencyMS float64 `json:"avg_latency_ms"`
	FailedRuns   int     `json:"failed_runs"`
	PendingRuns  int     `json:"pending_runs"`
}

type WebhookKey struct {
	ID          int       `json:"id"`
	Environment string    `json:"environment"`
	PublicKey   string    `json:"public_key"`
	Secret      string    `json:"secret,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type IngestionResult struct {
	EventID      uuid.UUID `json:"event_id"`
	MatchedRules int       `json:"matched_rules"`
	EnqueuedJobs int       `json:"enqueued_jobs"`
	Duplicate    bool      `json:"duplicate"`
}

package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/pierre/event-driven-automation-platform/apps/api/internal/models"
)

type Store interface {
	GetWebhookKeyByPublicKey(ctx context.Context, publicKey string) (models.WebhookKey, error)
	ListWebhookKeys(ctx context.Context) ([]models.WebhookKey, error)
	CreateEvent(ctx context.Context, event models.Event) error
	ListEnabledWorkflows(ctx context.Context) ([]models.Workflow, error)
	ListEnabledWorkflowsByTriggerAndKey(ctx context.Context, trigger string, webhookKey string) ([]models.Workflow, error)
	CreateRun(ctx context.Context, run models.Run) error
	UpdateRunStatus(ctx context.Context, runID uuid.UUID, status models.RunStatus, retries int) error
	FinalizeRun(ctx context.Context, runID uuid.UUID, status models.RunStatus, retries int, latencyMS int, response string, errMsg string, reqPreview []byte) error
	GetOverview(ctx context.Context) (models.Overview, error)
	ListWorkflows(ctx context.Context) ([]models.Workflow, error)
	GetWorkflow(ctx context.Context, workflowID uuid.UUID) (models.Workflow, error)
	ListExecutionLogs(ctx context.Context, workflowID string, limit int) ([]models.ExecutionLog, error)
}

type QueuePublisher interface {
	Publish(ctx context.Context, job models.QueueJob) error
	PublishDLQ(ctx context.Context, job models.QueueJob, reason string) error
}

type IdempotencyStore interface {
	Reserve(ctx context.Context, key string, ttl time.Duration) (bool, error)
}

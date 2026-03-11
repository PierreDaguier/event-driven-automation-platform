package tests

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pierre/event-driven-automation-platform/apps/api/internal/config"
	"github.com/pierre/event-driven-automation-platform/apps/api/internal/httpapi"
	"github.com/pierre/event-driven-automation-platform/apps/api/internal/models"
	"github.com/pierre/event-driven-automation-platform/apps/api/internal/service"
)

type fakeStore struct {
	mu        sync.Mutex
	webhook   models.WebhookKey
	workflows []models.Workflow
	runs      []models.Run
}

func (f *fakeStore) GetWebhookKeyByPublicKey(ctx context.Context, publicKey string) (models.WebhookKey, error) {
	return f.webhook, nil
}
func (f *fakeStore) ListWebhookKeys(ctx context.Context) ([]models.WebhookKey, error) {
	return []models.WebhookKey{f.webhook}, nil
}
func (f *fakeStore) CreateEvent(ctx context.Context, event models.Event) error { return nil }
func (f *fakeStore) ListEnabledWorkflows(ctx context.Context) ([]models.Workflow, error) {
	return f.workflows, nil
}
func (f *fakeStore) ListEnabledWorkflowsByTriggerAndKey(ctx context.Context, trigger string, webhookKey string) ([]models.Workflow, error) {
	return f.workflows, nil
}
func (f *fakeStore) CreateRun(ctx context.Context, run models.Run) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.runs = append(f.runs, run)
	return nil
}
func (f *fakeStore) UpdateRunStatus(ctx context.Context, runID uuid.UUID, status models.RunStatus, retries int) error {
	return nil
}
func (f *fakeStore) FinalizeRun(ctx context.Context, runID uuid.UUID, status models.RunStatus, retries int, latencyMS int, response string, errMsg string, reqPreview []byte) error {
	return nil
}
func (f *fakeStore) GetOverview(ctx context.Context) (models.Overview, error) {
	return models.Overview{}, nil
}
func (f *fakeStore) ListWorkflows(ctx context.Context) ([]models.Workflow, error) {
	return f.workflows, nil
}
func (f *fakeStore) GetWorkflow(ctx context.Context, workflowID uuid.UUID) (models.Workflow, error) {
	return f.workflows[0], nil
}
func (f *fakeStore) ListExecutionLogs(ctx context.Context, workflowID string, limit int) ([]models.ExecutionLog, error) {
	return nil, nil
}

type fakeQueue struct {
	jobs []models.QueueJob
}

func (f *fakeQueue) Publish(ctx context.Context, job models.QueueJob) error {
	f.jobs = append(f.jobs, job)
	return nil
}
func (f *fakeQueue) PublishDLQ(ctx context.Context, job models.QueueJob, reason string) error {
	return nil
}

type memoryIdempotency struct {
	set map[string]struct{}
}

func (m *memoryIdempotency) Reserve(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	if _, ok := m.set[key]; ok {
		return false, nil
	}
	m.set[key] = struct{}{}
	return true, nil
}

func TestWebhookIngestion(t *testing.T) {
	store := &fakeStore{
		webhook: models.WebhookKey{PublicKey: "whk_test", Secret: "whsec_test"},
		workflows: []models.Workflow{
			{
				ID:      uuid.New(),
				Enabled: true,
				Trigger: "lead.created",
				Conditions: []models.Condition{
					{Field: "lead.score", Operator: "gt", Value: 70},
				},
				Actions: []models.Action{{Name: "push", Method: "POST", URL: "https://example.com"}},
			},
		},
	}
	queue := &fakeQueue{}

	ingestor := service.IngestionService{
		Store:          store,
		Queue:          queue,
		Idempotency:    &memoryIdempotency{set: map[string]struct{}{}},
		MaxRetries:     3,
		IdempotencyTTL: time.Hour,
	}

	router := httpapi.NewRouter(httpapi.API{
		Cfg:      config.Config{WebhookHeader: "X-Signature", AllowedOrigins: "*"},
		Store:    store,
		Ingestor: ingestor,
	})

	payload := map[string]any{"event_type": "lead.created", "lead": map[string]any{"score": 91}}
	body, _ := json.Marshal(payload)

	mac := hmac.New(sha256.New, []byte("whsec_test"))
	_, _ = mac.Write(body)
	sig := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	req := httptest.NewRequest(http.MethodPost, "/api/v1/webhooks/whk_test", bytes.NewReader(body))
	req.Header.Set("X-Signature", sig)
	req.Header.Set("Idempotency-Key", "idem-123")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	if w.Code != http.StatusAccepted {
		t.Fatalf("expected 202, got %d", w.Code)
	}

	if len(queue.jobs) != 1 {
		t.Fatalf("expected 1 enqueued job, got %d", len(queue.jobs))
	}
}

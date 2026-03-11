package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pierre/event-driven-automation-platform/apps/api/internal/models"
	"github.com/pierre/event-driven-automation-platform/apps/api/internal/rules"
)

var ErrDuplicateEvent = errors.New("duplicate event")

type IngestionService struct {
	Store          Store
	Queue          QueuePublisher
	Idempotency    IdempotencyStore
	MaxRetries     int
	IdempotencyTTL time.Duration
}

func (s IngestionService) Ingest(ctx context.Context, source string, webhookKey string, idempotencyKey string, payloadBytes []byte) (models.IngestionResult, error) {
	result := models.IngestionResult{}
	if idempotencyKey == "" {
		return result, errors.New("missing idempotency key")
	}

	reserved, err := s.Idempotency.Reserve(ctx, fmt.Sprintf("idempotency:%s", idempotencyKey), s.IdempotencyTTL)
	if err != nil {
		return result, err
	}
	if !reserved {
		result.Duplicate = true
		return result, ErrDuplicateEvent
	}

	var payload map[string]any
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return result, errors.New("invalid json payload")
	}

	eventType := "incoming"
	if raw, ok := payload["event_type"]; ok {
		eventType = fmt.Sprintf("%v", raw)
	}

	event := models.Event{
		ID:             uuid.New(),
		Source:         source,
		IdempotencyKey: idempotencyKey,
		EventType:      eventType,
		Payload:        payloadBytes,
		ReceivedAt:     time.Now().UTC(),
	}
	if err := s.Store.CreateEvent(ctx, event); err != nil {
		return result, err
	}

	workflows, err := s.Store.ListEnabledWorkflowsByTriggerAndKey(ctx, eventType, webhookKey)
	if err != nil {
		return result, err
	}

	result.EventID = event.ID
	for _, workflow := range workflows {
		if !rules.Matches(workflow.Conditions, payload) {
			continue
		}
		result.MatchedRules++

		for _, action := range workflow.Actions {
			run := models.Run{
				ID:         uuid.New(),
				WorkflowID: workflow.ID,
				EventID:    event.ID,
				ActionName: action.Name,
				Status:     models.RunStatusPending,
				Retries:    0,
				CreatedAt:  time.Now().UTC(),
				UpdatedAt:  time.Now().UTC(),
			}
			if err := s.Store.CreateRun(ctx, run); err != nil {
				return result, err
			}
			job := models.QueueJob{
				RunID:      run.ID,
				WorkflowID: workflow.ID,
				Action:     action,
				Payload:    payload,
				Attempt:    0,
				MaxRetries: s.MaxRetries,
			}
			if err := s.Queue.Publish(ctx, job); err != nil {
				return result, err
			}
			result.EnqueuedJobs++
		}
	}

	return result, nil
}

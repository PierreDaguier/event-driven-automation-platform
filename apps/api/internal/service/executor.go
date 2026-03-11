package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/pierre/event-driven-automation-platform/apps/api/internal/models"
	"github.com/pierre/event-driven-automation-platform/apps/api/internal/utils"
)

type ActionExecutor struct {
	Store           Store
	Queue           QueuePublisher
	HTTPClient      *http.Client
	OutgoingTimeout time.Duration
}

func (e ActionExecutor) ProcessJob(ctx context.Context, job models.QueueJob) error {
	if err := e.Store.UpdateRunStatus(ctx, job.RunID, models.RunStatusRunning, job.Attempt); err != nil {
		return err
	}

	bodyObj := map[string]any{
		"payload": job.Payload,
		"action":  job.Action.Body,
	}

	redacted := utils.RedactPayload(bodyObj)
	reqBytes, _ := json.Marshal(redacted)

	timeout := e.OutgoingTimeout
	if job.Action.Timeout > 0 {
		timeout = time.Duration(job.Action.Timeout) * time.Millisecond
	}
	ctxReq, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	requestBody, _ := json.Marshal(bodyObj)
	req, err := http.NewRequestWithContext(ctxReq, strings.ToUpper(job.Action.Method), job.Action.URL, bytes.NewReader(requestBody))
	if err != nil {
		return e.finalizeWithRetry(ctx, job, 0, "", err, reqBytes)
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range job.Action.Headers {
		req.Header.Set(k, v)
	}

	start := time.Now()
	resp, err := e.HTTPClient.Do(req)
	if err != nil {
		return e.finalizeWithRetry(ctx, job, int(time.Since(start).Milliseconds()), "", err, reqBytes)
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
	latency := int(time.Since(start).Milliseconds())

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		if err := e.Store.FinalizeRun(ctx, job.RunID, models.RunStatusSucceeded, job.Attempt, latency, string(respBytes), "", reqBytes); err != nil {
			return err
		}
		return nil
	}

	err = errors.New(fmt.Sprintf("downstream status %d", resp.StatusCode))
	return e.finalizeWithRetry(ctx, job, latency, string(respBytes), err, reqBytes)
}

func (e ActionExecutor) finalizeWithRetry(ctx context.Context, job models.QueueJob, latency int, response string, runErr error, request []byte) error {
	if job.Attempt < job.MaxRetries {
		next := job
		next.Attempt = job.Attempt + 1
		backoff := time.Duration(1<<job.Attempt) * time.Second
		time.Sleep(backoff)
		if err := e.Store.UpdateRunStatus(ctx, job.RunID, models.RunStatusPending, next.Attempt); err != nil {
			return err
		}
		if err := e.Queue.Publish(ctx, next); err != nil {
			return err
		}
		return nil
	}

	if err := e.Store.FinalizeRun(ctx, job.RunID, models.RunStatusFailed, job.Attempt, latency, response, runErr.Error(), request); err != nil {
		return err
	}
	if err := e.Queue.PublishDLQ(ctx, job, runErr.Error()); err != nil {
		return err
	}
	return nil
}

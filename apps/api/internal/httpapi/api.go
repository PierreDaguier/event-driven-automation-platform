package httpapi

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/pierre/event-driven-automation-platform/apps/api/internal/config"
	custommiddleware "github.com/pierre/event-driven-automation-platform/apps/api/internal/middleware"
	"github.com/pierre/event-driven-automation-platform/apps/api/internal/service"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type API struct {
	Cfg      config.Config
	Store    service.Store
	Ingestor service.IngestionService
}

var (
	webhookCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "webhook_ingestions_total",
		Help: "Total webhook ingestions by status",
	}, []string{"status"})

	queueJobsCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "queue_jobs_enqueued_total",
		Help: "Total number of jobs published to queue",
	})
)

func init() {
	prometheus.MustRegister(webhookCounter, queueJobsCounter)
}

func NewRouter(api API) http.Handler {
	r := chi.NewRouter()
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(30 * time.Second))
	r.Use(custommiddleware.CORS(api.Cfg.AllowedOrigins))

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	r.Handle("/metrics", promhttp.Handler())

	r.Route("/api/v1", func(v1 chi.Router) {
		v1.Post("/webhooks/{publicKey}", api.handleWebhook)
		v1.Get("/overview", api.handleOverview)
		v1.Get("/workflows", api.handleWorkflows)
		v1.Get("/workflows/{workflowID}", api.handleWorkflowDetail)
		v1.Get("/logs", api.handleLogs)
		v1.Get("/settings", api.handleSettings)
	})

	return otelhttp.NewHandler(r, "api-router")
}

func (api API) handleWebhook(w http.ResponseWriter, r *http.Request) {
	publicKey := chi.URLParam(r, "publicKey")
	if publicKey == "" {
		writeError(w, http.StatusBadRequest, "missing webhook public key")
		return
	}

	webhookKey, err := api.Store.GetWebhookKeyByPublicKey(r.Context(), publicKey)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid webhook key")
		return
	}

	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		writeError(w, http.StatusBadRequest, "unable to read payload")
		return
	}

	sig := r.Header.Get(api.Cfg.WebhookHeader)
	if err := custommiddleware.VerifyHMACSHA256(sig, body, webhookKey.Secret); err != nil {
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}

	idempotencyKey := r.Header.Get("Idempotency-Key")
	result, err := api.Ingestor.Ingest(r.Context(), "webhook", publicKey, idempotencyKey, body)
	if err != nil {
		if errors.Is(err, service.ErrDuplicateEvent) {
			webhookCounter.WithLabelValues("duplicate").Inc()
			writeJSON(w, http.StatusAccepted, result)
			return
		}
		webhookCounter.WithLabelValues("failed").Inc()
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if result.EnqueuedJobs > 0 {
		queueJobsCounter.Add(float64(result.EnqueuedJobs))
	}
	webhookCounter.WithLabelValues("accepted").Inc()
	writeJSON(w, http.StatusAccepted, result)
}

func (api API) handleOverview(w http.ResponseWriter, r *http.Request) {
	overview, err := api.Store.GetOverview(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, overview)
}

func (api API) handleWorkflows(w http.ResponseWriter, r *http.Request) {
	workflows, err := api.Store.ListWorkflows(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, workflows)
}

func (api API) handleWorkflowDetail(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "workflowID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid workflow id")
		return
	}
	workflow, err := api.Store.GetWorkflow(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, workflow)
}

func (api API) handleLogs(w http.ResponseWriter, r *http.Request) {
	limit := 50
	if raw := r.URL.Query().Get("limit"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err == nil {
			limit = parsed
		}
	}
	logs, err := api.Store.ListExecutionLogs(r.Context(), r.URL.Query().Get("workflow_id"), limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, logs)
}

func (api API) handleSettings(w http.ResponseWriter, r *http.Request) {
	keys, err := api.Store.ListWebhookKeys(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	for i := range keys {
		if len(keys[i].Secret) > 8 {
			keys[i].Secret = keys[i].Secret[:4] + strings.Repeat("*", len(keys[i].Secret)-8) + keys[i].Secret[len(keys[i].Secret)-4:]
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"webhook_keys": keys,
		"environment":  api.Cfg.AppEnv,
	})
}

func writeError(w http.ResponseWriter, code int, message string) {
	writeJSON(w, code, map[string]string{"error": message})
}

func writeJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(payload)
}

package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/pierre/event-driven-automation-platform/apps/api/internal/config"
	"github.com/pierre/event-driven-automation-platform/apps/api/internal/httpapi"
	"github.com/pierre/event-driven-automation-platform/apps/api/internal/queue"
	"github.com/pierre/event-driven-automation-platform/apps/api/internal/repo"
	"github.com/pierre/event-driven-automation-platform/apps/api/internal/service"
	"github.com/pierre/event-driven-automation-platform/apps/api/internal/telemetry"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()

	shutdownTracer, err := telemetry.InitTracer(ctx, "automation-api", cfg.OTELExporterOTLP)
	if err != nil {
		log.Fatalf("failed to initialize telemetry: %v", err)
	}
	defer func() {
		_ = shutdownTracer(context.Background())
	}()

	store, err := repo.NewPostgresStore(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect postgres: %v", err)
	}
	defer store.Close()

	if err := store.Health(ctx); err != nil {
		log.Fatalf("postgres ping failed: %v", err)
	}

	idempotency := repo.NewRedisIdempotency(cfg.RedisAddr, cfg.RedisPassword)
	defer idempotency.Close()

	rabbit, err := queue.NewRabbitMQ(cfg.RabbitMQURL, cfg.RabbitMQQueue, cfg.RabbitMQDLQ)
	if err != nil {
		log.Fatalf("failed to connect rabbitmq: %v", err)
	}
	defer rabbit.Close()

	ingestor := service.IngestionService{
		Store:          store,
		Queue:          rabbit,
		Idempotency:    idempotency,
		MaxRetries:     cfg.MaxRetries,
		IdempotencyTTL: time.Duration(cfg.IdempotencyTTLHour) * time.Hour,
	}

	router := httpapi.NewRouter(httpapi.API{
		Cfg:      cfg,
		Store:    store,
		Ingestor: ingestor,
	})

	srv := &http.Server{
		Addr:         cfg.HTTPAddress(),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdownCtx)
	}()

	log.Printf("automation api listening on %s", cfg.HTTPAddress())
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server failed: %v", err)
	}
}

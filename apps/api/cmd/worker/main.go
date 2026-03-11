package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/pierre/event-driven-automation-platform/apps/api/internal/config"
	"github.com/pierre/event-driven-automation-platform/apps/api/internal/queue"
	"github.com/pierre/event-driven-automation-platform/apps/api/internal/repo"
	"github.com/pierre/event-driven-automation-platform/apps/api/internal/service"
	"github.com/pierre/event-driven-automation-platform/apps/api/internal/telemetry"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()

	shutdownTracer, err := telemetry.InitTracer(ctx, "automation-worker", cfg.OTELExporterOTLP)
	if err != nil {
		log.Fatalf("failed to init telemetry: %v", err)
	}
	defer func() {
		_ = shutdownTracer(context.Background())
	}()

	store, err := repo.NewPostgresStore(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect postgres: %v", err)
	}
	defer store.Close()

	rabbit, err := queue.NewRabbitMQ(cfg.RabbitMQURL, cfg.RabbitMQQueue, cfg.RabbitMQDLQ)
	if err != nil {
		log.Fatalf("failed to connect rabbitmq: %v", err)
	}
	defer rabbit.Close()

	executor := service.ActionExecutor{
		Store:           store,
		Queue:           rabbit,
		HTTPClient:      &http.Client{Timeout: time.Duration(cfg.OutgoingTimeoutMS) * time.Millisecond},
		OutgoingTimeout: time.Duration(cfg.OutgoingTimeoutMS) * time.Millisecond,
	}

	log.Println("worker consuming queue")
	if err := rabbit.Consume(ctx, "automation-worker", executor.ProcessJob); err != nil {
		log.Fatalf("worker consume failure: %v", err)
	}
}

package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	AppEnv             string
	HTTPPort           string
	DatabaseURL        string
	RedisAddr          string
	RedisPassword      string
	RabbitMQURL        string
	RabbitMQQueue      string
	RabbitMQDLQ        string
	OutgoingTimeoutMS  int
	MaxRetries         int
	OTELExporterOTLP   string
	AllowedOrigins     string
	WebhookHeader      string
	IdempotencyTTLHour int
}

func Load() Config {
	cfg := Config{
		AppEnv:             getEnv("APP_ENV", "development"),
		HTTPPort:           getEnv("HTTP_PORT", "8080"),
		DatabaseURL:        getEnv("DATABASE_URL", "postgres://postgres:postgres@postgres:5432/automation?sslmode=disable"),
		RedisAddr:          getEnv("REDIS_ADDR", "redis:6379"),
		RedisPassword:      getEnv("REDIS_PASSWORD", ""),
		RabbitMQURL:        getEnv("RABBITMQ_URL", "amqp://guest:guest@rabbitmq:5672/"),
		RabbitMQQueue:      getEnv("RABBITMQ_QUEUE", "workflow-actions"),
		RabbitMQDLQ:        getEnv("RABBITMQ_DLQ", "workflow-actions-dlq"),
		OutgoingTimeoutMS:  getEnvInt("OUTGOING_TIMEOUT_MS", 8000),
		MaxRetries:         getEnvInt("MAX_RETRIES", 3),
		OTELExporterOTLP:   getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "otel-collector:4317"),
		AllowedOrigins:     getEnv("ALLOWED_ORIGINS", "http://localhost:3000"),
		WebhookHeader:      getEnv("WEBHOOK_SIGNATURE_HEADER", "X-Signature"),
		IdempotencyTTLHour: getEnvInt("IDEMPOTENCY_TTL_HOURS", 24),
	}
	return cfg
}

func (c Config) HTTPAddress() string {
	return fmt.Sprintf(":%s", c.HTTPPort)
}

func getEnv(k, fallback string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(k string, fallback int) int {
	v := os.Getenv(k)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}

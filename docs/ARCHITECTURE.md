# Architecture

```mermaid
flowchart LR
  Client[Client systems\nCRM / Billing / Forms] -->|Signed webhook\nHMAC + Idempotency-Key| API[Go API\nWebhook Ingestion]

  API --> PG[(PostgreSQL 18)]
  API --> Redis[(Redis 8)]
  API --> MQ[(RabbitMQ 4.2 queue)]

  MQ --> Worker[Go Worker\nRetry + Backoff]
  Worker --> External[External HTTP Actions\nCRM / Slack / Billing]
  Worker --> PG
  Worker --> DLQ[(RabbitMQ DLQ)]

  API --> FE[Next.js 16 Dashboard]
  FE --> ClientUser[Business users]

  API --> OTel[OpenTelemetry Collector]
  Worker --> OTel
  OTel --> Prom[Prometheus]
  Prom --> Grafana[Grafana dashboards]
```

## Notes

- **Security**: HMAC SHA-256 signature verification + Redis idempotency reservation before persistence.
- **Reliability**: retries are exponential (1s, 2s, 4s by default) and failed jobs are routed to DLQ.
- **Auditability**: execution logs persist status, latency, redacted payload preview, and response/error snippets.
- **UX**: dashboard targets non-technical stakeholders with KPI-first information architecture.

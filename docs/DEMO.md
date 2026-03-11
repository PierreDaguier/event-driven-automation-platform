# How To Demo To A Client In 5 Minutes

## 1) Start the platform (45s)

```bash
./scripts/dev-up.sh
```

Open:
- Web app: http://localhost:3000
- RabbitMQ: http://localhost:15672
- Prometheus: http://localhost:9090
- Grafana: http://localhost:3001

## 2) Show value immediately (45s)

- Open landing page (`/`) and explain business outcomes:
  - secure ingestion,
  - reliable processing,
  - measurable SLA visibility.

## 3) Open dashboard KPIs (60s)

- Go to `/dashboard`.
- Highlight:
  - total runs,
  - success rate,
  - average latency,
  - failed runs to investigate.

## 4) Trigger live event (60s)

```bash
./scripts/demo-webhook.sh
```

- Explain HMAC signature and idempotency key used by the script.

## 5) Show traceability (90s)

- `/dashboard/logs`: timeline updates with status + latency + redacted payload preview.
- `/dashboard/workflows`: show rule logic and outbound actions.
- Mention retries + DLQ behavior for failed downstream endpoints.

## 6) Show production credibility (40s)

- `openapi/openapi.yaml`
- `docs/ARCHITECTURE.md` (diagram)
- CI workflows in `.github/workflows/`

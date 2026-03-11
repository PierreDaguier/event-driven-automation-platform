# Minimum PR Plan (as requested)

## PR #1 - scaffold architecture + docker + ci
- Branch: `chore/observability-stack`
- Scope: monorepo structure, compose stack, CI workflows, base docs.

## PR #2 - webhook ingestion + security
- Branch: `fix/webhook-signature-validation`
- Scope: HMAC validation, idempotency handling, webhook endpoint.

## PR #3 - queue + retry + dlq
- Branch: `feat/workflow-engine-core`
- Scope: RabbitMQ queueing, retry backoff, DLQ routing, worker runtime.

## PR #4 - dashboard UX + logs timeline
- Branch: `feat/ui-dashboard-runs`
- Scope: landing, overview metrics, workflow screens, logs/settings timeline.

## PR #5 - observability + docs + demo script
- Branch: `docs/client-demo-script`
- Scope: OTel + Prometheus + Grafana configs, OpenAPI, README, screenshots, 5-min demo flow.

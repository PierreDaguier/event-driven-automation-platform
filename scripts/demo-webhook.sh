#!/usr/bin/env bash
set -euo pipefail

PAYLOAD='{"event_type":"lead.created","lead":{"score":88,"region":"us","email":"client@example.com"}}'
SECRET='whsec_staging_demo_secret_01'
IDEMPOTENCY_KEY="demo-$(date +%s)"

SIGNATURE="sha256=$(printf '%s' "$PAYLOAD" | openssl dgst -sha256 -hmac "$SECRET" -binary | xxd -p -c 256)"

RESPONSE="$(curl -sS -X POST "http://localhost:8080/api/v1/webhooks/whk_staging_demo_01" \
  -H "Content-Type: application/json" \
  -H "X-Signature: ${SIGNATURE}" \
  -H "Idempotency-Key: ${IDEMPOTENCY_KEY}" \
  -d "$PAYLOAD")"

if command -v jq >/dev/null 2>&1; then
  printf '%s\n' "$RESPONSE" | jq .
else
  printf '%s\n' "$RESPONSE"
fi

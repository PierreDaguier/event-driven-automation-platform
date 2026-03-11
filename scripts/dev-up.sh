#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."
docker compose up -d --build
echo "Web: http://localhost:3000"
echo "API: http://localhost:8080"
echo "RabbitMQ: http://localhost:15672"
echo "Prometheus: http://localhost:9090"
echo "Grafana: http://localhost:3001 (admin/admin)"

.PHONY: install up down logs test seed demo lint

install:
	cd apps/web && npm ci
	cd apps/api && docker run --rm -v "$$(pwd)":/src -w /src golang:1.26-alpine go mod tidy

up:
	docker compose up -d --build

down:
	docker compose down -v

logs:
	docker compose logs -f --tail=200

test:
	docker compose run --rm api sh -lc "go test ./..."

demo:
	./scripts/demo-webhook.sh

lint:
	docker compose run --rm web sh -lc "npm run lint"

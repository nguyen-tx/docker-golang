.PHONY: help dev up down logs migrate-up migrate-down tidy

help: ## Hiển thị danh sách lệnh
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# ─── Docker ───────────────────────────────────────────────────────────────────

up: ## Khởi động toàn bộ dự án (Docker)
	docker compose up -d

down: ## Dừng toàn bộ dự án
	docker compose down

logs: ## Xem logs
	docker compose logs -f

db-up: ## Chỉ khởi động databases
	docker compose up -d postgres mongo

# ─── Development ──────────────────────────────────────────────────────────────

dev-backend: ## Chạy backend với hot-reload (Air)
	cd backend && air

dev-frontend: ## Chạy frontend dev server
	cd frontend && npm run dev

# ─── Migration ────────────────────────────────────────────────────────────────

MIGRATE_URL ?= postgres://utm_user:utm_pass@localhost:5432/utm_db?sslmode=disable
MIGRATE_DIR = backend/migrations/postgres

migrate-up: ## Chạy migration (up)
	migrate -path $(MIGRATE_DIR) -database "$(MIGRATE_URL)" up

migrate-down: ## Rollback migration
	migrate -path $(MIGRATE_DIR) -database "$(MIGRATE_URL)" down 1

migrate-status: ## Xem trạng thái migration
	migrate -path $(MIGRATE_DIR) -database "$(MIGRATE_URL)" version

# ─── Go ───────────────────────────────────────────────────────────────────────

tidy: ## go mod tidy
	cd backend && go mod tidy

build: ## Build backend binary
	cd backend && go build -o bin/utm-server ./cmd/server

test: ## Chạy tests
	cd backend && go test ./... -v

# ─── Setup ────────────────────────────────────────────────────────────────────

setup: ## Khởi tạo lần đầu
	cp -n .env.example .env || true
	cd frontend && npm install
	cd backend && go mod download
	@echo "✓ Setup done. Chỉnh sửa .env rồi chạy: make db-up && make migrate-up && make dev-backend"

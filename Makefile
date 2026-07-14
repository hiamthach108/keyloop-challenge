BACKEND_DIR := backend
FRONTEND_DIR := frontend
COMPOSE_FILE := docker-compose.yml
name ?= schema_change

setup:
	@test -f $(BACKEND_DIR)/.env || cp $(BACKEND_DIR)/.env.example $(BACKEND_DIR)/.env
	@test -f $(FRONTEND_DIR)/.env.local || cp $(FRONTEND_DIR)/.env.example $(FRONTEND_DIR)/.env.local
	npm --prefix $(FRONTEND_DIR) install

infra-up:
	docker-compose -f $(COMPOSE_FILE) up -d --wait postgres

infra-down:
	docker-compose -f $(COMPOSE_FILE) down

backend-run:
	$(MAKE) -C $(BACKEND_DIR) run

backend-test:
	$(MAKE) -C $(BACKEND_DIR) test

backend-lint:
	$(MAKE) -C $(BACKEND_DIR) lint

backend-tidy:
	cd $(BACKEND_DIR) && go mod tidy

migration-diff:
	cd $(BACKEND_DIR) && atlas migrate diff $(name) --env gorm

migration-hash:
	cd $(BACKEND_DIR) && atlas migrate hash --dir file://migrations

migration-validate:
	cd $(BACKEND_DIR) && atlas migrate validate --dir file://migrations

migration-apply:
	cd $(BACKEND_DIR) && atlas migrate apply --dir file://migrations --url 'postgres://postgres:postgres@localhost:5432/keyloop_inventory?search_path=public&sslmode=disable'

frontend-run:
	npm --prefix $(FRONTEND_DIR) run dev

frontend-lint:
	npm --prefix $(FRONTEND_DIR) run lint

frontend-typecheck:
	npm --prefix $(FRONTEND_DIR) run typecheck

frontend-build:
	npm --prefix $(FRONTEND_DIR) run build

docker-up:
	docker-compose -f $(COMPOSE_FILE) up --build

docker-up-detached:
	docker-compose -f $(COMPOSE_FILE) up --build -d

docker-down:
	docker-compose -f $(COMPOSE_FILE) down

docker-reset:
	docker-compose -f $(COMPOSE_FILE) down -v
	docker-compose -f $(COMPOSE_FILE) up --build

dev: setup infra-up migration-apply
	$(MAKE) -j2 backend-run frontend-run

verify: backend-test migration-validate frontend-lint frontend-typecheck frontend-build

.PHONY: setup infra-up infra-down backend-run backend-test backend-lint backend-tidy migration-diff migration-hash migration-validate migration-apply frontend-run frontend-lint frontend-typecheck frontend-build docker-up docker-up-detached docker-down docker-reset dev verify

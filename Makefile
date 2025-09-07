.PHONY: dev prod down logs dev-logs api dbshell redis-shell \
	migrate-up migrate-down migrate-force migrate-create migrate-status reset-db

# ---------- Load environment variables ----------
ifneq (,$(wildcard .env))
    include .env
    export $(shell sed 's/=.*//' .env)
endif

# ---------- Compose profiles ----------
dev:
	@docker compose up -d --build --remove-orphans

prod:
	@docker compose up -d --build --remove-orphans

down:
	@docker compose down --remove-orphans

logs:
	@docker compose logs -f

dev-logs:
	@docker compose logs -f api-dev

api:
	@docker compose build api && docker compose up -d api

# ---------- Shell helpers ----------
dbshell:
	@docker compose exec -e PGPASSWORD=$(POSTGRES_PASSWORD) postgres \
		psql -U $(POSTGRES_USER) -d $(POSTGRES_DB)

redis-shell:
	@docker compose exec redis redis-cli

# ---------- Migrations via 'migrator' service ----------
migrate-up:
	@docker compose run --rm -e DATABASE_URL=$(DATABASE_URL) migrator migrate \
		-path=/migrations \
		-database "$(DATABASE_URL)" \
		-verbose up

n ?= 1
migrate-down:
	@docker compose run --rm -e DATABASE_URL=$(DATABASE_URL) migrator migrate \
		-path=/migrations \
		-database "$(DATABASE_URL)" \
		-verbose down $(n)

migrate-force:
	@docker compose run --rm -e DATABASE_URL=$(DATABASE_URL) migrator migrate \
		-path=/migrations \
		-database "$(DATABASE_URL)" \
		-verbose force $(v)

migrate-status:
	@docker compose run --rm -e DATABASE_URL=$(DATABASE_URL) migrator migrate \
		-path=/migrations \
		-database "$(DATABASE_URL)" \
		-verbose version

migrate-create:
	@docker run --rm -v $(PWD)/infra/migrations:/migrations migrate/migrate:v4.17.0 \
		create -ext sql -dir /migrations -seq $(name)

reset-db:
	@docker compose down -v --remove-orphans
	@docker compose up -d postgres redis
	@sleep 5
	@make migrate-up

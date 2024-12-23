SHELL := /bin/bash

DB_CONFIG := config.json
MIGRATION_DIR := migration

test:
 go test ./test/

migrate-up:
	@echo "==> Migrating up..."
	go run cmd/migrate/main.go -config=$(DB_CONFIG) -dir=$(MIGRATION_DIR) up

migrate-down:
	@echo "==> Migrating down..."
	go run cmd/migrate/main.go -config=$(DB_CONFIG) -dir=$(MIGRATION_DIR) down

migrate-drop:
	@echo "==> Dropping all tables..."
	go run cmd/migrate/main.go -config=$(DB_CONFIG) -dir=$(MIGRATION_DIR) drop

migrate-version:
	@echo "==> Current migration version..."
	go run cmd/migrate/main.go -config=$(DB_CONFIG) -dir=$(MIGRATION_DIR) version

run:
	@echo "==> Running main app ..."
	go run cmd/main.go -config=$(DB_CONFIG)

help:
	@echo "Usage: make [target]"
	@echo "Targets:"
	@echo "  migrate-up      : Apply all available migrations"
	@echo "  migrate-down    : Rollback the last migration (or all, step by step)"
	@echo "  migrate-drop    : Drop all tables (careful!)"
	@echo "  migrate-version : Show current migration version"
	@echo "  run             : Run the main application"

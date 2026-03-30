.PHONY: migrate-create migrate-up-local migrate-down-local migrate-up migrate-down migrate-force migrate-version migrate-clean migrate-reset

# Путь к миграциям
MIGRATIONS_PATH = ./internal/migrations

# Загрузка .env файла
ifneq (,$(wildcard .env))
    include .env
    export
endif

# Переменные для подключения к БД
DB_USER ?= postgres
DB_PASSWORD ?= 1234
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_NAME ?= soar
DB_SSLMODE ?= disable

# Локальная строка подключения
DB_URL_LOCAL = postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

# Создание миграции
migrate-create:
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(name)

# Применить все миграции локально
migrate-up-local:
	@echo "🔄 Applying migrations locally..."
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL_LOCAL)" up
	@echo "✅ Migrations applied successfully!"

# Откатить одну миграцию локально
migrate-down-local:
	@echo "🔄 Rolling back last migration..."
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL_LOCAL)" down 1
	@echo "✅ Migration rolled back!"

# Откатить все миграции локально
migrate-down-local-all:
	@echo "⚠️  WARNING: This will rollback ALL migrations!"
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL_LOCAL)" down -all
	@echo "✅ All migrations rolled back!"

# Применить все миграции (для production)
migrate-up:
	@echo "🔄 Applying migrations to production..."
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" up
	@echo "✅ Migrations applied successfully!"

# Откатить N миграций
migrate-down:
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" down $(n)

# Принудительно установить версию миграции (ФИКС DIRTY)
migrate-force:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL_LOCAL)" force $(version)

# Показать текущую версию
migrate-version:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL_LOCAL)" version

# Очистить грязное состояние
migrate-clean:
	@echo "🧹 Cleaning dirty state..."
	psql "$(DB_URL_LOCAL)" -c "DROP TABLE IF EXISTS schema_migrations;"
	@echo "✅ Dirty state cleaned!"

# Полный сброс (ОСТОРОЖНО!)
migrate-reset:
	@echo "⚠️  WARNING: This will delete ALL data!"
	psql "$(DB_URL_LOCAL)" -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
	@echo "✅ Database reset!"
	@echo "🔄 Reapplying migrations..."
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL_LOCAL)" up
	@echo "✅ Migrations applied!"

# Показать статус
migrate-status:
	@echo "📊 Migration status:"
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL_LOCAL)" version
	@echo ""
	@echo "📋 List of migrations:"
	dir $(MIGRATIONS_PATH)

# Помощь
help:
	@echo "Available commands:"
	@echo ""
	@echo "  make migrate-create name=xxx     - Create new migration"
	@echo "  make migrate-up-local            - Apply migrations locally"
	@echo "  make migrate-down-local          - Rollback last migration"
	@echo "  make migrate-down-local-all      - Rollback ALL migrations"
	@echo "  make migrate-up                  - Apply migrations to production"
	@echo "  make migrate-down n=1            - Rollback N migrations"
	@echo "  make migrate-force version=1     - FIX dirty database state"
	@echo "  make migrate-clean               - Clean dirty state"
	@echo "  make migrate-reset               - Reset database (ALL DATA LOST)"
	@echo "  make migrate-version             - Show current version"
	@echo "  make migrate-status              - Show full status"
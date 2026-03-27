.PHONY: migrate-create migrate-up-local migrate-down-local migrate-up migrate-down migrate-force

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
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL_LOCAL)" up

# Откатить одну миграцию локально
migrate-down-local:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL_LOCAL)" down 1

# Применить все миграции (для production)
migrate-up:
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" up

# Откатить N миграций
migrate-down:
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" down $(n)

# Принудительно установить версию миграции
migrate-force:
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" force $(version)

# Показать текущую версию
migrate-version:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL_LOCAL)" version
.PHONY: migrate-create migrate-up migrate-down migrate-force
migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

migrate-up-local:
	migrate -path migrations -database "postgres://postgres:1234@localhost:5432/note_proj?sslmode=disable" up

migrate-down-local:
	migrate -path migrations -database "postgres://postgres:1234@localhost:5432/note_proj?sslmode=disable" down 1


migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down $(n)

migrate-force:
	migrate -path migrations -database "$(DATABASE_URL)" force $(version)

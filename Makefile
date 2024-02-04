.PHONY: dev prepare-db

dev:
	sqlc generate; go run cmd/server/main.go

prepare-db:
	go run cmd/prepare_db/main.go
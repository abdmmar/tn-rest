.PHONY: dev prepare-db

dev:
	sqlc generate;\
	go run \
		cmd/server/main.go \
		cmd/server/router.go \
		cmd/server/national_park.go

prepare-db:
	go run cmd/prepare_db/main.go
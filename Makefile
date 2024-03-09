DB_USER=root
DB_PASSWORD=secret
DB_HOST=localhost
DB_PORT=5432
DB_NAME=simple_bank

init-migration:
	@echo "Creating migration"
	migrate create -ext sql -dir db/migration -seq init_schema

migrate-up:
	@echo "Migrating up"
	migrate -path db/migration -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -verbose up

migrate-down:
	@echo "Migrating down"
	migrate -path db/migration -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -verbose down

create-db:
	@echo "Creating database"
	docker go-k8s-db-1 exec createdb -U $(DB_USER) --owner=$(DB_USER) $(DB_NAME)

drop-db:
	@echo "Dropping database"
	docker go-k8s-db-1 exec dropdb -U $(DB_USER) $(DB_NAME)

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: create-db drop-db init-migration migrate-up migrate-down sqlc generate
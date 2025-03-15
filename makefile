createdb:
	docker exec -it minhdang2803 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it minhdang2803 dropdb simple_bank

postgres:
	docker run --name minhdang2803 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14-alpine

migrate_up:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migrate_down:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

create_migration: 
	migrate create -ext sql - dir db/migration -seq init_schema
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
server:
	go run main.go

.PHONY: sqlc createdb dropdb postgres migrate_down migrate_up create_migration test server
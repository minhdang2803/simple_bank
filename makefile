createdb:
	docker exec -it minhdang2803 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it minhdang2803 dropdb simple_bank

postgres:
	docker run --name minhdang2803 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14-alpine

migrate_up:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migrate_up1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migrate_down:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

migrate_down1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

create_migration: 
	migrate create -ext sql - dir db/migration -seq init_schema
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
server:
	go run main.go
mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/minhdang2803/simple_bank/db/sqlc Store 
.PHONY: sqlc createdb dropdb postgres migrate_down migrate_up migrate_down1 migrate_up1 create_migration test server mockgen
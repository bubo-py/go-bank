run_postgres:
	docker run --name postgres16 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:16-alpine

createdb:
	docker exec -it postgres16 createdb --username=root --owner=root go-bank
dropdb:
	docker exec -it postgres16 dropdb go-bank

migrateup:
	migrate -path db/migration -database "postgresql://root:password@localhost:5432/go-bank?sslmode=disable" --verbose up
migratedown:
	migrate -path db/migration -database "postgresql://root:password@localhost:5432/go-bank?sslmode=disable" --verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: run_postgres postgres_createdb postgres_dropdb migrateup migratedown sqlc test
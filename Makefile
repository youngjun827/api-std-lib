format:
	gofmt -s -w .

postgres:
	docker run --name api-std-lib-postgres --network api-std-lib-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	docker exec -it api-std-lib-postgres createdb --username=root --owner=root api-std-lib-db

dropdb:
	docker exec -it api-std-lib-postgres dropdb api-std-lib-db

migrateup:
	migrate -path internal/database/migration -database "postgresql://root:secret@localhost:5432/api-std-lib-db?sslmode=disable" -verbose up

migrateup1:
	migrate -path internal/database/migration -database "postgresql://root:secret@localhost:5432/api-std-lib-db?sslmode=disable" -verbose up 1

migratedown:
	migrate -path internal/database/migration -database "postgresql://root:secret@localhost:5432/api-std-lib-db?sslmode=disable" -verbose down

migratedown1:
	migrate -path internal/database/migration -database "postgresql://root:secret@localhost:5432/api-std-lib-db?sslmode=disable" -verbose down 1

test:
	go test -v -cover ./...

server:
	go run main.go

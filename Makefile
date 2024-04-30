BINARY=inventory
.DEFAULT_GOAL := run

build:
	@go build -o bin/${BINARY}

run:build
	@./bin/${BINARY}

clean:
	rm -rf bin

dbup:
	@docker compose up -d --remove-orphans

dbdown:
	@docker compose down

db_shell:
	@docker exec -it inventory psql -U anomy inventory

createdb:
	@docker exec -it inventory createdb --username=anomy --owner=anomy inventory

dropdb:
	@docker exec -it inventory dropdb --username=anomy inventory

migrate_up:
	@migrate -path misc/migrations -database "postgresql://anomy:secret@localhost:5432/inventory?sslmode=disable" up

migrate_down:
	@migrate -path misc/migrations -database "postgresql://anomy:secret@localhost:5432/inventory?sslmode=disable" down

sqlc:
	@sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

.PHONY: dbup dbdown createdb dropdb migrate_up migrate_down test server
.PHONY: all deps run test docker-up docker-down db-exec logs

all: run

deps:
	go mod tidy
	go mod download

run:
	go run cmd/main.go

test:
	go test ./tests/... -v

docker-up:
	docker compose up -d

docker-down:
	docker compose down -v

db-exec:
	docker exec -it arcs-db psql -U admin arcs

logs:
	docker compose logs -f

logs-%:
	docker compose logs -f $*
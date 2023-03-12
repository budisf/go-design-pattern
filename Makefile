MPATH=app/databases/migrations
DSN=postgres://root:secret@localhost:5433/econolab_ethical?sslmode=disable
PG_USER=root
PG_PASSWORD=secret
PG_PORT=5433
DB_NAME=ethical
# Testing purpose on docker

run:
	nodemon --exec go run main.go --signal SIGTERM

postgres-start:
	docker run --name postgres12 -p $(PG_PORT):5432 -e POSTGRES_USER=$(PG_USER) -e POSTGRES_PASSWORD=$(PG_PASSWORD) -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=$(PG_USER) --owner=$(PG_USER) $(DB_NAME)

dropdb:
	docker exec -it postgres12 dropdb $(DB_NAME)

postgres-stop:
	docker stop postgres12

postgres-delete:
	docker rm postgres12

migrate-up:
	migrate -path $(MPATH) -database "$(DSN)" up

migrate-down:
	migrate -path $(MPATH) -database "$(DSN)" down

# .PHONY: postgres createdb dropdb
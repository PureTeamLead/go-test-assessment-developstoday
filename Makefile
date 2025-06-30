include .env

APP-NAME = cat-app
DB_USER = ${POSTGRES_USER}
DB_PASSWORD = ${POSTGRES_PASSWORD}
DB_NAME = ${POSTGRES_DB}

.PHONY:

mig-down: db/migrations
	migrate -path ./db/migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@localhost:5432/${DB_NAME}?sslmode=disable" down

build: cmd/${APP-NAME}/main.go internal
	go build -o ./bin/${APP-NAME} ./cmd/${APP-NAME}/main.go

run: build
	bin/${APP-NAME}

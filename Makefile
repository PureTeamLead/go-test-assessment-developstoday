APP-NAME = cat-app

.PHONY:

mig-down: db/migrations
	migrate -path ./db/migrations -database "postgres://themaxs:1234@localhost:5432/assessment-db?sslmode=disable" down

build: cmd/${APP-NAME}/main.go internal
	go build -o ./bin/${APP-NAME} ./cmd/${APP-NAME}/main.go

run: build
	bin/${APP-NAME}
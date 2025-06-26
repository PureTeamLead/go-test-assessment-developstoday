.PHONY:

mig-down: db/migrations
	migrate -path ./db/migrations -database "postgres://themaxs:1234@localhost:5432/assessment-db?sslmode=disable" down

build:

run:
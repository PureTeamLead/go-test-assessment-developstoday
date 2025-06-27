FROM golang:1.24-alpine AS modules

WORKDIR /app

COPY go.mod go.sum /app/

RUN go mod download

FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY --from=modules . /app
COPY . /app

RUN go build -o ./bin/app ./cmd/cat-app/main.go

FROM alpine

WORKDIR /app

COPY --from=builder ./app/bin/app /app

CMD ["./app"]


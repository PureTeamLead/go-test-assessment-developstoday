FROM golang:1.24-alpine AS modules

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY --from=modules /app/go.mod ./go.mod
COPY --from=modules /app/go.sum ./go.sum
COPY . .

RUN go build -o /bin/app ./cmd/cat-app/main.go

FROM alpine

WORKDIR /cat-app

COPY --from=builder ./bin/app ./bin/app
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/db ./db

CMD ["./bin/app"]


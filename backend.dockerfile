FROM golang:alpine AS builder

ENV CGO_ENABLED=1

RUN apk add --no-cache \
    gcc \
    musl-dev \
    sqlite-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN rm -rf web

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

RUN go build -o main main.go

FROM alpine

WORKDIR /app

COPY --from=builder /app/main ./
COPY --from=builder /go/bin/goose /usr/bin/goose
COPY --from=builder /app/migrations ./migrations

CMD goose -dir ./migrations sqlite3 "$DATABASE_FILE" up && ./main

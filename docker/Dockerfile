FROM golang:1.24.0

COPY config /app/config
COPY src /app/src

WORKDIR /app/src

RUN go build cmd/migrate/migrate.go
RUN go build cmd/main/main.go

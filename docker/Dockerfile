FROM golang:1.23.4

COPY config /app/config
COPY src /app/src

WORKDIR /app/src

RUN go build cmd/migrate.go
RUN go build cmd/main.go

########## BUILD STAGE ##########
FROM golang:1.24-alpine AS builder

WORKDIR /build

# Зависимости
COPY src/go.mod src/go.sum ./
RUN go mod download

# Исходники
COPY src .

# Сборка бинарников
RUN CGO_ENABLED=0 go build -o /out/api ./cmd/main
RUN CGO_ENABLED=0 go build -o /out/migrate ./cmd/migrate


########## RUNTIME STAGE ##########
FROM alpine:3.19

WORKDIR /app

# non-root пользователь
RUN addgroup -S nonroot \
 && adduser  -S nonroot -G nonroot

# Runtime файлы
COPY entrypoint.sh /app/entrypoint.sh
COPY config /app/config
COPY --from=builder /out/api /app/api
COPY --from=builder /out/migrate /app/migrate
COPY src/assets /app/assets

# Права
RUN chmod +x /app/entrypoint.sh \
 && chown -R nonroot:nonroot /app

USER nonroot

ENV CONFIG_PATH=/app/config/config.yaml

EXPOSE 8000

ENTRYPOINT ["/bin/sh", "/app/entrypoint.sh"]

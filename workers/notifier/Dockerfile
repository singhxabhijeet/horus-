# --- Build Stage ---
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /notifier-worker .

# --- Final Stage ---
FROM alpine:latest AS final
WORKDIR /
COPY --from=builder /notifier-worker /notifier-worker
ENTRYPOINT ["/notifier-worker"]
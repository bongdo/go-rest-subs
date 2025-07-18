# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.24.2-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server/main.go

# Final stage
FROM alpine:3.20


RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app


COPY --from=builder /server /app/server
COPY migrations /app/migrations
RUN chown -R appuser:appgroup /app


USER appuser

EXPOSE 8080

ENTRYPOINT ["/app/server"]

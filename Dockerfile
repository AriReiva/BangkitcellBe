# --- Build Stage ---
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh folder project
COPY . .

# Build binary
RUN go build -o app .

# --- Runtime Stage ---
FROM alpine:latest

WORKDIR /app

# Copy binary dari builder
COPY --from=builder /app/app .

EXPOSE 8000

CMD ["./app"]
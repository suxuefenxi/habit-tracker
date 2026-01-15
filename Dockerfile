# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o habit-tracker cmd/server/main.go

# Run stage
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/habit-tracker .
COPY --from=builder /app/web ./web

# Expose port
EXPOSE 8080

# Run the application
CMD ["./habit-tracker"]

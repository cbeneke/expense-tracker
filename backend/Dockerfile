# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install necessary build tools
RUN apk add --no-cache gcc musl-dev

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

# Development stage
FROM builder AS development

# Enable docker healthcheck
RUN apk add --no-cache curl

CMD ["go", "run", "./cmd/server"]

# Production stage
FROM alpine:latest AS production

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .

# Install necessary runtime dependencies
RUN apk add --no-cache curl

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"] 
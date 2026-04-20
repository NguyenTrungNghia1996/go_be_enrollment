# Stage 1: Build
FROM golang:1.25-alpine AS builder

# Set up work directory
WORKDIR /app

# Install git and tzdata
RUN apk add --no-cache git tzdata

# Copy go mod and sum files
COPY go.mod go.sum ./

# Check if vendor exists and use it, otherwise download
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app (using static linking for alpine)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api-server ./cmd/server/main.go

# Stage 2: Final image
FROM alpine:latest

# Install minimal certificates and timezone data
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app/

# Copy the binary from builder
COPY --from=builder /app/api-server .

# Copy environment template if needed or rely on docker-compose env vars
# We also copy migrations in case your app relies on them programmatically
COPY --from=builder /app/migrations ./migrations

# Expose the application port
EXPOSE 8080

# Run the app
CMD ["./api-server"]

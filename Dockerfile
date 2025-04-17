# Multi-stage build for the Go IP check service
FROM golang:1.22-alpine AS builder

# Set working directory for the build
WORKDIR /app

# Install git (needed to download Go module dependencies if not vendored)
RUN apk add --no-cache git

# Copy go.mod and go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application's source code
COPY . .

# Build the Go application (statically linked binary for Linux)
ENV CGO_ENABLED=0
RUN go build -ldflags="-s -w" -o /app/ip-check-service ./main.go

# Second stage: create a minimal image for running the service
FROM scratch

# Set working directory in the final image
WORKDIR /app

# Copy the compiled binary and GeoLite2 database from the builder stage
COPY --from=builder /app/ip-check-service .
COPY --from=builder /app/GeoLite2-Country.mmdb .

# Use an unprivileged user for security (use numeric UID since scratch has no /etc/passwd)
USER 65534

# Expose HTTP and gRPC ports
EXPOSE 8080 9090

# Run the service by default
ENTRYPOINT ["./ip-check-service"]
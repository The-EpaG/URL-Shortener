# =========================================================================
# Build stage
#
# This stage compiles the Go application.
# It uses a specific version of the golang:alpine image to ensure
# reproducible builds.
# =========================================================================
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies for CGO.
RUN apk add --no-cache gcc musl-dev

# Copy go.mod and go.sum first to leverage Docker's layer caching.
# This layer is only rebuilt when dependencies change.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY main.go .
COPY internal /app/internal

# Build the Go application.
# CGO_ENABLED=1 is required for go-sqlite3.
# -ldflags="-w -s" strips debug information, reducing the binary size.
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-w -s" -o /app/main .

# =========================================================================
# Final stage
#
# This stage creates the final, minimal, and secure image.
# It starts from a minimal Alpine base image.
# =========================================================================
FROM alpine

# Create a dedicated user and group for the application to run as.
# Running as a non-root user is a critical security best practice.
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Set the working directory
WORKDIR /app

# Install only the required runtime dependencies (sqlite-libs and curl for healthcheck).
RUN apk add --no-cache sqlite-libs curl

# Copy the compiled binary from the builder stage.
COPY --from=builder /app/main /app/main

# Give the non-root user ownership of the app directory and its contents.
# This is necessary so the application can create the SQLite database file at runtime.
RUN chown -R appuser:appgroup /app

# Switch to the non-root user
USER appuser

EXPOSE 8080

# This allows Docker to monitor the container's health.
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/healthz || exit 1

CMD ["/app/main"]

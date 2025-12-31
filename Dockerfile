# Stage 1: Build frontend
FROM node:20-alpine AS frontend-builder

WORKDIR /frontend

# Copy frontend files
COPY frontend/package*.json ./
RUN npm install

COPY frontend/ ./
RUN npm run build

# Stage 2: Build backend
FROM golang:1.21-alpine AS backend-builder

# Install build dependencies for CGO (needed for sqlite3)
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

# Copy Go module files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and built frontend
COPY . .
COPY --from=frontend-builder /frontend/dist ./frontend/dist

# Build the binary
ARG VERSION=docker
ENV CGO_CFLAGS="-D_LARGEFILE64_SOURCE"
RUN CGO_ENABLED=1 go build -tags "sqlite_omit_load_extension" -ldflags "-X main.Version=${VERSION}" -o days .

# Stage 3: Runtime
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates sqlite-libs curl

WORKDIR /app

# Copy binary from builder
COPY --from=backend-builder /app/days .

# Create data directory
RUN mkdir -p /app/data

# Expose port
EXPOSE 8080

# Set environment variables
ENV DB_PATH=/app/data/journal.db
ENV PORT=8080

# Healthcheck
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD curl -f http://localhost:8080/ || exit 1

# Run the application
CMD ["./days"]

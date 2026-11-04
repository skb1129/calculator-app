# ==========================================
# Stage 1: Build Frontend React SPA
# ==========================================
FROM node:22-alpine AS frontend-builder
WORKDIR /build

COPY frontend/package*.json ./
RUN npm ci || npm install

COPY frontend/ ./
RUN npm run build

# ==========================================
# Stage 2: Build Go Backend Binary
# ==========================================
FROM golang:1.24-alpine AS backend-builder
WORKDIR /build

# Install CA certificates for secure external calls if needed
RUN apk add --no-cache ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /build/server ./cmd/server

# ==========================================
# Stage 3: Minimal Production Runtime
# ==========================================
FROM alpine:3.21 AS runtime

WORKDIR /app

# Add non-root user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy CA certs
COPY --from=backend-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy Go binary and static frontend assets
COPY --from=backend-builder /build/server /app/server
COPY --from=frontend-builder /build/dist /app/frontend/dist

# Set ownership
RUN chown -R appuser:appgroup /app

USER appuser

ENV PORT=8080
ENV STATIC_DIR=/app/frontend/dist
ENV GIN_MODE=release

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/v1/health || exit 1

ENTRYPOINT ["/app/server"]

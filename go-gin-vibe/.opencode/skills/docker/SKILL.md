---
name: docker
description: Use when creating or modifying Docker configurations, Dockerfiles, docker-compose files, or containerization. Covers multi-stage builds, optimization, and Docker best practices.
---

# Docker Skill

## Dockerfile

### Multi-stage Build (Go)

```dockerfile
# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/api/main.go

# Run stage
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/.env.example .

EXPOSE 8080

CMD ["./main"]
```

### Docker Compose

```yaml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - APP_ENV=development
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_NAME=alfatihah_center
      - DB_USER=postgres
      - DB_PASSWORD=password
      - REDIS_HOST=redis
      - REDIS_PORT=6379
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    volumes:
      - .:/app
    command: go run cmd/api/main.go

  postgres:
    image: postgres:16-alpine
    ports:
      - "5432:5432"
    environment:
      - POSTGRES_DB=alfatihah_center
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=password
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 5s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      timeout: 5s
      retries: 5

  worker:
    build: .
    environment:
      - REDIS_HOST=redis
      - REDIS_PORT=6379
    depends_on:
      redis:
        condition: service_healthy
    command: go run cmd/worker/main.go

volumes:
  postgres_data:
  redis_data:
```

## Commands

### Build image

```bash
docker build -t alfatihah-center .
```

### Run container

```bash
docker run -p 8080:8080 alfatihah-center
```

### Docker Compose

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f app

# Stop all services
docker-compose down

# Rebuild
docker-compose up -d --build
```

## Best Practices

1. Use multi-stage builds to reduce image size
2. Copy go.mod/go.sum first for better caching
3. Use .dockerignore to exclude unnecessary files
4. Run as non-root user in production
5. Use health checks
6. Use named volumes for persistence
7. Pin image versions
8. Use alpine for smaller images
9. Set proper environment variables
10. Use docker-compose for local development

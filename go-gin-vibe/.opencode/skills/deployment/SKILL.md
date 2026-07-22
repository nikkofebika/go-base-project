---
name: deployment
description: Use when setting up CI/CD pipelines, deploying applications, or configuring production environments. Covers GitHub Actions, deployment strategies, and production best practices.
---

# Deployment Skill

## GitHub Actions

### CI Workflow

```yaml
name: CI

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    
    services:
      postgres:
        image: postgres:16
        env:
          POSTGRES_DB: test_db
          POSTGRES_USER: postgres
          POSTGRES_PASSWORD: password
        ports:
          - 5432:5432
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
      
      redis:
        image: redis:7
        ports:
          - 6379:6379
        options: >-
          --health-cmd "redis-cli ping"
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5

    steps:
      - uses: actions/checkout@v4
      
      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.24'
      
      - name: Download dependencies
        run: go mod download
      
      - name: Run linter
        uses: golangci/golangci-lint-action@v4
        with:
          version: latest
      
      - name: Run tests
        env:
          DB_HOST: localhost
          DB_PORT: 5432
          DB_NAME: test_db
          DB_USER: postgres
          DB_PASSWORD: password
          REDIS_HOST: localhost
          REDIS_PORT: 6379
        run: go test -race -coverprofile=coverage.out ./...
      
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          file: ./coverage.out
```

### CD Workflow

```yaml
name: CD

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v4
      
      - name: Build and push Docker image
        uses: docker/build-push-action@v5
        with:
          push: true
          tags: registry.alfatihah-center.com/api:${{ github.sha }}
      
      - name: Deploy to production
        run: |
          kubectl set image deployment/api \
            api=registry.alfatihah-center.com/api:${{ github.sha }}
```

## Environment Variables

### Production

```env
APP_ENV=production
APP_PORT=8080

DB_HOST=production-db-host
DB_PORT=5432
DB_NAME=alfatihah_center
DB_USER=production_user
DB_PASSWORD=production_password
DB_SSL_MODE=require

REDIS_HOST=production-redis-host
REDIS_PORT=6379

JWT_SECRET=production-secret-key
JWT_ISSUER=alfatihah-center

ACCESS_TOKEN_TTL=15m
REFRESH_TOKEN_TTL=720h
```

## Deployment Checklist

### Pre-deployment

- [ ] All tests passing
- [ ] Linter passing
- [ ] Migration files ready
- [ ] Environment variables configured
- [ ] Docker image built
- [ ] Database backed up

### Deployment

- [ ] Run migrations
- [ ] Deploy new version
- [ ] Health check passing
- [ ] Logs monitoring
- [ ] Error tracking

### Post-deployment

- [ ] Verify all endpoints
- [ ] Check error rates
- [ ] Monitor performance
- [ ] Notify team

## Production Best Practices

1. Use environment variables for configuration
2. Never commit secrets to git
3. Use secrets management (Vault, AWS Secrets Manager)
4. Implement health checks
5. Use structured logging
6. Set up error tracking (Sentry, Bugsnag)
7. Monitor application metrics
8. Use blue-green or canary deployments
9. Implement rollback strategy
10. Set up alerts for critical errors

# AGENTS.md

## Project

**Project Name:** Alfatihah Center

You are a **Senior Golang Backend Engineer** responsible for designing and implementing a production-ready REST API.

Your primary goal is to produce code that is:

- Clean
- Maintainable
- Modular
- Testable
- Secure
- Scalable

This project is expected to grow into a large production system. Every implementation must prioritize maintainability over speed.

---

# Core Principles

Always follow these principles:

- SOLID
- Clean Architecture
- Clean Code
- DRY
- KISS
- YAGNI
- Composition over Inheritance
- Interface Segregation
- Dependency Injection

Never sacrifice architecture quality for shorter code.

---

# Tech Stack

Use only these technologies unless explicitly instructed otherwise.

- Go
- Gin
- GORM
- PostgreSQL
- golang-migrate
- Zerolog
- Swaggo
- go-playground/validator
- golang-jwt
- bcrypt
- hibiken/asynq
- caarlos0/env

---

# Architecture

Follow **Clean Architecture**.

Every feature must be separated into layers.

```
HTTP Request

↓

Handler

↓

Service

↓

Repository

↓

Database
```

Dependencies must only point downward.

Never bypass layers.

Example:

❌ Handler → Repository

❌ Handler → Database

❌ Service → HTTP

✅ Handler → Service

✅ Service → Repository

---

# Feature Module Structure

Every module must follow this structure.

```
internal/

    user/

        dto.go

        entity.go

        repository.go

        repository_impl.go

        service.go

        service_impl.go

        handler.go

        routes.go

        mapper.go

        errors.go
```

If a module grows significantly, split files into subfolders while preserving separation of concerns.

---

# Responsibility of Each Layer

## Handler

Handler is responsible only for:

- parsing request
- validating request
- calling service
- returning response

Handler must NOT contain:

- business logic
- SQL
- GORM query
- transaction logic

---

## Service

Service contains all business logic.

Service may:

- call repositories
- call other services
- perform validation
- perform authorization
- manage transactions

Service must NOT know anything about HTTP.

Never return HTTP status codes from service.

---

## Repository

Repository only communicates with database.

Repository responsibilities:

- CRUD
- Query
- Pagination
- Filtering

Repository must NEVER contain business logic.

---

# Dependency Injection

Always use dependency injection.

Never instantiate dependencies inside handlers.

Good

```
UserHandler
    ↓

UserService
    ↓

UserRepository
```

Bad

```
handler creates repository

handler creates database

handler creates service
```

---

# Interface Rules

Every Service must have an interface.

Every Repository must have an interface.

Use interfaces where mocking/testing is expected.

Avoid unnecessary abstractions.

---

# Context

Always propagate context.Context.

Example

```
Handler

↓

Service(ctx)

↓

Repository(ctx)
```

Never use context.Background() inside business logic.

---

# Error Handling

Never panic.

Return errors.

Wrap errors with context.

Example

```
fmt.Errorf("create user: %w", err)
```

Never ignore errors.

Never return raw database errors directly to API consumers.

Create domain-specific errors where appropriate.

---

# Logging

Use Zerolog.

Every request must include:

- request_id
- user_id (if authenticated)
- IP address
- method
- path
- latency
- status code

Log unexpected errors.

Never log:

- password
- token
- refresh token
- JWT
- secret
- authorization header

Sensitive information must never appear in logs.

---

# Validation

Use:

go-playground/validator

Validation must happen before business logic.

Never perform repetitive manual validation.

Create reusable validators whenever possible.

---

# Authentication

Authentication uses JWT.

Access Token:

- short-lived

Refresh Token:

- stored in database
- hashed before storage

Never store plain refresh tokens.

Forgot Password tokens follow the same principle.

---

# Authorization

Use custom RBAC middleware.

Permission format:

```
user.read

user.create

user.update

user.delete
```

Never hardcode permissions inside handlers.

Authorization belongs in middleware and service layer.

---

# Database

Database:

PostgreSQL

ORM:

GORM

Migration:

golang-migrate

Rules:

- never edit executed migration
- always create new migration
- migration must be reversible
- foreign keys must be explicit
- indexes must be added where appropriate

---

# Transactions

Whenever multiple tables are modified:

Use database transactions.

Example:

Create User

↓

Create User Detail

↓

Assign Default Role

All operations must succeed or rollback.

---

# Naming Convention

Packages

```
user
auth
role
permission
```

Do not use:

```
helpers
utils
common
misc
```

unless absolutely necessary.

Use meaningful names.

Bad

```
Do()

Handle()

Manager()

```

Good

```
CreateUser()

AssignRole()

RefreshToken()

```

---

# DTO

Separate DTO from entity.

Never expose database models directly.

Always map:

Request DTO

↓

Entity

↓

Response DTO

---

# API Response

Success

```json
{
  "data": {}
}
```

Collection

```json
{
  "data": []
}
```

Error

```json
{
  "message": "Validation failed.",
  "errors": {
    "email": ["The email field is required."]
  }
}
```

Keep response format consistent.

---

# Pagination

List endpoint must support:

- page
- per_page

Response:

```json
{
  "data": [],
  "meta": {
    "current_page": 1,
    "per_page": 10,
    "total": 100,
    "last_page": 10,
    "from": 1,
    "to": 10
  }
}
```

Avoid loading unnecessary records.

---

# Swagger

Every endpoint must contain Swaggo annotations.

Swagger documentation must always be updated after endpoint changes.

---

# Queue

Background jobs use:

hibiken/asynq

Use queue for:

- forgot password email
- email notification
- future background jobs

Never perform long-running jobs inside HTTP handlers.

---

# Security

Always:

- hash passwords using bcrypt
- hash refresh tokens
- validate JWT signature
- validate input
- sanitize user input
- use parameterized queries through GORM
- avoid leaking internal errors

Never trust client input.

---

# Code Quality

Prefer small functions.

Recommended function size:

20–40 lines.

Split functions when they become difficult to read.

Avoid deep nesting.

Prefer early return.

Bad

```
if ...
    if ...
        if ...
```

Good

```
if err != nil {
    return err
}
```

---

# Constants

Avoid magic strings.

Prefer constants.

Example

```
TokenTypeRefresh

TokenTypeForgotPassword

RoleAdmin

PermissionUserCreate
```

---

# Configuration

Use:

caarlos0/env

Environment variables must be loaded only once during application startup.

Do not access environment variables directly throughout the application.

---

# Testing

Write code that is testable.

Business logic must not depend on HTTP.

Prefer dependency injection to enable mocking.

Avoid tightly coupled implementations.

---

# Performance

Avoid N+1 queries.

Use preload only when needed.

Select only required columns.

Use indexes appropriately.

Avoid unnecessary allocations.

---

# General Rules

Always write production-ready code.

Always think about future scalability.

Prefer readability over clever code.

If there are multiple possible implementations:

Choose the one that is easiest to maintain.

When generating new code:

- follow existing project structure
- reuse existing abstractions
- avoid duplication
- keep consistency across modules

Never rewrite unrelated files.

Never introduce breaking changes unless explicitly requested.

Always explain architectural decisions briefly when introducing a new pattern or dependency.

The codebase should be understandable by another senior Go engineer without additional explanation.

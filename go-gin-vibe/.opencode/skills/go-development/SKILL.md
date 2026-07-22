---
name: go-development
description: Use when writing, reviewing, or refactoring Go code in this project. Covers Clean Architecture patterns, module structure, dependency injection, and idiomatic Go conventions.
---

# Go Development Skill

## Project Structure

This project uses Modular Monolith with Clean Architecture:

```
internal/
├── auth/          # Authentication module
├── user/          # User management module
├── role/          # Role management module
├── permission/    # Permission management module
├── token/         # Token management module
├── common/        # Shared packages
│   ├── config/
│   ├── database/
│   ├── logger/
│   ├── middleware/
│   ├── response/
│   ├── validator/
│   ├── pagination/
│   ├── constants/
│   └── errors/
└── bootstrap/
```

## Module Structure

Each module follows:

```
module/
├── handler/
│   └── handler.go
├── service/
│   ├── service.go
│   └── service_impl.go
├── repository/
│   ├── repository.go
│   └── repository_impl.go
├── entity/
│   └── entity.go
├── dto/
│   └── dto.go
├── mapper/
│   └── mapper.go
├── routes.go
└── module.go
```

## Rules

1. Handler only: parse request, validate, call service, return response
2. Service: all business logic, no HTTP awareness
3. Repository: only database access
4. Use dependency injection via constructors
5. Always propagate context.Context
6. Use DTO for request/response, never expose entities
7. Use typed constants, no magic strings
8. Wrap errors with context: fmt.Errorf("operation: %w", err)

## Code Style

- snake_case for files
- lowercase for packages
- Meaningful function names: CreateUser(), FindByEmail()
- Early return pattern
- Max 30-50 lines per function
- Use gofumpt for formatting

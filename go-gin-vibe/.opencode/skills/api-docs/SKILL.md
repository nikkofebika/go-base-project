---
name: api-docs
description: Use when creating or updating Swagger documentation for API endpoints. Covers Swaggo annotations and API response format.
---

# API Documentation Skill

## Tool

Swaggo

## Generate Swagger

```bash
swag init -g cmd/api/main.go
```

## Annotation Template

```go
// @Summary      Create user
// @Description  Create a new user account
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request  body      dto.CreateUserRequest  true  "User data"
// @Success      201      {object}  response.SuccessResponse{data=dto.UserResponse}
// @Failure      422      {object}  response.ValidationErrorResponse
// @Failure      409      {object}  response.ErrorResponse
// @Failure      500      {object}  response.ErrorResponse
// @Router       /api/v1/users [post]
// @Security     BearerAuth
```

## Response Format (Laravel-style)

### Success

```json
{"data": {}}
```

### Collection

```json
{"data": []}
```

### Pagination

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

### Error

```json
{
  "message": "Validation failed.",
  "errors": {
    "email": ["The email field is required."]
  }
}
```

## Tags Convention

Group endpoints by module:

- auth
- users
- roles
- permissions

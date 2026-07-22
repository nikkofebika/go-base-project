---
name: go-error-handling
description: Use when handling errors, creating custom errors, or wrapping errors. Covers error patterns, sentinel errors, and error handling best practices.
---

# Go Error Handling Skill

## Rules

1. Always check errors
2. Wrap errors with context
3. Use sentinel errors for expected errors
4. Use custom error types for complex errors
5. Never panic in library code
6. Return early on error

## Error Wrapping

```go
// Good
func CreateUser(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
    user, err := repo.FindByEmail(ctx, req.Email)
    if err != nil {
        return nil, fmt.Errorf("find user by email: %w", err)
    }
    // ...
}

// Bad
func CreateUser(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
    user, err := repo.FindByEmail(ctx, req.Email)
    if err != nil {
        return nil, err // No context
    }
    // ...
}
```

## Sentinel Errors

```go
// errors/errors.go
var (
    ErrUserNotFound       = errors.New("user not found")
    ErrEmailAlreadyExists = errors.New("email already exists")
    ErrInvalidCredential  = errors.New("invalid credential")
    ErrUnauthorized       = errors.New("unauthorized")
    ErrForbidden          = errors.New("forbidden")
)
```

## Custom Error Types

```go
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// Usage
if user.Email == "" {
    return &ValidationError{
        Field:   "email",
        Message: "is required",
    }
}
```

## Error Checking

```go
// Check specific error
if errors.Is(err, ErrUserNotFound) {
    // Handle not found
}

// Check error type
var validationErr *ValidationError
if errors.As(err, &validationErr) {
    // Handle validation error
    fmt.Println(validationErr.Field)
}
```

## Error Handling Patterns

```go
// Early return
func CreateUser(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
    if err := validate(req); err != nil {
        return nil, fmt.Errorf("validate request: %w", err)
    }

    user, err := repo.FindByEmail(ctx, req.Email)
    if err != nil {
        return nil, fmt.Errorf("find user: %w", err)
    }

    if user != nil {
        return nil, ErrEmailAlreadyExists
    }

    // Continue with creation
    return createNewUser(ctx, req)
}
```

## Never Ignore Errors

```go
// Bad
result, _ := doSomething()

// Good
result, err := doSomething()
if err != nil {
    return fmt.Errorf("do something: %w", err)
}
```

## Error Handling in Service Layer

```go
func (s *service) Create(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
    // Business validation
    existing, err := s.repo.FindByEmail(ctx, req.Email)
    if err != nil && !errors.Is(err, ErrUserNotFound) {
        return nil, fmt.Errorf("check existing user: %w", err)
    }

    if existing != nil {
        return nil, ErrEmailAlreadyExists
    }

    // Create user
    user := &entity.User{
        Email:    req.Email,
        Password: hashPassword(req.Password),
    }

    if err := s.repo.Create(ctx, user); err != nil {
        return nil, fmt.Errorf("create user: %w", err)
    }

    return mapper.ToResponse(user), nil
}
```

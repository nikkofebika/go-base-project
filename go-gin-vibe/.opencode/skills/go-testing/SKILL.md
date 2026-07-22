---
name: go-testing
description: Use when writing unit tests, integration tests, or benchmarks. Covers table-driven tests, test helpers, mocking, and test patterns.
---

# Go Testing Skill

## Commands

### Run all tests

```bash
go test ./...
```

### Run with race detector

```bash
go test -race ./...
```

### Run with coverage

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Run specific test

```bash
go test -run TestCreateUser ./internal/user/...
```

### Run benchmarks

```bash
go test -bench=. ./...
```

## Table-Driven Tests

```go
func TestCreateUser(t *testing.T) {
    tests := []struct {
        name    string
        input   dto.CreateUserRequest
        want    *dto.UserResponse
        wantErr error
    }{
        {
            name: "success",
            input: dto.CreateUserRequest{
                Email:    "test@example.com",
                Password: "password123",
            },
            want: &dto.UserResponse{
                Email: "test@example.com",
            },
            wantErr: nil,
        },
        {
            name: "email already exists",
            input: dto.CreateUserRequest{
                Email:    "existing@example.com",
                Password: "password123",
            },
            want:    nil,
            wantErr: ErrEmailAlreadyExists,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := service.Create(context.Background(), tt.input)
            if tt.wantErr != nil {
                assert.ErrorIs(t, err, tt.wantErr)
                return
            }
            assert.NoError(t, err)
            assert.Equal(t, tt.want, got)
        })
    }
}
```

## Test Helpers

```go
func setupTestDB(t *testing.T) *gorm.DB {
    t.Helper()
    db, err := gorm.Open(postgres.Open(testDSN), &gorm.Config{})
    require.NoError(t, err)
    t.Cleanup(func() {
        sqlDB, _ := db.DB()
        sqlDB.Close()
    })
    return db
}
```

## Mocking

Use interfaces for mocking:

```go
type UserRepository interface {
    Create(ctx context.Context, user *entity.User) error
    FindByEmail(ctx context.Context, email string) (*entity.User, error)
}

type mockUserRepository struct {
    createFn      func(ctx context.Context, user *entity.User) error
    findByEmailFn func(ctx context.Context, email string) (*entity.User, error)
}

func (m *mockUserRepository) Create(ctx context.Context, user *entity.User) error {
    return m.createFn(ctx, user)
}

func (m *mockUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
    return m.findByEmailFn(ctx, email)
}
```

## Assertions

Use testify for assertions:

```go
import (
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

// assert - continues test even on failure
assert.Equal(t, expected, actual)
assert.NoError(t, err)
assert.NotNil(t, result)

// require - stops test on failure
require.NoError(t, err)
require.NotNil(t, result)
```

## Best Practices

1. Test name should describe the scenario
2. One assertion per test case (or related assertions)
3. Use t.Helper() for test helpers
4. Use t.Cleanup() for cleanup
5. Run tests with -race flag
6. Use table-driven tests for multiple scenarios
7. Mock external dependencies, not internal logic
8. Test error cases as thoroughly as success cases

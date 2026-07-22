---
name: go-security
description: Use when implementing security features, handling authentication, authorization, or protecting against common vulnerabilities. Covers JWT, bcrypt, input validation, and security best practices.
---

# Go Security Skill

## Authentication

### JWT Best Practices

```go
// Token generation
func GenerateAccessToken(userID int64, secret string, ttl time.Duration) (string, error) {
    claims := jwt.MapClaims{
        "user_id": userID,
        "exp":     time.Now().Add(ttl).Unix(),
        "iat":     time.Now().Unix(),
        "iss":     "alfatihah-center",
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}

// Token validation
func ValidateAccessToken(tokenString, secret string) (jwt.MapClaims, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(secret), nil
    })

    if err != nil {
        return nil, fmt.Errorf("parse token: %w", err)
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok || !token.Valid {
        return nil, errors.New("invalid token")
    }

    return claims, nil
}
```

### Refresh Token Storage

```go
// Hash refresh token before storage
func HashToken(token string) (string, error) {
    hash := sha256.Sum256([]byte(token))
    return hex.EncodeToString(hash[:]), nil
}

// Verify refresh token
func VerifyToken(token, hashedToken string) bool {
    hash := sha256.Sum256([]byte(token))
    return hex.EncodeToString(hash[:]) == hashedToken
}
```

## Password Hashing

```go
// Hash password
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", fmt.Errorf("hash password: %w", err)
    }
    return string(bytes), nil
}

// Check password
func CheckPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

## Input Validation

```go
// Use go-playground/validator
type CreateUserRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8,max=72"`
    Name     string `json:"name" validate:"required,max=150"`
}

// Custom validation
func ValidateStruct(s interface{}) error {
    validate := validator.New()
    if err := validate.Struct(s); err != nil {
        for _, err := range err.(validator.ValidationErrors) {
            return fmt.Errorf("field %s failed on %s validation", err.Field(), err.Tag())
        }
    }
    return nil
}
```

## SQL Injection Prevention

```go
// Good - use parameterized queries
func FindByEmail(ctx context.Context, email string) (*User, error) {
    var user User
    err := db.WithContext(ctx).Where("email = ?", email).First(&user).Error
    return &user, err
}

// Bad - string concatenation (NEVER DO THIS)
func FindByEmail(ctx context.Context, email string) (*User, error) {
    var user User
    query := "SELECT * FROM users WHERE email = '" + email + "'"
    err = db.WithContext(ctx).Raw(query).Scan(&user).Error
    return &user, err
}
```

## CORS Configuration

```go
func SetupCORS() gin.HandlerFunc {
    return cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000"},
        AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    })
}
```

## Rate Limiting

```go
import "golang.org/x/time/rate"

func RateLimiter(rps int) gin.HandlerFunc {
    limiter := rate.NewLimiter(rate.Limit(rps), rps*2)
    return func(c *gin.Context) {
        if !limiter.Allow() {
            c.JSON(http.StatusTooManyRequests, gin.H{
                "message": "Too many requests",
            })
            c.Abort()
            return
        }
        c.Next()
    }
}
```

## Security Headers

```go
func SecurityHeaders() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("X-Content-Type-Options", "nosniff")
        c.Header("X-Frame-Options", "DENY")
        c.Header("X-XSS-Protection", "1; mode=block")
        c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        c.Header("Content-Security-Policy", "default-src 'self'")
        c.Next()
    }
}
```

## Environment Variables

```go
// Never log sensitive data
func LogRequest(c *gin.Context) {
    logger.Info().
        Str("request_id", c.GetString("request_id")).
        Str("method", c.Request.Method).
        Str("path", c.Request.URL.Path).
        Str("client_ip", c.ClientIP()).
        Int("status", c.Writer.Status()).
        Dur("latency", time.Since(start)).
        Msg("request completed")
    // Never log: passwords, tokens, secrets, authorization headers
}
```

## Best Practices

1. Always hash passwords with bcrypt
2. Hash refresh tokens before storage
3. Use parameterized queries
4. Validate all input
5. Use HTTPS in production
6. Implement rate limiting
7. Set security headers
8. Never log sensitive data
9. Use environment variables for secrets
10. Rotate secrets regularly

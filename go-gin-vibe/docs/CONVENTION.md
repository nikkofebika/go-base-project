# CONVENTION.md

# Alfatihah Center Go Coding Convention

Version: 1.0

Language: Go

Go Version: 1.24+

---

# 1. Purpose

Dokumen ini berisi standar penulisan kode yang wajib diikuti oleh seluruh developer maupun AI.

Tujuan utama:

- Konsisten
- Mudah dibaca
- Mudah di-maintain
- Mudah di-test
- Mengikuti idiomatic Go
- Production Ready

---

# 2. General Principles

Seluruh implementasi wajib mengikuti:

- SOLID
- Clean Code
- DRY
- KISS
- YAGNI
- Composition over Inheritance
- Interface Segregation
- Dependency Injection

---

# 3. Package Convention

Gunakan nama package yang singkat.

✅

```text
user

auth

role

permission

repository

service
```

❌

```text
User

UserService

UserRepository

Helpers

Utilities
```

Seluruh nama package menggunakan lowercase.

---

# 4. File Naming

Gunakan snake_case.

✅

```text
user_service.go

user_handler.go

jwt_middleware.go
```

❌

```text
UserService.go

UserHandler.go
```

---

# 5. Interface Convention

Gunakan interface hanya jika memang diperlukan.

Contoh

```go
type Service interface {
    Create(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error)
}
```

Implementasi

```go
type service struct {
    repository repository.Repository
}
```

Jangan membuat interface tanpa alasan.

---

# 6. Constructor

Seluruh dependency wajib menggunakan constructor.

```go
func New(
    repository repository.Repository,
) Service
```

Jangan membuat object menggunakan global variable.

---

# 7. Dependency Injection

Dependency harus di-inject melalui constructor.

Jangan membuat dependency di dalam Service maupun Handler.

❌

```go
repo := repository.New(db)
```

di dalam service.

---

# 8. Context

Seluruh layer wajib menerima

```go
context.Context
```

Contoh

```go
func (s *service) Create(ctx context.Context, req dto.CreateUserRequest)
```

Jangan menggunakan

```go
context.Background()
```

di business logic.

---

# 9. Handler

Handler hanya bertugas:

- Bind Request
- Validate Request
- Memanggil Service
- Mengembalikan Response

Handler tidak boleh memiliki business logic.

Handler tidak boleh mengakses database.

---

# 10. Service

Seluruh business logic berada pada Service.

Service boleh:

- Memanggil Repository
- Memanggil Service module lain
- Menjalankan Transaction
- Validasi Business Rule

Service tidak boleh mengetahui HTTP.

---

# 11. Repository

Repository hanya bertanggung jawab terhadap database.

Repository tidak boleh:

- Validasi
- Business Logic
- HTTP Response

---

# 12. DTO

Gunakan DTO untuk seluruh request dan response.

Entity tidak boleh dikembalikan langsung.

---

# 13. Mapper

Seluruh konversi Entity ↔ DTO dilakukan pada Mapper.

Jangan melakukan mapping di Handler.

---

# 14. Entity

Entity hanya digunakan untuk representasi database.

Entity tidak boleh memiliki business logic.

---

# 15. Error Handling

Seluruh function mengembalikan

```go
(value, error)
```

Gunakan early return.

```go
if err != nil {
    return nil, err
}
```

Jangan menggunakan panic untuk business error.

---

# 16. Custom Error

Gunakan custom error untuk business error.

Contoh

```go
ErrUserNotFound

ErrEmailAlreadyExists

ErrInvalidCredential
```

Jangan membandingkan error menggunakan string.

---

# 17. Logging

Gunakan Zerolog.

Yang wajib dicatat:

- Request ID
- User ID
- Method
- Path
- Status Code
- Duration
- Error

Jangan pernah mencatat:

- Password
- JWT
- Refresh Token
- Secret
- OTP

---

# 18. Validation

Gunakan

```text
go-playground/validator
```

Seluruh validasi request dilakukan pada DTO.

Business validation dilakukan pada Service.

---

# 19. Transaction

Transaction hanya boleh dimulai pada Service Layer.

Repository tidak boleh memulai transaction sendiri.

---

# 20. Constant

Gunakan typed constant.

```go
type TokenType string

const (
    TokenTypeRefresh TokenType = "refresh_token"
    TokenTypeForgotPassword TokenType = "forgot_password"
)
```

Jangan menggunakan magic string.

---

# 21. Configuration

Seluruh environment dibaca satu kali saat startup.

Gunakan

```text
caarlos0/env
```

Business logic tidak boleh memanggil

```go
os.Getenv()
```

---

# 22. Response

Seluruh response menggunakan helper yang sama.

Contoh

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

Validation Error

```json
{
  "message": "Validation failed.",
  "errors": {
    "email": ["The email field is required."]
  }
}
```

---

# 23. Pagination

Gunakan format

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

---

# 24. Naming Convention

Gunakan nama yang jelas.

✅

```go
CreateUser

UpdateRole

FindByEmail
```

❌

```go
Do

Exec

Run

HandleData
```

---

# 25. Function

Satu function hanya memiliki satu tanggung jawab.

Idealnya maksimal sekitar 30–50 baris.

Jika terlalu panjang, pecah menjadi function kecil.

---

# 26. Variable

Gunakan nama yang jelas.

✅

```go
user

role

permission
```

❌

```go
u

r

x

tmp
```

Gunakan nama pendek hanya untuk scope yang sangat kecil, misalnya iterator (`i`, `j`) atau receiver (`s`, `r`, `h`).

---

# 27. Comment

Gunakan comment hanya jika benar-benar diperlukan.

Kode yang baik harus dapat menjelaskan dirinya sendiri.

---

# 28. Testing

Business logic harus mudah di-unit test.

Prioritas testing:

- Service
- Repository
- Middleware

---

# 29. Don't

- Jangan menggunakan global variable.
- Jangan menggunakan panic untuk business error.
- Jangan membuat helper yang tidak jelas tanggung jawabnya.
- Jangan membuat package `utils` atau `helpers` sebagai tempat fungsi acak.
- Jangan menggunakan magic number.
- Jangan menggunakan magic string.
- Jangan melakukan query database di Handler.
- Jangan mengembalikan Entity langsung ke client.

---

# 30. Definition of Done

Kode dianggap memenuhi standar apabila:

- Mengikuti idiomatic Go.
- Menggunakan Dependency Injection.
- Menggunakan constructor.
- Menggunakan context.Context.
- Menggunakan DTO.
- Menggunakan Repository Pattern.
- Menggunakan Service Layer.
- Menggunakan typed constant.
- Menggunakan custom error.
- Menggunakan response wrapper.
- Menggunakan logging yang konsisten.
- Mudah di-test.
- Mudah dipindahkan menjadi microservice.

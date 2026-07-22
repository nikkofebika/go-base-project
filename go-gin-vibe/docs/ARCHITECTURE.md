# ARCHITECTURE.md

# Alfatihah Center Backend Architecture

Version: 1.0

---

# 1. Goals

Backend Alfatihah Center dibangun menggunakan **Modular Monolith Architecture**.

Target utama arsitektur ini adalah:

- Modular
- Scalable
- Maintainable
- Testable
- Production Ready
- Microservice Ready

Seluruh module harus dapat dipindahkan menjadi microservice dengan perubahan seminimal mungkin.

---

# 2. Architecture Principles

Seluruh implementasi wajib mengikuti prinsip berikut.

- SOLID
- Clean Architecture
- Separation of Concerns
- High Cohesion
- Low Coupling
- Dependency Injection
- Interface Driven Design
- Composition over Inheritance
- DRY
- KISS
- YAGNI

---

# 3. Modular Monolith

Project menggunakan **Package by Feature (Domain)**, bukan **Package by Layer**.

## ❌ Jangan

```text
handlers/
services/
repositories/
models/
```

Karena seluruh business logic akan bercampur.

## ✅ Gunakan

```text
auth/
user/
role/
permission/
token/
```

Setiap folder merupakan sebuah **Business Module**.

Module dianggap sebagai aplikasi kecil yang independen.

---

# 4. Module Boundary

Module adalah boundary.

Setiap module memiliki:

- Entity
- DTO
- Handler
- Service
- Repository
- Route
- Dependency Registration

Module lain tidak boleh mengakses implementasi internal module tersebut.

Komunikasi antar module hanya melalui **Service Interface**.

Contoh

```text
Auth Module

↓

User Service Interface

↓

User Module
```

Bukan

```text
Auth Module

↓

User Repository
```

Dengan pendekatan ini, implementasi User dapat diganti menjadi HTTP Client atau gRPC Client tanpa mengubah business logic Auth.

---

# 5. Project Structure

```text
cmd/
└── api/
    └── main.go

internal/

├── auth/
│   ├── handler/
│   ├── service/
│   ├── repository/
│   ├── entity/
│   ├── dto/
│   ├── mapper/
│   ├── routes.go
│   └── module.go
│
├── user/
│   ├── handler/
│   ├── service/
│   ├── repository/
│   ├── entity/
│   ├── dto/
│   ├── mapper/
│   ├── routes.go
│   └── module.go
│
├── role/
├── permission/
├── token/
│
├── common/
│   ├── config/
│   ├── database/
│   ├── middleware/
│   ├── logger/
│   ├── response/
│   ├── pagination/
│   ├── validator/
│   ├── constants/
│   └── errors/
│
└── bootstrap/
    ├── app.go
    ├── router.go
    └── provider.go

docs/

migrations/

scripts/

pkg/
```

---

# 6. Module Structure

Setiap module memiliki struktur berikut.

```text
user/

    handler/

    service/

    repository/

    entity/

    dto/

    mapper/

    routes.go

    module.go
```

Jika module berkembang besar, setiap package dapat memiliki beberapa file.

Contoh

```text
handler/

    create.go

    update.go

    delete.go

    list.go

    detail.go
```

---

# 7. Dependency Rule

Dependency hanya boleh mengarah ke bawah.

```text
HTTP

↓

Handler

↓

Service

↓

Repository

↓

Database
```

Tidak boleh ada dependency sebaliknya.

Repository tidak boleh mengetahui Service.

Handler tidak boleh mengetahui Database.

---

# 8. Layer Responsibility

## Handler

Bertanggung jawab untuk:

- Parse Request
- Validate Request
- Memanggil Service
- Mengembalikan Response

Tidak boleh memiliki business logic.

---

## Service

Berisi seluruh business logic.

Service dapat:

- Memanggil Repository
- Memanggil Service module lain
- Menjalankan transaction
- Melakukan authorization
- Menjalankan orchestration

Service tidak boleh mengetahui HTTP maupun Gin.

---

## Repository

Repository hanya bertugas mengakses database.

Repository tidak boleh memiliki business rule.

---

# 9. Dependency Injection

Seluruh dependency dibuat saat startup.

Contoh

```text
Database

↓

Repository

↓

Service

↓

Handler

↓

Router
```

Handler tidak boleh membuat Service.

Service tidak boleh membuat Repository.

Semua dependency diberikan melalui constructor.

---

# 10. Module Registration

Setiap module memiliki file

```text
module.go
```

yang bertugas melakukan dependency wiring.

Contoh

```text
User Module

↓

Repository

↓

Service

↓

Handler

↓

Routes
```

Bootstrap hanya bertugas memanggil module tersebut.

Dengan demikian bootstrap tidak mengetahui detail implementasi setiap module.

---

# 11. Request Lifecycle

```text
HTTP Request

↓

Recovery Middleware

↓

Request ID

↓

Logger

↓

Authentication

↓

Authorization

↓

Handler

↓

Service

↓

Repository

↓

PostgreSQL

↓

Response
```

---

# 12. Database Access

Seluruh query database harus melalui Repository.

Tidak boleh ada query GORM di:

- Handler
- Middleware
- Service

---

# 13. Transaction

Transaction hanya boleh dimulai pada Service Layer.

Contoh

Create User

↓

Create User Detail

↓

Assign Default Role

↓

Commit

Jika salah satu gagal maka seluruh perubahan harus rollback.

---

# 14. DTO

Seluruh Request dan Response menggunakan DTO.

Flow

```text
Request DTO

↓

Entity

↓

Business Logic

↓

Response DTO
```

Entity tidak boleh langsung dikirim ke client.

---

# 15. Shared Package

Package shared hanya digunakan untuk komponen lintas module.

Contoh

```text
common/

config/

database/

logger/

middleware/

pagination/

response/

validator/

constants/

errors/
```

Business logic tidak boleh berada di common.

---

# 16. Authentication

Authentication menggunakan JWT.

Refresh Token disimpan di database dalam bentuk hash.

Authentication berada pada Auth Module.

Module lain tidak boleh mengimplementasikan authentication sendiri.

---

# 17. Authorization

Authorization menggunakan Role Based Access Control (RBAC).

Flow

```text
User

↓

Role

↓

Permission

↓

Middleware

↓

Handler
```

Permission menggunakan format

```text
user.read
user.create
role.update
permission.delete
```

---

# 18. Queue

Background Job menggunakan Asynq.

Flow

```text
Service

↓

Enqueue Job

↓

Redis

↓

Worker

↓

Business Logic
```

HTTP Handler tidak boleh menjalankan proses berat secara langsung.

---

# 19. Logging

Seluruh request harus memiliki Request ID.

Minimal informasi yang dicatat:

- Request ID
- User ID
- Method
- Path
- Status Code
- Duration
- Client IP

Sensitive information tidak boleh dicatat.

---

# 20. Configuration

Configuration dibaca satu kali saat startup menggunakan `caarlos0/env`.

Business logic tidak boleh memanggil environment variable secara langsung.

---

# 21. Coding Guidelines

- Gunakan constructor untuk seluruh dependency.
- Gunakan interface pada Service dan Repository.
- Gunakan context.Context pada seluruh layer.
- Hindari global variable.
- Gunakan early return.
- Hindari function yang terlalu panjang.
- Hindari circular dependency.
- Hindari package `helpers` atau `utils` sebagai tempat fungsi acak.

---

# 22. Microservice Migration Strategy

Setiap module harus dianggap sebagai calon microservice.

Misalnya

```text
internal/user
```

nantinya dapat dipindahkan menjadi

```text
user-service
```

Langkah migrasi yang diharapkan hanya:

1. Memindahkan folder module ke repository baru.
2. Mengganti implementasi Service Interface menjadi HTTP atau gRPC Client.
3. Memindahkan database jika diperlukan.
4. Menyesuaikan dependency injection.

Business logic tidak boleh berubah.

---

# 23. Definition of Done

Sebuah implementasi dianggap sesuai arsitektur apabila:

- Tidak ada business logic di Handler.
- Tidak ada query database di Handler maupun Service.
- Tidak ada module yang mengakses Repository module lain.
- Seluruh komunikasi antar module melalui Service Interface.
- Seluruh dependency menggunakan constructor injection.
- Seluruh perubahan multi-table menggunakan transaction.
- Setiap module dapat dipindahkan menjadi microservice tanpa refactor besar.
- Struktur dan pola implementasi konsisten di seluruh project.

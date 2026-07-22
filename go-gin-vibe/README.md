# README.md

# Alfatihah Center Backend

Backend REST API untuk **Alfatihah Center**.

Project ini dibangun menggunakan **Golang** dengan arsitektur **Modular Monolith**, mengikuti prinsip **Clean Architecture**, **SOLID**, dan dirancang agar mudah dipisahkan menjadi **Microservices** di masa depan.

---

# Features

Current MVP

- JWT Authentication
- Refresh Token
- Forgot Password
- User Management
- Role Management
- Permission Management
- Role Based Access Control (RBAC)
- Swagger Documentation
- Background Job (Asynq)

Future

- Organization
- Donation
- Program
- Volunteer
- Event
- Notification
- Audit Log
- File Storage
- Payment Integration

---

# Tech Stack

| Component      | Technology              |
| -------------- | ----------------------- |
| Language       | Go                      |
| HTTP Framework | Gin                     |
| ORM            | GORM                    |
| Database       | PostgreSQL              |
| Migration      | golang-migrate          |
| Queue          | Hibiken Asynq           |
| Cache          | Redis                   |
| Authentication | JWT                     |
| Password Hash  | bcrypt                  |
| Validation     | go-playground/validator |
| Logging        | Zerolog                 |
| Documentation  | Swaggo                  |
| Configuration  | caarlos0/env            |

---

# Project Structure

```text
cmd/
└── api/
    └── main.go

internal/

├── auth/
├── user/
├── role/
├── permission/
├── token/

├── common/
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

docs/

migrations/

scripts/

pkg/
```

---

# Architecture

Project menggunakan **Modular Monolith**.

Setiap module merupakan business domain yang independen.

Communication Flow

```text
HTTP

↓

Handler

↓

Service

↓

Repository

↓

PostgreSQL
```

Business logic hanya berada pada Service Layer.

Repository hanya bertanggung jawab terhadap database.

Handler hanya menangani HTTP Request dan Response.

---

# Documentation

Seluruh dokumentasi project tersedia pada folder `docs/`.

| File                    | Description                       |
| ----------------------- | --------------------------------- |
| docs/AGENTS.md          | AI Coding Guideline               |
| docs/PRD.md             | Product Requirement Document      |
| docs/ARCHITECTURE.md    | Technical Architecture            |
| docs/CONVENTION.md      | Go Coding Convention              |
| docs/DATABASE.md        | Database Convention               |
| docs/API.md             | API Convention                    |
| docs/ERD.md             | Database Schema                   |

---

# Requirements

- Go 1.24+
- PostgreSQL 16+
- Redis 7+
- Docker (optional)

---

# Installation

Clone repository

```bash
git clone <repository-url>

cd alfatihah-center
```

Install dependency

```bash
go mod tidy
```

---

# Environment

Copy environment file

```bash
cp .env.example .env
```

Example

```env
APP_NAME=Alfatihah Center
APP_ENV=development
APP_PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_NAME=alfatihah_center
DB_USER=postgres
DB_PASSWORD=password

REDIS_HOST=localhost
REDIS_PORT=6379

JWT_SECRET=your-secret-key

ACCESS_TOKEN_TTL=15m
REFRESH_TOKEN_TTL=720h
```

---

# Run Application

Development

```bash
go run cmd/api/main.go
```

or

```bash
air
```

---

# Database Migration

Create migration

```bash
migrate create -ext sql -dir migrations create_users_table
```

Run migration

```bash
migrate -path migrations \
-database "$DATABASE_URL" \
up
```

Rollback

```bash
migrate -path migrations \
-database "$DATABASE_URL" \
down 1
```

---

# Swagger

Generate documentation

```bash
swag init -g cmd/api/main.go
```

Swagger UI

```text
http://localhost:8080/swagger/index.html
```

---

# Queue Worker

Run Asynq Worker

```bash
go run cmd/worker/main.go
```

---

# Development Rules

Before creating new code:

- Read `docs/AGENTS.md`
- Read `docs/ARCHITECTURE.md`
- Read `docs/DATABASE.md`

Every implementation must follow these documents.

---

# Coding Principles

The project follows:

- SOLID
- Clean Architecture
- Clean Code
- DRY
- KISS
- YAGNI

---

# Module Rules

Every feature must be implemented as a business module.

Example

```text
user/

role/

permission/

auth/
```

Do not create global business logic packages.

Communication between modules must happen through Service Interfaces.

---

# Database Rules

- PostgreSQL only
- BIGINT Identity Primary Key
- Snake Case
- TIMESTAMPTZ
- Soft Delete
- Audit Columns
- Partial Unique Index
- Migration First

Read `docs/DATABASE.md` before modifying the schema.

---

# API Rules

- RESTful API
- JSON Response
- JWT Authentication
- RBAC Authorization
- Request Validation
- Swagger Documentation

---

# Logging

Logging uses Zerolog.

Every request must include:

- Request ID
- User ID (Authenticated Request)
- Method
- Path
- Status Code
- Duration
- Client IP

Sensitive information must never be logged.

---

# Security

- Passwords are hashed using bcrypt.
- Refresh Tokens are hashed before storage.
- JWT is used for authentication.
- Input validation is mandatory.
- SQL Injection protection relies on GORM parameterized queries.

---

# Testing

Business logic should be testable independently from HTTP.

Prefer unit testing for:

- Service
- Repository
- Middleware

---

# Future Roadmap

- Organization Module
- Donation Module
- Program Module
- Volunteer Module
- Event Module
- Notification Module
- File Storage
- Payment Integration
- Audit Log
- Microservices Migration

---

# Contribution

Before submitting code:

- Follow `docs/AGENTS.md`.
- Follow `docs/ARCHITECTURE.md`.
- Follow `docs/DATABASE.md`.
- Keep code modular.
- Do not introduce breaking changes without discussion.

Consistency is more important than personal coding style.

---

# License

Private Project

Copyright © Alfatihah Center.

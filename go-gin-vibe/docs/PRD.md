# PRD.md

# Alfatihah Center

Version: 1.0

Status: MVP

---

# 1. Overview

## Background

Alfatihah Center adalah sebuah platform yang digunakan untuk mengelola berbagai aktivitas organisasi secara terpusat.

Versi pertama (MVP) berfokus pada pembangunan fondasi sistem, yaitu Authentication, User Management, dan Role Based Access Control (RBAC). Seluruh module berikutnya akan dibangun di atas fondasi ini.

Backend akan menyediakan REST API yang nantinya dapat digunakan oleh Web, Mobile, maupun aplikasi pihak ketiga.

---

# 2. Objectives

Tujuan MVP adalah menyediakan backend yang memiliki kemampuan:

- Authentication menggunakan JWT
- Refresh Token
- Forgot Password
- User Management
- Role Management
- Permission Management
- Role Permission Management
- User Role Management
- REST API
- Swagger Documentation
- Background Job

---

# 3. Out of Scope

Fitur berikut belum termasuk dalam MVP.

- Donation
- Program Management
- Volunteer Management
- Payment Gateway
- Dashboard Analytics
- Notification Center
- Audit Log
- Email Template Management
- Multi Organization
- File Storage
- Two Factor Authentication

---

# 4. Technology Stack

Backend

- Go

Framework

- Gin

ORM

- GORM

Database

- PostgreSQL

Migration

- golang-migrate

Queue

- hibiken/asynq

Authentication

- JWT
- bcrypt

Validation

- go-playground/validator

Logging

- Zerolog

Configuration

- caarlos0/env

Documentation

- Swaggo

---

# 5. Functional Requirements

## Authentication

System harus menyediakan endpoint:

- Login
- Refresh Token
- Forgot Password
- Reset Password
- Logout
- Get Current User

### Login

Input

- Email
- Password

Output

- Access Token
- Refresh Token
- User Profile

---

### Refresh Token

Refresh token digunakan untuk membuat access token baru.

Refresh token disimpan pada database.

Refresh token harus dapat dicabut (revoke).

---

### Forgot Password

User dapat meminta reset password menggunakan email.

System akan membuat token reset password.

Token memiliki masa berlaku.

Token hanya dapat digunakan satu kali.

---

### Logout

Logout akan mencabut refresh token aktif.

Access token tidak disimpan di database.

---

# 6. User Management

Admin dapat:

- Create User
- Update User
- Delete User
- Restore User
- Activate User
- Deactivate User
- View User Detail
- View User List

User dapat melihat profile miliknya sendiri.

---

# 7. Role Management

Admin dapat:

- Create Role
- Update Role
- Delete Role
- View Role
- List Role

---

# 8. Permission Management

Admin dapat:

- Create Permission
- Update Permission
- Delete Permission
- View Permission
- List Permission

Permission menggunakan format:

```text
user.read
user.create
user.update
user.delete
role.read
role.create
```

---

# 9. Role Permission

Admin dapat:

- Assign Permission ke Role
- Remove Permission dari Role
- List Permission berdasarkan Role

Satu Role memiliki banyak Permission.

Satu Permission dapat dimiliki banyak Role.

---

# 10. User Role

Admin dapat:

- Assign Role ke User
- Remove Role dari User
- List Role User

Satu User dapat memiliki lebih dari satu Role.

---

# 11. Database Design

## users

| Column        | Type      |
| ------------- | --------- |
| id            | bigint    |
| email         | varchar   |
| password      | varchar   |
| is_active     | boolean   |
| last_login_at | timestamp |
| created_at    | timestamp |
| updated_at    | timestamp |
| deleted_at    | timestamp |

---

## user_details

| Column     | Type      |
| ---------- | --------- |
| id         | bigint    |
| user_id    | bigint    |
| full_name  | varchar   |
| phone      | varchar   |
| avatar     | text      |
| birth_date | date      |
| gender     | varchar   |
| address    | text      |
| created_at | timestamp |
| updated_at | timestamp |

Relationship

User has one User Detail.

---

## roles

| Column      | Type    |
| ----------- | ------- |
| id          | bigint  |
| name        | varchar |
| slug        | varchar |
| description | text    |

---

## permissions

| Column      | Type    |
| ----------- | ------- |
| id          | bigint  |
| name        | varchar |
| slug        | varchar |
| description | text    |

---

## role_permissions

| Column        | Type   |
| ------------- | ------ |
| role_id       | bigint |
| permission_id | bigint |

Many-to-Many

---

## user_roles

| Column  | Type   |
| ------- | ------ |
| user_id | bigint |
| role_id | bigint |

Many-to-Many

---

## tokens

| Column     | Type      |
| ---------- | --------- |
| id         | bigint    |
| user_id    | bigint    |
| type       | varchar   |
| token      | varchar   |
| expired_at | timestamp |
| used_at    | timestamp |
| created_at | timestamp |

Supported token type

- refresh_token
- forgot_password

Token disimpan dalam bentuk hash.

Plain token tidak boleh disimpan.

---

# 12. API Convention

Semua endpoint menggunakan prefix

```text
/api/v1
```

REST Convention

| Method | Action         |
| ------ | -------------- |
| GET    | Read           |
| POST   | Create         |
| PUT    | Replace        |
| PATCH  | Partial Update |
| DELETE | Delete         |

---

# 13. Response Format

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

Pagination

```json
{
  "data": [],
  "meta": {
    "current_page": 1,
    "per_page": 10,
    "total": 100,
    "last_page": 5
  }
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

Internal Error

```json
{
  "message": "Internal server error."
}
```

---

# 14. Authentication Flow

Login

↓

Generate Access Token

↓

Generate Refresh Token

↓

Hash Refresh Token

↓

Save Refresh Token

↓

Return Tokens

---

Refresh Token

↓

Validate Token

↓

Check Database

↓

Generate New Access Token

↓

(Optional) Rotate Refresh Token

↓

Return New Tokens

---

Forgot Password

↓

Generate Random Token

↓

Hash Token

↓

Save Token

↓

Send Email

↓

Reset Password

↓

Invalidate Token

---

# 15. Authorization

Authorization menggunakan Role Based Access Control (RBAC).

Flow

User

↓

Role

↓

Permission

Endpoint dapat memiliki lebih dari satu permission.

Middleware akan melakukan pengecekan sebelum request masuk ke handler.

---

# 16. Pagination

Semua endpoint list wajib mendukung:

- page
- per_page
- search
- sort
- order

Response

```json
{
  "data": [],
  "meta": {
    "current_page": 1,
    "per_page": 20,
    "total": 100,
    "last_page": 5
  }
}
```

---

# 17. Logging

System harus mencatat:

- Request ID
- User ID
- Method
- URL
- IP Address
- Response Status
- Latency

Sensitive data tidak boleh dicatat.

---

# 18. Queue

Asynq digunakan untuk background jobs.

MVP menggunakan queue untuk:

- Forgot Password Email

Future

- Email Notification
- Push Notification
- Scheduled Job

---

# 19. Documentation

Semua endpoint harus tersedia pada Swagger UI.

Swagger menjadi dokumentasi resmi API.

---

# 20. Future Roadmap

Module berikutnya akan dibangun setelah MVP selesai.

- Organization
- Donation
- Donation Category
- Program
- Volunteer
- Event
- Gallery
- News
- Settings
- Notification
- Audit Log
- File Upload
- Payment Integration
- Multi Tenancy

---

# 21. Non Functional Requirements

- RESTful API
- Modular Architecture
- Clean Architecture
- High Maintainability
- Production Ready
- Scalable
- Secure
- Easy to Test
- Consistent Response
- Comprehensive Logging
- API Documentation
- Background Job Support

---

# 22. MVP Success Criteria

Project dianggap selesai apabila telah memenuhi seluruh kriteria berikut:

- User dapat login menggunakan email dan password.
- User dapat melakukan refresh access token menggunakan refresh token.
- User dapat melakukan logout sehingga refresh token tidak lagi valid.
- User dapat melakukan forgot password dan reset password.
- Admin dapat melakukan CRUD User.
- Admin dapat melakukan CRUD Role.
- Admin dapat melakukan CRUD Permission.
- Admin dapat mengelola relasi User ↔ Role.
- Admin dapat mengelola relasi Role ↔ Permission.
- Seluruh endpoint terlindungi oleh JWT Authentication.
- Seluruh endpoint dapat dibatasi menggunakan Permission (RBAC).
- Seluruh endpoint terdokumentasi pada Swagger.
- Background job untuk Forgot Password berjalan menggunakan Asynq.
- Seluruh migration database dapat dijalankan menggunakan golang-migrate.

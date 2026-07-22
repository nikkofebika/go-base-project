# API.md

# Alfatihah Center API Convention

Version: 1.0

Protocol: REST API

Content-Type: application/json

Authentication: JWT Bearer Token

---

# 1. Purpose

Dokumen ini menjelaskan standar implementasi REST API yang wajib diikuti oleh seluruh developer maupun AI.

Seluruh endpoint harus memiliki format request dan response yang konsisten.

---

# 2. REST Convention

Gunakan HTTP Method sesuai standar.

| Method | Purpose          |
| ------ | ---------------- |
| GET    | Retrieve data    |
| POST   | Create resource  |
| PUT    | Replace resource |
| PATCH  | Partial update   |
| DELETE | Delete resource  |

---

# 3. URL Convention

Gunakan:

```text
/api/v1
```

Contoh

```text
/api/v1/auth/login

/api/v1/users

/api/v1/users/{id}

/api/v1/roles
```

Gunakan:

- lowercase
- plural noun
- kebab-case bila diperlukan

---

# 4. JSON Convention

Semua endpoint wajib menggunakan

```text
application/json
```

---

# 5. Success Response

Seluruh response berhasil menggunakan format berikut.

```json
{
  "data": {}
}
```

---

# 6. Collection Response

```json
{
  "data": []
}
```

---

# 7. Pagination Response

Endpoint yang menggunakan pagination wajib menggunakan format berikut.

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

# 8. Validation Error

Jika request tidak valid.

HTTP Status

```text
422 Unprocessable Entity
```

Response

```json
{
  "message": "Validation failed.",
  "errors": {
    "email": ["The email field is required."]
  }
}
```

---

# 9. Business Error

Jika business rule gagal.

Contoh

```json
{
  "message": "Email already exists."
}
```

---

# 10. Authentication Error

HTTP

```text
401 Unauthorized
```

Response

```json
{
  "message": "Unauthenticated."
}
```

---

# 11. Authorization Error

HTTP

```text
403 Forbidden
```

Response

```json
{
  "message": "You do not have permission to perform this action."
}
```

---

# 12. Not Found

HTTP

```text
404 Not Found
```

Response

```json
{
  "message": "Data not found."
}
```

---

# 13. Internal Server Error

HTTP

```text
500 Internal Server Error
```

Response

```json
{
  "message": "Internal server error."
}
```

Stack trace tidak boleh dikembalikan ke client.

---

# 14. HTTP Status Convention

| Status | Usage                 |
| ------ | --------------------- |
| 200    | Success               |
| 201    | Created               |
| 204    | No Content            |
| 400    | Bad Request           |
| 401    | Unauthorized          |
| 403    | Forbidden             |
| 404    | Not Found             |
| 409    | Conflict              |
| 422    | Validation Error      |
| 500    | Internal Server Error |

---

# 15. Pagination Query

Gunakan query berikut.

```text
?page=1

?per_page=10
```

Contoh

```text
GET /users?page=1&per_page=20
```

Default

```text
page = 1

per_page = 10
```

Maximum

```text
per_page = 100
```

---

# 16. Sorting

Gunakan

```text
sort
order
```

Contoh

```text
/users?sort=created_at&order=desc
```

Nilai order

```text
asc

desc
```

---

# 17. Filtering

Gunakan query parameter.

Contoh

```text
/users?name=john

/users?email=gmail

/users?is_active=true
```

Hindari nested query parameter yang kompleks.

---

# 18. Search

Gunakan parameter

```text
search
```

Contoh

```text
/users?search=john
```

Search bersifat case-insensitive.

---

# 19. Soft Delete

Secara default endpoint tidak mengembalikan data yang sudah di-soft delete.

---

# 20. Authentication

Gunakan

```http
Authorization: Bearer <access_token>
```

Access Token dikirim melalui Authorization Header.

Refresh Token tidak dikirim melalui query parameter.

---

# 21. Request Validation

Seluruh request wajib divalidasi menggunakan:

- go-playground/validator
- DTO

Handler tidak boleh melakukan validasi manual.

---

# 22. Response DTO

Entity tidak boleh dikembalikan langsung.

Flow

```text
Entity

↓

Response DTO

↓

JSON
```

---

# 23. Resource ID

Gunakan

```text
{id}
```

Contoh

```text
/users/1

/roles/2
```

Jangan menggunakan nama parameter lain seperti

```text
/user/{userId}

/role/{roleId}
```

---

# 24. File Upload

Gunakan

```text
multipart/form-data
```

Response tetap mengikuti response wrapper.

---

# 25. Swagger

Seluruh endpoint wajib memiliki dokumentasi Swagger.

Minimal mencakup:

- Summary
- Description
- Tags
- Parameters
- Request Body
- Response
- Security

---

# 26. API Versioning

Gunakan URL Versioning.

```text
/api/v1
```

Versi baru

```text
/api/v2
```

---

# 27. Idempotency

Method

```text
GET

PUT

DELETE
```

harus bersifat idempotent.

---

# 28. API Principles

Seluruh endpoint wajib:

- Menggunakan RESTful Convention.
- Menggunakan JSON.
- Menggunakan Response Wrapper.
- Menggunakan DTO.
- Menggunakan Validation.
- Menggunakan Authentication.
- Menggunakan Authorization.
- Memiliki dokumentasi Swagger.
- Memiliki response yang konsisten.

Response consistency lebih penting daripada preferensi implementasi masing-masing developer.

# DATABASE.md

# Alfatihah Center Database Convention

Version: 1.0

Database Engine: PostgreSQL 16+

ORM: GORM

Migration: golang-migrate

---

# 1. Purpose

Dokumen ini berisi standar dan konvensi database yang wajib diikuti oleh seluruh developer maupun AI ketika membuat atau mengubah schema database.

Tujuan utama:

- Konsisten
- Mudah dipahami
- Mudah di-maintain
- Performa tinggi
- Mudah dikembangkan
- Siap dipisahkan menjadi Microservices

Dokumen ini **tidak menjelaskan struktur tabel secara detail**, tetapi menjelaskan aturan desain database.

---

# 2. General Principles

- Gunakan normalisasi hingga minimal Third Normal Form (3NF).
- Hindari duplikasi data.
- Gunakan foreign key bila memang diperlukan untuk menjaga integritas data.
- Jangan membuat index yang tidak digunakan.
- Hindari premature optimization.
- Database harus mudah dikembangkan tanpa breaking changes.
- Migration harus bersifat immutable.

---

# 3. Naming Convention

## Table

Gunakan **plural noun**.

✅ Benar

```text
users
user_details
roles
permissions
user_roles
role_permissions
tokens
```

❌ Salah

```text
user
User
tbl_users
mst_users
```

---

## Column

Gunakan snake_case.

✅

```text
full_name

created_at

updated_by_id

deleted_at

last_login_at
```

❌

```text
FullName

fullName

fullname
```

---

## Foreign Key

Gunakan format

```text
{table_singular}_id
```

Contoh

```text
user_id

role_id

permission_id

created_by_id

updated_by_id

deleted_by_id
```

---

## Index

Gunakan format

```text
idx_{table}_{column}
```

Contoh

```text
idx_users_created_at

idx_tokens_user_id

idx_tokens_expired_at
```

Composite Index

```text
idx_tokens_user_id_type
```

---

## Unique Index

Gunakan format

```text
uniq_{table}_{column}
```

Contoh

```text
uniq_users_email

uniq_roles_slug

uniq_permissions_slug
```

---

# 4. Primary Key

Seluruh tabel menggunakan satu primary key.

Nama primary key selalu

```text
id
```

Gunakan

```sql
BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY
```

Jangan menggunakan UUID sebagai primary key.

Alasan:

- Lebih cepat
- Index lebih kecil
- Join lebih cepat
- Lebih efisien untuk PostgreSQL

Jika diperlukan identifier publik, gunakan kolom tambahan.

Contoh

```text
public_id UUID UNIQUE
```

---

# 5. Foreign Key

Foreign key menggunakan tipe

```sql
BIGINT
```

Nama mengikuti format

```text
user_id

role_id

permission_id
```

Gunakan foreign key hanya jika memang diperlukan.

Untuk audit column (`created_by_id`, `updated_by_id`, `deleted_by_id`) **tidak perlu menggunakan foreign key constraint** agar lebih fleksibel ketika aplikasi berkembang menjadi microservice.

---

# 6. Timestamp

Gunakan PostgreSQL

```sql
TIMESTAMPTZ
```

Semua waktu disimpan dalam UTC.

Timezone dikonversi pada level aplikasi.

Gunakan nama berikut

```text
created_at

updated_at

deleted_at
```

---

# 7. Audit Columns

Seluruh **master table** wajib memiliki audit column berikut.

```text
created_at

created_by_id

updated_at

updated_by_id

deleted_at

deleted_by_id
```

Tipe

```text
created_at      TIMESTAMPTZ NOT NULL

created_by_id   BIGINT NULL

updated_at      TIMESTAMPTZ NOT NULL

updated_by_id   BIGINT NULL

deleted_at      TIMESTAMPTZ NULL

deleted_by_id   BIGINT NULL
```

Audit column tidak menggunakan foreign key.

Pivot table tidak menggunakan audit column kecuali memang memiliki kebutuhan bisnis.

---

# 8. Boolean

Gunakan prefix

```text
is_

has_

can_
```

Contoh

```text
is_active

is_verified

has_password

can_login
```

---

# 9. Enum

Project **tidak menggunakan PostgreSQL ENUM**.

Gunakan

```sql
VARCHAR
```

Contoh

```text
gender

type

status
```

Enum dikelola pada level aplikasi menggunakan Go Typed Constants.

Contoh

```go
type TokenType string

const (

    TokenTypeRefresh TokenType = "refresh_token"

    TokenTypeForgotPassword TokenType = "forgot_password"

)
```

Validasi enum dilakukan menggunakan Validator dan Business Logic.

---

# 10. Text Convention

Gunakan tipe data sekecil mungkin.

| Type        | Digunakan Untuk                         |
| ----------- | --------------------------------------- |
| VARCHAR(n)  | Data dengan panjang maksimum yang jelas |
| TEXT        | Konten panjang                          |
| BOOLEAN     | True / False                            |
| INTEGER     | Bilangan bulat kecil                    |
| BIGINT      | Identifier                              |
| DATE        | Tanggal                                 |
| TIMESTAMPTZ | Tanggal dan waktu                       |

Gunakan panjang yang sesuai.

Contoh

```text
email          VARCHAR(254)

phone          VARCHAR(30)

slug           VARCHAR(100)

name           VARCHAR(150)

password       VARCHAR(255)
```

Jangan menggunakan

```text
VARCHAR(255)
```

untuk seluruh kolom.

---

# 11. Index Convention

Buat index hanya pada kolom yang sering digunakan.

Contoh

- WHERE
- JOIN
- ORDER BY
- UNIQUE

Contoh

```text
users.email

users.created_at

tokens.user_id

tokens.expired_at

user_roles.role_id

role_permissions.permission_id
```

Hindari index yang tidak diperlukan.

Semakin banyak index akan memperlambat INSERT dan UPDATE.

---

# 12. Unique Index

Gunakan UNIQUE hanya jika memang merupakan business rule.

Contoh

```text
users.email

roles.slug

permissions.slug
```

Untuk tabel yang menggunakan soft delete, gunakan **Partial Unique Index**.

Contoh

```sql
CREATE UNIQUE INDEX uniq_users_email
ON users(email)
WHERE deleted_at IS NULL;
```

---

# 13. Soft Delete

Gunakan soft delete hanya pada tabel master.

Contoh

```text
users

roles

permissions

user_details
```

Jangan menggunakan soft delete pada tabel pivot.

Contoh

```text
user_roles

role_permissions
```

---

# 14. Pivot Table

Pivot table tidak menggunakan kolom

```text
id
```

Gunakan Composite Primary Key.

Contoh

```text
PRIMARY KEY (user_id, role_id)
```

atau

```text
PRIMARY KEY (role_id, permission_id)
```

Pivot table tidak memiliki audit column.

---

# 15. UUID

UUID bukan primary key.

UUID digunakan hanya jika dibutuhkan sebagai identifier publik.

Contoh

```text
public_id UUID UNIQUE
```

Internal application tetap menggunakan BIGINT.

---

# 16. Default Values

Gunakan default value bila memang masuk akal.

Contoh

```sql
is_active BOOLEAN DEFAULT TRUE
```

Jangan membuat default value yang dapat menyembunyikan bug aplikasi.

---

# 17. Migration Rules

Migration menggunakan

```text
golang-migrate
```

Rules

- Migration bersifat immutable.
- Jangan mengubah migration lama.
- Selalu buat migration baru.
- Setiap migration wajib memiliki rollback.
- Migration harus dapat dijalankan pada database kosong.

---

# 18. Query Guidelines

Selalu

- Gunakan SELECT kolom yang diperlukan.
- Hindari SELECT \*.
- Gunakan pagination.
- Gunakan transaction untuk operasi multi-table.
- Gunakan batch insert jika diperlukan.
- Gunakan preload hanya jika memang dibutuhkan.

---

# 19. Performance Guidelines

- Gunakan BIGINT sebagai primary key.
- Gunakan index sesuai query.
- Hindari N+1 Query.
- Hindari over-indexing.
- Hindari join yang tidak diperlukan.
- Simpan data seminimal mungkin.
- Hindari kolom yang tidak digunakan.

---

# 20. Security

- Password harus di-hash menggunakan bcrypt.
- Refresh token harus di-hash sebelum disimpan.
- Jangan pernah menyimpan plain password.
- Jangan pernah menyimpan plain refresh token.
- Jangan menyimpan secret pada database kecuali benar-benar diperlukan.

---

# 21. Module Ownership

Setiap module memiliki ownership terhadap tabelnya sendiri.

Contoh

```text
Auth

users

tokens
```

```text
RBAC

roles

permissions

user_roles

role_permissions
```

Module lain tidak boleh mengubah tabel milik module lain tanpa alasan bisnis yang jelas.

---

# 22. Future Scalability

Database harus dirancang agar mudah dipisahkan menjadi microservice.

Prinsip yang digunakan:

- Setiap module memiliki ownership terhadap tabelnya.
- Hindari dependency yang tidak perlu antar module.
- Hindari circular reference.
- Gunakan audit column tanpa foreign key agar lebih fleksibel saat migrasi.
- Seluruh business rule berada di aplikasi, bukan di database.

---

# 23. Definition of Done

Sebuah desain database dianggap memenuhi standar apabila:

- Menggunakan naming convention yang konsisten.
- Menggunakan BIGINT Identity sebagai primary key.
- Menggunakan snake_case.
- Menggunakan TIMESTAMPTZ.
- Memiliki audit column pada seluruh master table.
- Menggunakan composite primary key pada pivot table.
- Menggunakan partial unique index pada tabel soft delete.
- Tidak menggunakan PostgreSQL ENUM.
- Tidak menggunakan UUID sebagai primary key.
- Memiliki index yang sesuai dengan pola query.
- Seluruh perubahan dilakukan melalui migration.
- Siap dikembangkan menjadi arsitektur microservice tanpa perubahan besar.

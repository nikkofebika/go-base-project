# ERD.md

# Alfatihah Center Entity Relationship Diagram

Version: 1.0

Database: PostgreSQL 16+

---

# 1. Overview

Document ini menjelaskan seluruh tabel, kolom, tipe data, relasi, dan index yang digunakan pada project Alfatihah Center.

---

# 2. Tables

## users

| Column         | Type                            | Constraint                     |
| -------------- | ------------------------------- | ------------------------------ |
| id             | BIGINT                          | GENERATED ALWAYS AS IDENTITY PK |
| email          | VARCHAR(254)                    | NOT NULL UNIQUE                |
| password       | VARCHAR(255)                    | NOT NULL                       |
| is_active      | BOOLEAN                         | NOT NULL DEFAULT TRUE          |
| last_login_at  | TIMESTAMPTZ                     | NULL                           |
| created_at     | TIMESTAMPTZ                     | NOT NULL                       |
| created_by_id  | BIGINT                          | NULL                           |
| updated_at     | TIMESTAMPTZ                     | NOT NULL                       |
| updated_by_id  | BIGINT                          | NULL                           |
| deleted_at     | TIMESTAMPTZ                     | NULL                           |
| deleted_by_id  | BIGINT                          | NULL                           |

Indexes

```sql
CREATE UNIQUE INDEX uniq_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_created_at ON users(created_at);
CREATE INDEX idx_users_is_active ON users(is_active);
```

---

## user_details

| Column     | Type         | Constraint                |
| ---------- | ------------ | ------------------------- |
| id         | BIGINT       | GENERATED ALWAYS AS IDENTITY PK |
| user_id    | BIGINT       | NOT NULL UNIQUE           |
| full_name  | VARCHAR(150) | NOT NULL                  |
| phone      | VARCHAR(30)  | NULL                      |
| avatar     | TEXT         | NULL                      |
| birth_date | DATE         | NULL                      |
| gender     | VARCHAR(10)  | NULL                      |
| address    | TEXT         | NULL                      |
| created_at | TIMESTAMPTZ  | NOT NULL                  |
| updated_at | TIMESTAMPTZ  | NOT NULL                  |

Relationship

- user_details.user_id → users.id (One-to-One)

Indexes

```sql
CREATE INDEX idx_user_details_user_id ON user_details(user_id);
```

---

## roles

| Column      | Type         | Constraint                     |
| ----------- | ------------ | ------------------------------ |
| id          | BIGINT       | GENERATED ALWAYS AS IDENTITY PK |
| name        | VARCHAR(100) | NOT NULL                       |
| slug        | VARCHAR(100) | NOT NULL UNIQUE                |
| description | TEXT         | NULL                           |
| created_at  | TIMESTAMPTZ  | NOT NULL                       |
| created_by_id | BIGINT    | NULL                           |
| updated_at  | TIMESTAMPTZ  | NOT NULL                       |
| updated_by_id | BIGINT    | NULL                           |
| deleted_at  | TIMESTAMPTZ  | NULL                           |
| deleted_by_id | BIGINT    | NULL                           |

Indexes

```sql
CREATE UNIQUE INDEX uniq_roles_slug ON roles(slug) WHERE deleted_at IS NULL;
CREATE INDEX idx_roles_created_at ON roles(created_at);
```

---

## permissions

| Column      | Type         | Constraint                     |
| ----------- | ------------ | ------------------------------ |
| id          | BIGINT       | GENERATED ALWAYS AS IDENTITY PK |
| name        | VARCHAR(100) | NOT NULL                       |
| slug        | VARCHAR(100) | NOT NULL UNIQUE                |
| description | TEXT         | NULL                           |
| created_at  | TIMESTAMPTZ  | NOT NULL                       |
| created_by_id | BIGINT    | NULL                           |
| updated_at  | TIMESTAMPTZ  | NOT NULL                       |
| updated_by_id | BIGINT    | NULL                           |
| deleted_at  | TIMESTAMPTZ  | NULL                           |
| deleted_by_id | BIGINT    | NULL                           |

Indexes

```sql
CREATE UNIQUE INDEX uniq_permissions_slug ON permissions(slug) WHERE deleted_at IS NULL;
CREATE INDEX idx_permissions_created_at ON permissions(created_at);
```

---

## role_permissions

| Column        | Type   | Constraint             |
| ------------- | ------ | ---------------------- |
| role_id       | BIGINT | NOT NULL               |
| permission_id | BIGINT | NOT NULL               |

Primary Key

```sql
PRIMARY KEY (role_id, permission_id)
```

Relationship

- role_permissions.role_id → roles.id
- role_permissions.permission_id → permissions.id

Indexes

```sql
CREATE INDEX idx_role_permissions_role_id ON role_permissions(role_id);
CREATE INDEX idx_role_permissions_permission_id ON role_permissions(permission_id);
```

---

## user_roles

| Column  | Type   | Constraint             |
| ------- | ------ | ---------------------- |
| user_id | BIGINT | NOT NULL               |
| role_id | BIGINT | NOT NULL               |

Primary Key

```sql
PRIMARY KEY (user_id, role_id)
```

Relationship

- user_roles.user_id → users.id
- user_roles.role_id → roles.id

Indexes

```sql
CREATE INDEX idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX idx_user_roles_role_id ON user_roles(role_id);
```

---

## tokens

| Column     | Type         | Constraint                     |
| ---------- | ------------ | ------------------------------ |
| id         | BIGINT       | GENERATED ALWAYS AS IDENTITY PK |
| user_id    | BIGINT       | NOT NULL                       |
| type       | VARCHAR(20)  | NOT NULL                       |
| token      | VARCHAR(255) | NOT NULL                       |
| expired_at | TIMESTAMPTZ  | NOT NULL                       |
| used_at    | TIMESTAMPTZ  | NULL                           |
| created_at | TIMESTAMPTZ  | NOT NULL                       |

Relationship

- tokens.user_id → users.id

Indexes

```sql
CREATE INDEX idx_tokens_user_id ON tokens(user_id);
CREATE INDEX idx_tokens_expired_at ON tokens(expired_at);
CREATE INDEX idx_tokens_user_id_type ON tokens(user_id, type);
```

Supported Token Types

- `refresh_token`
- `forgot_password`

Token harus di-hash sebelum disimpan. Plain token tidak boleh disimpan.

---

# 3. Relationships

```text
users (1) ──── (1) user_details

users (M) ──── (M) roles          via user_roles

roles (M) ──── (M) permissions    via role_permissions

users (1) ──── (M) tokens
```

---

# 4. Module Ownership

| Module    | Tables                              |
| --------- | ----------------------------------- |
| Auth      | users, tokens                       |
| RBAC      | roles, permissions, user_roles, role_permissions |
| User      | users, user_details                 |

---

# 5. Soft Delete

Soft delete hanya digunakan pada master tables:

- users
- roles
- permissions

Tidak menggunakan soft delete pada pivot tables:

- user_roles
- role_permissions

---

# 6. Audit Columns

Master tables memiliki audit columns:

```text
created_at
created_by_id
updated_at
updated_by_id
deleted_at
deleted_by_id
```

Audit columns tidak menggunakan foreign key constraint.

Pivot tables tidak memiliki audit columns.

---

# 7. Naming Convention

| Object     | Format                          |
| ---------- | ------------------------------- |
| Table      | plural noun (snake_case)        |
| Column     | snake_case                      |
| Index      | idx_{table}_{column}            |
| Unique     | uniq_{table}_{column}           |
| Foreign Key | {table_singular}_id            |
| Primary Key | id                             |

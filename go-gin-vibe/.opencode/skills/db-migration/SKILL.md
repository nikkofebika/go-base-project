---
name: db-migration
description: Use when creating, modifying, or rolling back database migrations. Covers golang-migrate workflow, naming conventions, and PostgreSQL best practices.
---

# Database Migration Skill

## Tool

golang-migrate CLI

## Commands

### Create new migration

```bash
migrate create -ext sql -dir migrations -seq <name>
```

### Run migrations

```bash
migrate -path migrations -database "$DATABASE_URL" up
```

### Rollback

```bash
migrate -path migrations -database "$DATABASE_URL" down 1
```

## Naming Convention

```
migrations/
├── 000001_create_users_table.up.sql
├── 000001_create_users_table.down.sql
├── 000002_create_roles_table.up.sql
└── 000002_create_roles_table.down.sql
```

## Rules

1. Migration is immutable - never edit executed migration
2. Always create new migration for changes
3. Every migration must have rollback
4. Use TIMESTAMPTZ for timestamps
5. Use BIGINT GENERATED ALWAYS AS IDENTITY for PK
6. Use snake_case for table and column names
7. Add indexes for WHERE, JOIN, ORDER BY columns
8. Use partial unique index for soft delete tables
9. Pivot tables use composite primary key
10. No audit columns on pivot tables

## PostgreSQL Types

| Use Case | Type |
|----------|------|
| Primary Key | BIGINT GENERATED ALWAYS AS IDENTITY |
| Foreign Key | BIGINT |
| Email | VARCHAR(254) |
| Name | VARCHAR(150) |
| Slug | VARCHAR(100) |
| Password | VARCHAR(255) |
| Boolean | BOOLEAN |
| Timestamp | TIMESTAMPTZ |

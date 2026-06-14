# go-base-project — agent instructions

## Repo layout

Two standalone Go Fiber v3 API apps, each in its own directory. No shared modules, no workspace file, no monorepo tooling.

```
go-fiber-mysql/      # GORM + go-sql-driver/mysql
go-fiber-postgres/   # GORM + pgx (postgres)
```

Both apps are **structurally identical**. Only the DB driver, migration SQL dialect, and DSN format differ.

## Run & dev

Commands below assume `cd go-fiber-{mysql,postgres}` first.

```powershell
# run the app (reads .env)
go run main.go

# seed test data (permissions, roles, admin/user accounts)
go run main.go --seed
```

## Migrations

Uses `golang-migrate/migrate` CLI (not the Go library).

```powershell
# up
migrate -database "<dsn>" -path database/migrations up

# create new migration
migrate create -ext sql -dir database/migrations <name>
```

DSN examples (from `commands.md`):
- MySQL: `mysql://root:password@tcp(localhost:3306)/go_fiber_mysql`
- Postgres: `postgres://postgres:password@localhost:5432/go_fiber_postgres?sslmode=disable`

## Architecture

4-layer, no DI framework — manual constructor injection wired in `internal/router/api_router.go`:

```
Controller → Service (interface) → Repository (interface) → GORM
```

- **Config** (`internal/config/app.go`): `godotenv.Load(".env")` (skipped when `APP_ENV=production`). All values via `os.Getenv` helper wrappers.
- **Validation** (`internal/pkg/validator/`): `go-playground/validator/v10` with custom `exists`, `unique`, `enum`, `length` tag validators.
- **Auth**: JWT (golang-jwt/jwt/v5) + refresh token rotation (stateful via DB). Auth middleware at `internal/app/middlewares/auth_middleware.go`.
- **RBAC**: Permission middleware checks user permissions. Admin bypass. Permissions use dotted convention (`user.read`, `role.create`).
- **Soft deletes**: GORM `DeletedAt` on entities.
- **Error handling**: Custom exception types with factory functions; Fiber global `ErrorHandler` in `internal/app/exception/`.
- **Query building**: `internal/pkg/request/` — `filter[]`, `include[]`, `sort`, `page`/`per_page` DSL.

## Key conventions

- **No tests in this repo** — none exist.
- **No Docker** — no Dockerfile or compose file.
- **No linter/formatter config** — none committed.
- Imports use module names matching directory names: `go-fiber-mysql/...` and `go-fiber-postgres/...`.
- 7 migrations per project: media, users, tokens, roles, permissions, role_permissions, user_roles.

## Gotchas

- The `commands.md` `--seed` example is a comment—use `go run main.go --seed`.
- `.env` is gitignored. Always copy `.env.example` to `.env` and fill in values before running.
- The `database/` directory contains subdirectories for both projects. Ensure you're in the correct project directory before running migration commands.

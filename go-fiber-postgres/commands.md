## MIGRATE

migrate create -ext sql -dir database/migrations create_users_table

migrate -database "postgres://postgres:password@localhost:5432/go_fiber_postgres?sslmode=disable" -path database/migrations up

migrate -database "postgres://postgres:password@localhost:5432/go_fiber_postgres?sslmode=disable" -path database/migrations version

migrate -database "postgres://postgres:password@localhost:5432/go_fiber_postgres?sslmode=disable" -path database/migrations force version

migrate -database "postgres://postgres:passwordword_yang_kuat@tcp(localhost:5432)/go_fiber_postgres?sslmode=disable?sslmode=disable" -path database/migrations up

- ticket-support
- database
- docker
- internal
  - app
    - controllers
    - entities
    - enums
    - exception
    - helpers
    - middlewares
    - repositories
    - requests
    - responses
    - services
  - config
  - pkg
  - router
- storage
  - logs
  - media
  - tmp

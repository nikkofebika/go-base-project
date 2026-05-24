## MIGRATE

migrate create -ext sql -dir database/migrations create_users_table

migrate -database "mysql://root:password@tcp(localhost:3306)/go_fiber_mysql" -path database/migrations up

migrate -database "mysql://root:password@tcp(localhost:3306)/go_fiber_mysql" -path database/migrations version

migrate -database "mysql://root:password@tcp(localhost:3306)/go_fiber_mysql" -path database/migrations force version

migrate -database "mysql://ticket_user:password_yang_kuat@tcp(localhost:3306)/go_fiber_mysql" -path database/migrations up

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

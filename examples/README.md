# HandyGo Web Scaffold

This example is a copy-ready Web project scaffold built on HandyGo, Gin, Gorm, Zap, and Viper.

## Structure

- `main.go`: application entrypoint. It embeds `manifest/config.yaml` into the binary.
- `manifest`: build-time embedded configuration files.
- `internal/bootstrap`: dependency wiring, migration, seed data, server registration.
- `internal/config`: typed application configuration.
- `internal/model`: Gorm models for users, roles, permissions, and join tables.
- `internal/repository`: database access layer.
- `internal/service`: business logic for auth, users, and RBAC.
- `internal/http/api`: public API handlers.
- `internal/http/backend`: backend management handlers.
- `internal/http/middleware`: JWT authentication and permission checks.
- `internal/http/router`: route grouping and permission binding.
- `internal/tasks`: Asynq task definitions and handlers.

## Routes

Public API routes are mounted under `/api/v1`.

- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `GET /api/v1/me`

Backend routes are mounted under `/backend/v1`.

- `POST /backend/v1/auth/login`
- `GET /backend/v1/me`
- `GET /backend/v1/users`
- `POST /backend/v1/users`
- `PUT /backend/v1/users/:id`
- `GET /backend/v1/roles`
- `POST /backend/v1/roles`
- `GET /backend/v1/permissions`

## Default Admin

On first start, the scaffold migrates SQLite tables and creates:

- Email: `admin@example.com`
- Password: `admin123456`

The default admin user is bound to the `admin` role, and the role is bound to all backend permissions.

## Workers

The scaffold demonstrates HandyGo's background infrastructure:

- Local worker pool and local scheduler are enabled by default.
- Asynq client, worker server, and scheduler are configured under `asynq` and disabled by default.
- Set `asynq.enabled` to `true` and provide Redis settings to enable distributed background tasks.

## Run

```bash
cd examples
go run .
```

## Build

```bash
cd examples
go build -o handygo-example .
./handygo-example
```

`manifest/config.yaml` is embedded into the binary at build time. Deployment only needs the built executable; the YAML file does not need to be shipped with it.

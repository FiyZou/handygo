# HandyGo Web App

This project was created with HandyGo. It is built on HandyGo, Gin, Gorm, Zap, and Viper.

## Quick Start

```bash
make install-tools
go mod tidy
make generate
make dev
```

`make dev` uses `manifest/config.local.yaml`, starts the server in debug mode, and creates a local SQLite database on first run.

## Structure

- `main.go`: application entrypoint. It embeds `manifest/config.yaml` into the binary.
- `manifest/config.yaml`: embedded build-time configuration.
- `manifest/config.local.yaml`: editable local development configuration.
- `manifest/gen.yaml`: database-to-model generation configuration.
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
make dev
```

Use `APP_CONFIG=/path/to/config.yaml go run .` for production-like runs. The embedded `manifest/config.yaml` intentionally leaves `auth.jwtSecret` empty, so production secrets must be supplied by deployment configuration.

## Generate Models

```bash
make generate
```

The command reads `manifest/gen.yaml` and writes generated models to `internal/model`. It does not generate repository/query code. Generated `*.gen.go` files are overwritten on each run; keep hand-written model helpers in separate non-generated files.

You can run the CLI directly:

```bash
handygo gen model -c manifest/gen.yaml
```

## Build

```bash
make build
./handygo-example
```

`manifest/config.yaml` is embedded into the binary at build time. Deployment only needs the built executable when required settings such as `auth.jwtSecret` are supplied by an external config path through `APP_CONFIG`.

## Common Commands

```bash
make help
make test
make smoke
make clean
```

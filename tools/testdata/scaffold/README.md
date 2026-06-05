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

## AI Collaboration

The scaffold includes an automatic collaboration workspace for CLI agents and human review. Users describe the goal; agents maintain the collaboration documents and role handoffs.

- `AGENTS.md`: global engineering and handoff rules
- `.codex/agents/collaboration-runner.md`: default entrypoint for user goals
- `docs/ai-collaboration.md`: bilingual usage guide for the closed-loop workflow
- `docs/handoff.md`: current role handoff summary
- `docs/decision-log.md`: architecture and technical decisions
- `docs/tasks.md`: active and upcoming work
- `docs/product/PRD.md`: product requirements
- `docs/tech/ARCHITECTURE.md`: system design notes
- `docs/review/`: review artifacts
- `docs/qa/`: QA artifacts
- `.codex/agents/`: project-local Codex agent notes and starter specialist agents

## User Goal Flow

For a fresh project, tell the agent the desired outcome:

```text
实现注册登录
```

The collaboration runner handles the rest:

1. PM scopes the goal and updates `docs/product/PRD.md`, `docs/tasks.md`, and `docs/handoff.md`.
2. Architect updates `docs/tech/ARCHITECTURE.md`, records decisions, and hands off to Developer.
3. Developer implements the scoped change, tests it, and updates task and handoff state.
4. Reviewer writes `docs/review/quality-report.md`, updates QA notes when needed, and closes or continues the handoff.

## Structure

- `main.go`: application entrypoint. It embeds `manifest/config.yaml` into the binary.
- `api/.../v1`: request and response DTOs. This is the transport contract layer, not the handler/controller layer.
- `manifest/config.yaml`: embedded build-time configuration.
- `manifest/config.local.yaml`: editable local development configuration.
- `manifest/gen.yaml`: database-to-model generation configuration.
- `docs`: collaboration memory, decisions, handoffs, and QA/review notes.
- `docs/ai-collaboration.md`: English and Chinese guide for user-facing agent collaboration.
- `.codex/agents`: project-local agent instructions and conventions.
- `.codex/agents/collaboration-runner.md`: automatic PM -> Architect -> Developer -> Reviewer coordinator.
- `.codex/agents/gorm-expert.md`: data access and transaction guardrails.
- `.codex/agents/gin-architect.md`: HTTP layering and API boundary guardrails.
- `.codex/agents/api-designer.md`: API contract design guardrails.
- `.codex/agents/quality-gate.md`: quality report checklist and output rules.
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

## Structure Conventions

- `api/.../v1` defines request and response DTOs.
- `internal/http/...` is the current HTTP boundary: handlers, middleware, and routing.
- `manifest/` is the active config directory for both runtime config and model generation. Treat it as the current equivalent of a future `configs/` directory.
- The scaffold may evolve toward `cmd/`, `internal/api`, and `configs/`, but the current generated layout is the source of truth until the CLI and templates are migrated together.
- QA planning defaults live in `docs/qa/test_cases.md`.
- Quality gate output defaults live in `docs/review/quality-report.md`.

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
- Set both `cache.enabled` and `asynq.enabled` to `true` to enable distributed background tasks. Asynq reuses `cache.redis` for its client, server, and scheduler.

## Redis

Redis cache is optional and disabled by default:

```yaml
cache:
  enabled: false
```

Configure `cache.redis` once for both cache usage and Asynq. `cache.enabled` declares that Redis is available; `asynq.enabled` requires `cache.enabled: true`.

## Logging

The default production logger uses JSON output, formatted timestamps, and rotating log files:

```yaml
logger:
  encoding: json
  timeFormat: "2006-01-02 15:04:05"
  fileOutputPath: logs/app.log
  errorFileOutputPath: logs/error.log
  rotation:
    enabled: true
    maxSizeMB: 100
    maxAgeDays: 30
    maxBackups: 10
```

`encoding` supports `json` and `console`. Log rotation uses size, age, backup count, local time, and compression settings. The scaffold does not include external log shipping.

## Database Options

SQLite is the default database:

```yaml
database:
  driver: sqlite
  dsn: file:handygo-example.db?cache=shared
```

`manifest/config.yaml` and `manifest/config.local.yaml` include commented MySQL and PostgreSQL examples next to the default SQLite settings. Uncomment one database driver and DSN pair when switching databases.

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

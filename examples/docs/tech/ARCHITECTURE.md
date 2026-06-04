# Architecture

## Stack

- Go
- Gin
- GORM
- HandyGo shared packages

## Layering

- API/handler layer: request binding, response shaping, auth boundary
- Service layer: business logic, transactions, orchestration
- Repository layer: data access
- Model layer: persistence models

## Current Structure Notes

- `api/.../v1`: request and response DTOs, not controllers
- `internal/http/...`: handlers, routing, middleware, and HTTP boundary
- `internal/service`: business logic
- `internal/repository`: data access
- `internal/model`: GORM models
- `manifest/`: runtime and generation config

## Terminology Mapping

This scaffold intentionally keeps the current working structure instead of renaming everything up front.

| Current Path | Current Meaning | Target Concept |
| --- | --- | --- |
| `api/.../v1` | DTO definitions for request/response payloads | transport contract / API schema |
| `internal/http/api` | public HTTP handlers | API/controller layer |
| `internal/http/backend` | admin/backend HTTP handlers | API/controller layer |
| `internal/http/router` | route registration | router layer |
| `internal/http/middleware` | auth and HTTP middleware | middleware layer |
| `manifest/` | runtime config and model-generation config | `configs/` equivalent for now |

## Structure Guardrails

- Do not treat `api/.../v1` as the place for business logic.
- Do not move handler code into `api/.../v1`.
- Keep business logic in `internal/service`.
- Keep data access in `internal/repository`.
- Keep `manifest/` as the active config directory until the CLI, README, Makefile, and generated scaffold are migrated together.

## Planned Evolution

The long-term target may evolve toward:

```text
cmd/
internal/
  api/
  service/
  repository/
  model/
  middleware/
pkg/
configs/
```

Until that migration is explicitly planned, documented, and verified end to end, the current layout is the source of truth.

## Decisions To Confirm

- Authentication mode
- Tenant isolation strategy
- Background task requirements

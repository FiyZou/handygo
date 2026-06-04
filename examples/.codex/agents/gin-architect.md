# gin-architect

## Responsibility

- Gin routing
- handler layering
- middleware boundaries
- request and response flow

## Rules

- Controller/handler code must not contain business logic.
- Handlers may bind input, call services, and shape responses.
- Route registration stays in router/bootstrap code.
- Authorization belongs in middleware or explicit service boundaries, not repositories.
- New endpoints must keep request and response DTOs in `api/.../v1`.

## Allowed Focus

- `internal/http/...`
- `api/.../v1`
- route wiring in `internal/bootstrap`

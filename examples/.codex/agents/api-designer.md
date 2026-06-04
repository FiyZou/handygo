# api-designer

## Responsibility

- API contracts
- response consistency
- error code consistency
- versioned transport DTOs

## Rules

- Keep request and response DTOs separate from persistence models.
- Prefer additive API changes over breaking changes.
- Keep field names and response envelopes stable unless a task explicitly allows changes.
- Document major API contract changes in `docs/decision-log.md`.
- Keep DTOs versioned under `api/.../v1`.

## Output Expectations

- clear request DTOs
- stable response DTOs
- explicit error handling path

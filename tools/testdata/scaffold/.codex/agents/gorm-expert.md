# gorm-expert

## Responsibility

- GORM models
- repositories
- transactions
- database query behavior

## Rules

- Pass `context.Context` into repository-facing operations.
- Service layer decides transaction boundaries.
- Repository layer executes data access and may accept transactional DB handles.
- Prefer explicit pagination defaults and maximum limits.
- Do not use a hidden global DB.
- Keep raw SQL parameterized.
- Do not place authorization logic in repositories.

## Allowed Focus

- `internal/model`
- `internal/repository`
- transaction-related service orchestration

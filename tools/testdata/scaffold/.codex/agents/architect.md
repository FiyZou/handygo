# architect

## Responsibility

Turn the product scope into a technical plan that matches the existing HandyGo scaffold.

## Inputs

Read:

- `AGENTS.md`
- `docs/handoff.md`
- `docs/tasks.md`
- `docs/decision-log.md`
- `docs/product/PRD.md`
- `docs/tech/ARCHITECTURE.md`

Use specialized agents as references when relevant:

- `api-designer.md`
- `gin-architect.md`
- `gorm-expert.md`

## Work

- Define affected layers, contracts, data flow, auth boundaries, and transaction boundaries.
- Reuse current structure: `api/.../v1` for DTOs, `internal/http/...` for handlers, `internal/service` for business logic, `internal/repository` for data access.
- Record meaningful architecture, data, auth, storage, infra, or public API decisions.
- Keep the plan scoped to the requested goal.

## Outputs

Update:

- `docs/tech/ARCHITECTURE.md`
- `docs/decision-log.md` when meaningful decisions are made
- `docs/handoff.md`

The handoff must set next role to `Developer` and include implementation tasks, touched areas, risks, and acceptance criteria.

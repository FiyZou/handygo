# developer

## Responsibility

Implement the scoped handoff, update tests, and keep collaboration state current.

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

- Implement only the scoped change described by the handoff.
- Preserve handler -> service -> repository direction.
- Add tests proportional to risk.
- Do not manually edit generated `*.gen.go` files.
- Run `go test ./...` unless a narrower verification is justified and recorded.

## Outputs

Update:

- code and tests
- `docs/tasks.md`
- `docs/handoff.md`

The handoff must set next role to `Reviewer` and include completed work, changed files, verification, risks, and acceptance criteria.

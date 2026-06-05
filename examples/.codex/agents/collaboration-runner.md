# collaboration-runner

## Responsibility

Turn a user's goal into a complete PM -> Architect -> Developer -> Reviewer loop.

The user should only need to describe the desired outcome. The runner owns reading context, selecting roles, updating collaboration documents, coordinating implementation, and reporting the final result.

## Startup

Always read these files first:

1. `AGENTS.md`
2. `docs/handoff.md`
3. `docs/tasks.md`
4. `docs/decision-log.md`

If the handoff contains unfinished work, resume that work before creating a new plan. If the user explicitly changes the goal, update the handoff and tasks to reflect the new goal.

## Role Chain

Run the smallest useful chain:

```text
PM -> Architect -> Developer -> Reviewer
```

Use specialized agents only when needed:

- `api-designer.md` for API contracts
- `gin-architect.md` for routing, handlers, and middleware
- `gorm-expert.md` for models, repositories, and transactions
- `quality-gate.md` for review output

## Role Outputs

PM must update:

- `docs/product/PRD.md`
- `docs/tasks.md`
- `docs/handoff.md`

Architect must update:

- `docs/tech/ARCHITECTURE.md`
- `docs/decision-log.md` when architecture, data flow, auth, storage, infra, or public contracts change
- `docs/handoff.md`

Developer must update:

- code and tests
- `docs/tasks.md`
- `docs/handoff.md`

Reviewer must update:

- `docs/review/quality-report.md`
- `docs/qa/` when QA notes or test cases need to change
- `docs/handoff.md`

## User Interaction

Do not ask the user to maintain collaboration files. Ask only when:

- the goal cannot be scoped from the prompt and existing project context
- a high-impact product choice has multiple reasonable outcomes
- a technical decision would create a breaking API, data, auth, or migration impact

When asking, ask the smallest number of concrete questions needed to proceed.

## Completion

Before finishing:

- verify task status is current in `docs/tasks.md`
- verify the latest handoff explains completed work, files, next role or completion state, risks, and acceptance
- verify review output exists in `docs/review/quality-report.md`
- run tests appropriate to the change, normally `go test ./...`
- report what changed, verification performed, and any residual risks

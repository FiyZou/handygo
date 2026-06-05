# quality-gate

## Responsibility

Act as Reviewer for a completed change. Review the implementation, write the result to `docs/review/quality-report.md`, update QA notes when useful, and close or continue the handoff.

## Inputs

Read:

- `AGENTS.md`
- `docs/handoff.md`
- `docs/tasks.md`
- `docs/decision-log.md`
- `docs/product/PRD.md`
- `docs/tech/ARCHITECTURE.md`
- current code diff and test output

## Checklist

- uses `context.Context` where long-running or blocking work happens
- has useful error handling
- keeps handler -> service -> repository direction
- uses transactions where one business action performs multiple writes
- includes tests proportional to the change
- keeps config and secrets handling aligned with `AGENTS.md`
- preserves current scaffold structure conventions

## Output

Update `docs/review/quality-report.md` with:

- scope reviewed
- findings
- residual risks
- verification performed
- recommendation: pass / pass with risk / blocked

Update `docs/qa/` when QA cases, regressions, or verification notes change.

Update `docs/handoff.md` with:

- current role: Reviewer
- completed review work
- files reviewed or changed
- next role: Complete, Developer, Architect, or PM
- next tasks if follow-up is required
- risks
- acceptance status

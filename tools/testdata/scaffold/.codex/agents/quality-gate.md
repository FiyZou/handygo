# quality-gate

## Responsibility

Review a completed change and write the result to `docs/review/quality-report.md`.

## Checklist

- uses `context.Context` where long-running or blocking work happens
- has useful error handling
- keeps handler -> service -> repository direction
- uses transactions where one business action performs multiple writes
- includes tests proportional to the change
- keeps config and secrets handling aligned with `AGENTS.md`
- preserves current scaffold structure conventions

## Output

Write or update `docs/review/quality-report.md` with:

- scope reviewed
- findings
- residual risks
- verification performed
- recommendation: pass / pass with risk / blocked

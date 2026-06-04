# AI Development Rules

This project is expected to be developed with Codex or other CLI coding agents. Follow these rules for every generated change.

## 1. Reuse Before Creating

- Search the codebase before adding new code.
- Prefer existing functions, services, repositories, middleware, config structs, response helpers, and validation helpers.
- Do not create duplicate helpers with different names for the same behavior.
- If existing code is close but not enough, extend it narrowly instead of creating a parallel implementation.
- Before adding a new abstraction, confirm it removes real duplication or clarifies a repeated workflow.

## 2. Keep Data Flow Directional

- Data must flow from outer layers to inner layers:
  - HTTP handler -> service -> repository -> database/model
  - bootstrap -> service wiring -> route registration
- Inner layers must not import or depend on outer layers.
- Repository code must not call HTTP handlers or services.
- Service code must not depend on HTTP request/response objects unless those types are explicitly API DTOs.
- Handler code should translate request/response concerns and delegate business behavior to services.
- Shared model and config types may be imported inward as needed, but avoid circular ownership.

## 3. Log For Debugging

- Log important lifecycle events, failures, and boundary decisions.
- Use structured logs where possible.
- Include enough context to locate bugs quickly:
  - operation name
  - resource id or key
  - user id when available
  - error value
  - external dependency name
- Do not log passwords, tokens, secrets, or raw credentials.
- Use log levels consistently:
  - debug: detailed development diagnostics
  - info: successful lifecycle and business milestones
  - warn: recoverable abnormal conditions
  - error: failed operations that require attention

## 4. Comments Must Help Review

- Add comments for non-obvious business rules, edge cases, concurrency, transactions, and security-sensitive code.
- Do not add comments that merely repeat the code.
- Generated code should be marked as generated when appropriate.
- Public functions and exported types should have clear names; add comments when the intent is not obvious from the name.
- Keep comments short and review-oriented.

## 5. Tests Are Required

- Every functional change must include tests unless there is a clear reason it cannot be tested.
- Cover normal behavior, failure behavior, and boundary cases.
- For validation logic, test empty values, invalid values, and valid edge values.
- For pagination, test page/size defaults, upper limits, and empty results.
- For auth and permission changes, test unauthorized, forbidden, and allowed cases.
- For repository/database changes, test query filters, not-found behavior, and transaction boundaries where practical.
- Run `go test ./...` before considering the change complete.

## 6. Configuration And Secrets

- Prefer typed config structs over ad hoc environment reads.
- Keep local development config in `manifest/config.local.yaml`.
- Do not commit real secrets.
- Support explicit config paths through `APP_CONFIG` for runtime and `manifest/gen.yaml` for model generation.

## 7. Generated Files

- Do not manually edit `*.gen.go` files.
- Put hand-written model constants, helpers, or methods in separate non-generated files.
- Regenerate models with:

```bash
make generate
```

## 8. Change Discipline

- Keep changes scoped to the requested behavior.
- Do not refactor unrelated code while implementing a feature.
- Preserve existing public APIs unless the task explicitly requires changing them.
- If a change crosses layers, update all affected tests and documentation.

## 9. Verification Checklist

Before finishing a task, verify:

- Existing helpers were reused where possible.
- Data flow still follows handler -> service -> repository.
- Logs include useful context and no secrets.
- Comments explain only non-obvious decisions.
- Tests cover normal, failure, and boundary behavior.
- `go test ./...` passes.

## 10. Security Boundaries

- Validate every external input, including body, query, path, header, and task payloads.
- Do not trust client-provided user ids, roles, permissions, prices, ownership, status, or audit fields.
- Authorization checks must live at explicit middleware, handler, or service boundaries.
- Repository code must never make authorization decisions.
- Never log tokens, passwords, secrets, raw credentials, or full `Authorization` headers.
- Authentication, authorization, password, token, and signature changes require negative tests.

## 11. Error Handling

- Do not ignore errors.
- Repository functions should return storage errors without converting them into HTTP concepts.
- Services should translate storage errors into business meaning where needed.
- Handlers should translate service errors into stable HTTP responses and error codes.
- Do not return database errors, stack traces, or internal dependency details to clients.
- Log internal details on the server side when they are needed for debugging.
- Every external dependency failure must have a clear failure path.

## 12. Transactions And Consistency

- Multiple database writes that belong to one business action must run in one transaction.
- Service layer decides transaction boundaries.
- Repository layer performs data operations and must accept transaction handles where needed.
- Do not perform irreversible external side effects inside a database transaction.
- Keep transactions short; do not put slow network calls or long-running work inside them.
- Test transaction rollback behavior when a later step fails.

## 13. Concurrency And Resource Management

- Every goroutine must have an exit path.
- Prefer passing `context.Context` into long-running or blocking operations.
- Do not create unbounded goroutines, queues, retries, or worker pools.
- External calls, database work, Redis operations, and task processing must support timeout or cancellation where practical.
- Close files, response bodies, tickers, timers, database rows, and other resources.
- Retries must have a maximum attempt count and backoff strategy.

## 14. Production Configuration

- Production must not rely on `manifest/config.local.yaml`.
- Development defaults are allowed only in local config.
- Production secrets must be provided explicitly by deployment configuration or environment.
- New config values must be added to typed config structs.
- Config required for startup must be validated during bootstrap.
- Do not hardcode database, Redis, queue, third-party API, or storage endpoints in business code.
- Do not commit real credentials.

## 15. Database And Model Safety

- Do not manually edit `*.gen.go`.
- List endpoints must enforce pagination and maximum page size.
- Delete behavior must be explicit: soft delete or hard delete.
- Unique constraints, status fields, and relationship constraints must match business rules.
- Repository code must handle not-found, duplicate, empty result, and filter boundary cases.
- Raw SQL must be parameterized; string concatenation for SQL is not allowed.

## 16. API Compatibility

- Do not change existing response field names, status codes, or error codes without explicit approval.
- Keep request and response DTOs separate from database models.
- Never expose password hashes, secrets, tokens, or internal-only fields in API responses.
- Additive response fields are preferred over breaking changes.
- Breaking API changes require README or changelog updates.

## 17. Scaffold And CLI Verification

- Changes to the scaffold must be verified with `handygo new`.
- Changes to model generation must be verified with `make generate`.
- Changes to Makefile commands must be verified from the generated project directory.
- Changes to Cobra commands must verify `--help`, success path, and failure path.
- Do not leave generated binaries, local databases, or temporary files in the repository.

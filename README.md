# HandyGo

HandyGo is a reusable Go web foundation package. It wraps popular production libraries with a small set of stable APIs so new web services can share the same app lifecycle, configuration, logging, HTTP server, database, cache, health, middleware, and response conventions.

## Features

- Gin based HTTP server and middleware.
- Gorm based database initialization and transaction helpers.
- Zap based structured logging.
- Viper based YAML/TOML configuration loading with environment override support.
- go-redis based Redis client wrapper.
- Shared app lifecycle management and graceful shutdown.
- Health checks, unified JSON responses, and example web service wiring.
- Safe goroutine helpers, worker pools, local schedulers, and Asynq based background queues.

## Scaffold CLI

HandyGo also ships a project scaffold CLI:

```bash
go install github.com/FiyZou/handygo/cmd/handygo@latest
handygo new myapp --module github.com/you/myapp
```

`handygo new` runs `go mod tidy` in the created project by default. Use `--skip-tidy` for offline or scripted runs.

The generated project includes both the runnable web scaffold and a default collaboration workspace:

- `AGENTS.md`: global engineering and handoff rules
- `.codex/agents/collaboration-runner.md`: default agent entrypoint for user goals
- `docs/ai-collaboration.md`: bilingual user guide for the closed-loop agent workflow
- `docs/collaboration-config.yaml`: optional frontend workflow and style skill settings
- `docs/handoff.md`: current role handoff state
- `docs/tasks.md`: backlog, in-progress, and done items
- `docs/decision-log.md`: architectural and public-contract decisions
- `docs/product/PRD.md`: product requirements
- `docs/tech/ARCHITECTURE.md`: system design notes
- `docs/review/`: review outputs
- `docs/qa/`: QA planning and verification notes
- `.codex/agents/`: project-local specialist agent prompts

This is the stage-three scaffold: users describe the desired outcome, and the collaboration runner advances PM, Architect, Developer, and Reviewer roles while maintaining PRD, architecture, tasks, decisions, handoff, review, and QA notes.

Frontend agent workflow is optional and disabled by default. Enable it during `handygo new` with `--frontend --frontend-style-skill <skill>` or later in `docs/collaboration-config.yaml`.

## Quick Start

```go
import (
    handyconfig "github.com/FiyZou/handygo/config"
    handylogger "github.com/FiyZou/handygo/logger"
    handyresponse "github.com/FiyZou/handygo/response"
    handyserver "github.com/FiyZou/handygo/server"
    "github.com/gin-gonic/gin"
)

cfg := struct {
    Server handyserver.Config `mapstructure:"server"`
    Logger handylogger.Config `mapstructure:"logger"`
}{}

if err := handyconfig.Load("config.yaml", &cfg); err != nil {
    panic(err)
}

log, _ := handylogger.New(cfg.Logger)
srv := handyserver.New(cfg.Server, log)
srv.Register(func(r *gin.Engine) {
    r.GET("/ping", func(c *gin.Context) {
        handyresponse.OK(c, gin.H{"message": "pong"})
    })
})
```

See `examples` for a complete composition example.

See [examples/README.md](examples/README.md) for the generated project layout and workflow details.
See [examples/docs/ai-collaboration.md](examples/docs/ai-collaboration.md) for the bilingual collaboration guide.

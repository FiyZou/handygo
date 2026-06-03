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

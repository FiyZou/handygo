package bootstrap

import (
	exampleconfig "github.com/FiyZou/handygo/examples/internal/config"
	"github.com/FiyZou/handygo/examples/internal/http/router"
	"github.com/FiyZou/handygo/health"
	handymiddleware "github.com/FiyZou/handygo/middleware"
	"github.com/FiyZou/handygo/server"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func newHTTPServer(cfg exampleconfig.AppConfig, zapLogger *zap.Logger, healthChecker *health.Health, services services) *server.Server {
	httpServer := server.New(cfg.Server, zapLogger)
	httpServer.Use(
		handymiddleware.RequestID(),
		handymiddleware.Recover(zapLogger),
		handymiddleware.AccessLog(zapLogger),
		handymiddleware.CORS(handymiddleware.CORSConfig{}),
	)
	httpServer.Register(func(engine *gin.Engine) {
		router.Register(engine, router.Dependencies{
			Health: healthChecker,
			Auth:   services.auth,
			Users:  services.user,
			RBAC:   services.rbac,
		})
	})
	return httpServer
}

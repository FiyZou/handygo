package router

import (
	apihttp "github.com/FiyZou/handygo/examples/internal/http/api"
	backendhttp "github.com/FiyZou/handygo/examples/internal/http/backend"
	authmw "github.com/FiyZou/handygo/examples/internal/http/middleware"
	"github.com/FiyZou/handygo/examples/internal/service"
	"github.com/FiyZou/handygo/health"
	"github.com/gin-gonic/gin"
)

type Dependencies struct {
	Health *health.Health
	Auth   *service.AuthService
	Users  *service.UserService
	RBAC   *service.RBACService
}

func Register(engine *gin.Engine, deps Dependencies) {
	engine.GET("/health", deps.Health.Handler())

	apiAuth := apihttp.NewAuthHandler(deps.Auth)
	apiProfile := apihttp.NewProfileHandler(deps.Auth)

	api := engine.Group("/api/v1")
	api.POST("/auth/register", apiAuth.Register)
	api.POST("/auth/login", apiAuth.Login)
	api.GET("/me", authmw.Auth(deps.Auth), apiProfile.Me)

	backendAuth := backendhttp.NewAuthHandler(deps.Auth)
	backendUsers := backendhttp.NewUserHandler(deps.Users)
	backendRBAC := backendhttp.NewRBACHandler(deps.RBAC)

	backend := engine.Group("/backend/v1")
	backend.POST("/auth/login", backendAuth.Login)

	protected := backend.Group("")
	protected.Use(authmw.Auth(deps.Auth), authmw.RequirePermission("backend:access"))
	protected.GET("/me", backendAuth.Me)
	protected.GET("/users", authmw.RequirePermission("user:list"), backendUsers.List)
	protected.POST("/users", authmw.RequirePermission("user:create"), backendUsers.Create)
	protected.PUT("/users/:id", authmw.RequirePermission("user:update"), backendUsers.Update)
	protected.GET("/roles", authmw.RequirePermission("role:list"), backendRBAC.Roles)
	protected.POST("/roles", authmw.RequirePermission("role:create"), backendRBAC.CreateRole)
	protected.GET("/permissions", authmw.RequirePermission("permission:list"), backendRBAC.Permissions)
}

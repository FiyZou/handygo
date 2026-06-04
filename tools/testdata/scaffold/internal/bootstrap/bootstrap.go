package bootstrap

import (
	"fmt"
	"time"

	"github.com/FiyZou/handygo/app"
	"github.com/FiyZou/handygo/database"
	exampleconfig "github.com/FiyZou/handygo/examples/internal/config"
	"github.com/FiyZou/handygo/examples/internal/repository"
	"github.com/FiyZou/handygo/examples/internal/service"
	"github.com/FiyZou/handygo/health"
	handylogger "github.com/FiyZou/handygo/logger"
	"go.uber.org/zap"
)

type services struct {
	auth *service.AuthService
	user *service.UserService
	rbac *service.RBACService
}

func New(cfg exampleconfig.AppConfig) (*app.App, error) {
	if err := cfg.NormalizeAndValidate(); err != nil {
		return nil, err
	}

	zapLogger, err := newLogger(cfg)
	if err != nil {
		return nil, err
	}

	db, err := newDatabase(cfg)
	if err != nil {
		return nil, err
	}

	services := newServices(cfg, db)
	healthChecker := newHealth(db)
	httpServer := newHTTPServer(cfg, zapLogger, healthChecker, services)
	localPool, localScheduler, err := newLocalWorkers(cfg, zapLogger)
	if err != nil {
		return nil, err
	}

	application := app.New(cfg.App.Name, app.WithLogger(zapLogger.Sugar()))
	application.Register(db, localPool, localScheduler, httpServer)
	if cfg.Asynq.Enabled {
		if err := registerAsynq(application, cfg, zapLogger); err != nil {
			return nil, err
		}
	}
	return application, nil
}

func newLogger(cfg exampleconfig.AppConfig) (*zap.Logger, error) {
	zapLogger, err := handylogger.New(cfg.Logger)
	if err != nil {
		return nil, fmt.Errorf("new logger: %w", err)
	}
	handylogger.SetDefault(zapLogger)
	return zapLogger, nil
}

func newDatabase(cfg exampleconfig.AppConfig) (*database.Database, error) {
	db, err := database.New(cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("new database: %w", err)
	}
	if err := migrateAndSeed(db.DB(), cfg.Auth.PasswordCost); err != nil {
		return nil, fmt.Errorf("migrate and seed: %w", err)
	}
	return db, nil
}

func newServices(cfg exampleconfig.AppConfig, db *database.Database) services {
	userRepo := repository.NewUserRepository(db.DB())
	rbacRepo := repository.NewRBACRepository(db.DB())
	return services{
		auth: service.NewAuthService(cfg.Auth, userRepo, rbacRepo),
		user: service.NewUserService(cfg.Auth.PasswordCost, userRepo, rbacRepo),
		rbac: service.NewRBACService(rbacRepo),
	}
}

func newHealth(db *database.Database) *health.Health {
	healthChecker := health.New(defaultHealthTimeout)
	healthChecker.Register(health.NewCheck("database", db.Start))
	return healthChecker
}

const defaultHealthTimeout = 2 * time.Second

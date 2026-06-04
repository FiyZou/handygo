package config

import (
	"errors"
	"strings"
	"time"

	"github.com/FiyZou/handygo/database"
	handylogger "github.com/FiyZou/handygo/logger"
	"github.com/FiyZou/handygo/queue"
	"github.com/FiyZou/handygo/scheduler"
	"github.com/FiyZou/handygo/server"
	"github.com/FiyZou/handygo/workerpool"
)

type AppConfig struct {
	App      App                `mapstructure:"app"`
	Server   server.Config      `mapstructure:"server"`
	Logger   handylogger.Config `mapstructure:"logger"`
	Database database.Config    `mapstructure:"database"`
	Auth     Auth               `mapstructure:"auth"`
	Worker   Worker             `mapstructure:"worker"`
	Asynq    Asynq              `mapstructure:"asynq"`
}

type App struct {
	Name string `mapstructure:"name"`
}

type Auth struct {
	JWTSecret    string        `mapstructure:"jwtSecret"`
	TokenTTL     time.Duration `mapstructure:"tokenTTL"`
	PasswordCost int           `mapstructure:"passwordCost"`
}

type Worker struct {
	Pool      workerpool.Config `mapstructure:"pool"`
	Scheduler scheduler.Config  `mapstructure:"scheduler"`
}

type Asynq struct {
	Enabled   bool                  `mapstructure:"enabled"`
	Client    queue.ClientConfig    `mapstructure:"client"`
	Server    queue.ServerConfig    `mapstructure:"server"`
	Scheduler queue.SchedulerConfig `mapstructure:"scheduler"`
}

const placeholderJWTSecret = "replace-this-secret-in-production"

func (cfg *AppConfig) NormalizeAndValidate() error {
	if cfg.App.Name == "" {
		cfg.App.Name = "handygo-web"
	}
	if strings.TrimSpace(cfg.Auth.JWTSecret) == "" {
		return errors.New("auth.jwtSecret cannot be empty")
	}
	if isProductionMode(cfg.Server.Mode) && cfg.Auth.JWTSecret == placeholderJWTSecret {
		return errors.New("auth.jwtSecret must be changed for production")
	}
	if cfg.Auth.TokenTTL <= 0 {
		return errors.New("auth.tokenTTL must be greater than zero")
	}
	if cfg.Auth.PasswordCost <= 0 {
		return errors.New("auth.passwordCost must be greater than zero")
	}
	return nil
}

func isProductionMode(mode string) bool {
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case "release", "prod", "production":
		return true
	default:
		return false
	}
}

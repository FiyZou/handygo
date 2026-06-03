package logger

import (
	"errors"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Level             string   `mapstructure:"level" json:"level" yaml:"level"`
	Mode              string   `mapstructure:"mode" json:"mode" yaml:"mode"`
	Encoding          string   `mapstructure:"encoding" json:"encoding" yaml:"encoding"`
	OutputPaths       []string `mapstructure:"outputPaths" json:"outputPaths" yaml:"outputPaths"`
	ErrorOutputPaths  []string `mapstructure:"errorOutputPaths" json:"errorOutputPaths" yaml:"errorOutputPaths"`
	DisableCaller     bool     `mapstructure:"disableCaller" json:"disableCaller" yaml:"disableCaller"`
	DisableStacktrace bool     `mapstructure:"disableStacktrace" json:"disableStacktrace" yaml:"disableStacktrace"`
}

var (
	defaultLogger = zap.NewNop()
	defaultMu     sync.RWMutex
)

func New(cfg Config, opts ...zap.Option) (*zap.Logger, error) {
	zapCfg := buildConfig(cfg)
	if err := parseLevel(cfg.Level, &zapCfg.Level); err != nil {
		return nil, err
	}
	return zapCfg.Build(opts...)
}

func MustNew(cfg Config, opts ...zap.Option) *zap.Logger {
	logger, err := New(cfg, opts...)
	if err != nil {
		panic(err)
	}
	return logger
}

func SetDefault(logger *zap.Logger) {
	if logger == nil {
		logger = zap.NewNop()
	}
	defaultMu.Lock()
	defaultLogger = logger
	defaultMu.Unlock()
}

func L() *zap.Logger {
	defaultMu.RLock()
	defer defaultMu.RUnlock()
	return defaultLogger
}

func S() *zap.SugaredLogger {
	return L().Sugar()
}

func Sync() error {
	return L().Sync()
}

func buildConfig(cfg Config) zap.Config {
	mode := strings.ToLower(cfg.Mode)
	var zapCfg zap.Config
	if mode == "prod" || mode == "production" {
		zapCfg = zap.NewProductionConfig()
	} else {
		zapCfg = zap.NewDevelopmentConfig()
	}

	if cfg.Encoding != "" {
		zapCfg.Encoding = cfg.Encoding
	}
	if len(cfg.OutputPaths) > 0 {
		zapCfg.OutputPaths = cfg.OutputPaths
	}
	if len(cfg.ErrorOutputPaths) > 0 {
		zapCfg.ErrorOutputPaths = cfg.ErrorOutputPaths
	}
	zapCfg.DisableCaller = cfg.DisableCaller
	zapCfg.DisableStacktrace = cfg.DisableStacktrace
	return zapCfg
}

func parseLevel(level string, target *zap.AtomicLevel) error {
	if level == "" {
		return nil
	}
	parsed := zapcore.InfoLevel
	if err := parsed.UnmarshalText([]byte(strings.ToLower(level))); err != nil {
		return errors.New("invalid logger level: " + level)
	}
	*target = zap.NewAtomicLevelAt(parsed)
	return nil
}

package logger

import (
	"errors"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Config struct {
	Level               string         `mapstructure:"level" json:"level" yaml:"level"`
	Mode                string         `mapstructure:"mode" json:"mode" yaml:"mode"`
	Encoding            string         `mapstructure:"encoding" json:"encoding" yaml:"encoding"`
	TimeFormat          string         `mapstructure:"timeFormat" json:"timeFormat" yaml:"timeFormat"`
	OutputPaths         []string       `mapstructure:"outputPaths" json:"outputPaths" yaml:"outputPaths"`
	ErrorOutputPaths    []string       `mapstructure:"errorOutputPaths" json:"errorOutputPaths" yaml:"errorOutputPaths"`
	FileOutputPath      string         `mapstructure:"fileOutputPath" json:"fileOutputPath" yaml:"fileOutputPath"`
	ErrorFileOutputPath string         `mapstructure:"errorFileOutputPath" json:"errorFileOutputPath" yaml:"errorFileOutputPath"`
	Rotation            RotationConfig `mapstructure:"rotation" json:"rotation" yaml:"rotation"`
	DisableCaller       bool           `mapstructure:"disableCaller" json:"disableCaller" yaml:"disableCaller"`
	DisableStacktrace   bool           `mapstructure:"disableStacktrace" json:"disableStacktrace" yaml:"disableStacktrace"`
}

type RotationConfig struct {
	Enabled    bool `mapstructure:"enabled" json:"enabled" yaml:"enabled"`
	MaxSizeMB  int  `mapstructure:"maxSizeMB" json:"maxSizeMB" yaml:"maxSizeMB"`
	MaxAgeDays int  `mapstructure:"maxAgeDays" json:"maxAgeDays" yaml:"maxAgeDays"`
	MaxBackups int  `mapstructure:"maxBackups" json:"maxBackups" yaml:"maxBackups"`
	LocalTime  bool `mapstructure:"localTime" json:"localTime" yaml:"localTime"`
	Compress   bool `mapstructure:"compress" json:"compress" yaml:"compress"`
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
	if usesRotatingFiles(cfg) {
		return buildRotatingLogger(cfg, zapCfg, opts...)
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
	if cfg.TimeFormat != "" {
		zapCfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(cfg.TimeFormat)
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

func usesRotatingFiles(cfg Config) bool {
	return cfg.Rotation.Enabled && (cfg.FileOutputPath != "" || cfg.ErrorFileOutputPath != "")
}

func buildRotatingLogger(cfg Config, zapCfg zap.Config, opts ...zap.Option) (*zap.Logger, error) {
	level := zapCfg.Level

	var cores []zapcore.Core
	if len(zapCfg.OutputPaths) > 0 {
		sink, _, err := zap.Open(zapCfg.OutputPaths...)
		if err != nil {
			return nil, err
		}
		cores = append(cores, zapcore.NewCore(encoder(zapCfg), sink, level))
	}
	if cfg.FileOutputPath != "" {
		cores = append(cores, zapcore.NewCore(encoder(zapCfg), zapcore.AddSync(rotatingWriter(cfg.FileOutputPath, cfg.Rotation)), level))
	}
	if cfg.ErrorFileOutputPath != "" {
		cores = append(cores, zapcore.NewCore(encoder(zapCfg), zapcore.AddSync(rotatingWriter(cfg.ErrorFileOutputPath, cfg.Rotation)), zapcore.ErrorLevel))
	}
	if len(cores) == 0 {
		cores = append(cores, zapcore.NewCore(encoder(zapCfg), zapcore.AddSync(rotatingWriter("logs/app.log", cfg.Rotation)), level))
	}

	options, err := buildOptions(zapCfg, opts...)
	if err != nil {
		return nil, err
	}
	logger := zap.New(zapcore.NewTee(cores...), options...)
	return logger, nil
}

func encoder(cfg zap.Config) zapcore.Encoder {
	switch strings.ToLower(cfg.Encoding) {
	case "console":
		return zapcore.NewConsoleEncoder(cfg.EncoderConfig)
	default:
		return zapcore.NewJSONEncoder(cfg.EncoderConfig)
	}
}

func rotatingWriter(path string, cfg RotationConfig) *lumberjack.Logger {
	if cfg.MaxSizeMB <= 0 {
		cfg.MaxSizeMB = 100
	}
	return &lumberjack.Logger{
		Filename:   path,
		MaxSize:    cfg.MaxSizeMB,
		MaxAge:     cfg.MaxAgeDays,
		MaxBackups: cfg.MaxBackups,
		LocalTime:  cfg.LocalTime,
		Compress:   cfg.Compress,
	}
}

func buildOptions(cfg zap.Config, opts ...zap.Option) ([]zap.Option, error) {
	options := make([]zap.Option, 0, len(opts)+4)
	if !cfg.DisableCaller {
		options = append(options, zap.AddCaller())
	}
	if !cfg.DisableStacktrace {
		options = append(options, zap.AddStacktrace(zapcore.ErrorLevel))
	}
	if len(cfg.ErrorOutputPaths) > 0 {
		sink, _, err := zap.Open(cfg.ErrorOutputPaths...)
		if err != nil {
			return nil, err
		}
		options = append(options, zap.ErrorOutput(sink))
	}
	options = append(options, opts...)
	return options, nil
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

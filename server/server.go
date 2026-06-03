package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Config struct {
	Name              string        `mapstructure:"name" json:"name" yaml:"name"`
	Addr              string        `mapstructure:"addr" json:"addr" yaml:"addr"`
	Mode              string        `mapstructure:"mode" json:"mode" yaml:"mode"`
	ReadTimeout       time.Duration `mapstructure:"readTimeout" json:"readTimeout" yaml:"readTimeout"`
	ReadHeaderTimeout time.Duration `mapstructure:"readHeaderTimeout" json:"readHeaderTimeout" yaml:"readHeaderTimeout"`
	WriteTimeout      time.Duration `mapstructure:"writeTimeout" json:"writeTimeout" yaml:"writeTimeout"`
	IdleTimeout       time.Duration `mapstructure:"idleTimeout" json:"idleTimeout" yaml:"idleTimeout"`
}

type RegisterFunc func(*gin.Engine)

type Server struct {
	cfg    Config
	engine *gin.Engine
	server *http.Server
	logger *zap.Logger
}

func New(cfg Config, logger *zap.Logger, opts ...Option) *Server {
	if cfg.Name == "" {
		cfg.Name = "http"
	}
	if cfg.Addr == "" {
		cfg.Addr = ":8080"
	}
	if cfg.Mode != "" {
		gin.SetMode(cfg.Mode)
	}
	if logger == nil {
		logger = zap.NewNop()
	}

	s := &Server{
		cfg:    cfg,
		engine: gin.New(),
		logger: logger,
	}
	for _, opt := range opts {
		opt(s)
	}
	s.server = &http.Server{
		Addr:              cfg.Addr,
		Handler:           s.engine,
		ReadTimeout:       cfg.ReadTimeout,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
	}
	return s
}

type Option func(*Server)

func WithEngine(engine *gin.Engine) Option {
	return func(s *Server) {
		if engine != nil {
			s.engine = engine
		}
	}
}

func (s *Server) Name() string {
	return s.cfg.Name
}

func (s *Server) Router() *gin.Engine {
	return s.engine
}

func (s *Server) Use(middleware ...gin.HandlerFunc) {
	s.engine.Use(middleware...)
}

func (s *Server) Register(register RegisterFunc) {
	if register != nil {
		register(s.engine)
	}
}

func (s *Server) Start(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		s.logger.Info("http server started", zap.String("addr", s.server.Addr))
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		return err
	case <-time.After(50 * time.Millisecond):
		return nil
	}
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("http server stopping", zap.String("addr", s.server.Addr))
	return s.server.Shutdown(ctx)
}

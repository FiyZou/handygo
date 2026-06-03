package queue

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
)

type RedisConfig struct {
	Addr     string `mapstructure:"addr" json:"addr" yaml:"addr"`
	Username string `mapstructure:"username" json:"username" yaml:"username"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
	DB       int    `mapstructure:"db" json:"db" yaml:"db"`
}

type ClientConfig struct {
	Name  string      `mapstructure:"name" json:"name" yaml:"name"`
	Redis RedisConfig `mapstructure:"redis" json:"redis" yaml:"redis"`
}

type ServerConfig struct {
	Name        string         `mapstructure:"name" json:"name" yaml:"name"`
	Redis       RedisConfig    `mapstructure:"redis" json:"redis" yaml:"redis"`
	Concurrency int            `mapstructure:"concurrency" json:"concurrency" yaml:"concurrency"`
	Queues      map[string]int `mapstructure:"queues" json:"queues" yaml:"queues"`
}

type SchedulerConfig struct {
	Name     string      `mapstructure:"name" json:"name" yaml:"name"`
	Redis    RedisConfig `mapstructure:"redis" json:"redis" yaml:"redis"`
	Location string      `mapstructure:"location" json:"location" yaml:"location"`
}

type Client struct {
	cfg    ClientConfig
	client *asynq.Client
}

func NewClient(cfg ClientConfig) *Client {
	if cfg.Name == "" {
		cfg.Name = "asynq-client"
	}
	return &Client{
		cfg:    cfg,
		client: asynq.NewClient(redisOpt(cfg.Redis)),
	}
}

func (c *Client) Name() string {
	return c.cfg.Name
}

func (c *Client) Start(ctx context.Context) error {
	if c.client == nil {
		return errors.New("asynq client is not initialized")
	}
	return nil
}

func (c *Client) Stop(ctx context.Context) error {
	if c.client == nil {
		return nil
	}
	return c.client.Close()
}

func (c *Client) Raw() *asynq.Client {
	return c.client
}

func (c *Client) Enqueue(ctx context.Context, task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	if c.client == nil {
		return nil, errors.New("asynq client is not initialized")
	}
	return c.client.EnqueueContext(ctx, task, opts...)
}

type Server struct {
	cfg    ServerConfig
	server *asynq.Server
	mux    *asynq.ServeMux
}

func NewServer(cfg ServerConfig) *Server {
	if cfg.Name == "" {
		cfg.Name = "asynq-server"
	}
	if cfg.Concurrency <= 0 {
		cfg.Concurrency = 10
	}
	if len(cfg.Queues) == 0 {
		cfg.Queues = map[string]int{"default": 1}
	}
	return &Server{
		cfg: cfg,
		server: asynq.NewServer(redisOpt(cfg.Redis), asynq.Config{
			Concurrency: cfg.Concurrency,
			Queues:      cfg.Queues,
		}),
		mux: asynq.NewServeMux(),
	}
}

func (s *Server) Name() string {
	return s.cfg.Name
}

func (s *Server) Handle(taskType string, handler asynq.Handler) {
	s.mux.Handle(taskType, handler)
}

func (s *Server) HandleFunc(taskType string, handler func(context.Context, *asynq.Task) error) {
	s.mux.HandleFunc(taskType, handler)
}

func (s *Server) Start(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		errCh <- s.server.Run(s.mux)
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
	done := make(chan struct{})
	go func() {
		s.server.Shutdown()
		close(done)
	}()
	select {
	case <-ctx.Done():
		s.server.Stop()
		return ctx.Err()
	case <-done:
		return nil
	}
}

type Scheduler struct {
	cfg       SchedulerConfig
	scheduler *asynq.Scheduler
}

func NewScheduler(cfg SchedulerConfig) (*Scheduler, error) {
	if cfg.Name == "" {
		cfg.Name = "asynq-scheduler"
	}
	location := time.Local
	if cfg.Location != "" {
		parsed, err := time.LoadLocation(cfg.Location)
		if err != nil {
			return nil, fmt.Errorf("load scheduler location: %w", err)
		}
		location = parsed
	}
	return &Scheduler{
		cfg: cfg,
		scheduler: asynq.NewScheduler(redisOpt(cfg.Redis), &asynq.SchedulerOpts{
			Location: location,
		}),
	}, nil
}

func (s *Scheduler) Name() string {
	return s.cfg.Name
}

func (s *Scheduler) Register(cronspec string, task *asynq.Task, opts ...asynq.Option) (string, error) {
	return s.scheduler.Register(cronspec, task, opts...)
}

func (s *Scheduler) Start(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		errCh <- s.scheduler.Run()
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

func (s *Scheduler) Stop(ctx context.Context) error {
	done := make(chan struct{})
	go func() {
		s.scheduler.Shutdown()
		close(done)
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	}
}

func redisOpt(cfg RedisConfig) asynq.RedisClientOpt {
	if cfg.Addr == "" {
		cfg.Addr = "127.0.0.1:6379"
	}
	return asynq.RedisClientOpt{
		Addr:     cfg.Addr,
		Username: cfg.Username,
		Password: cfg.Password,
		DB:       cfg.DB,
	}
}

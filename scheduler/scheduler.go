package scheduler

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/FiyZou/handygo/safego"
	"github.com/robfig/cron/v3"
)

type Job func(context.Context) error

type Logger interface {
	Errorw(msg string, keysAndValues ...any)
	Infow(msg string, keysAndValues ...any)
}

type Config struct {
	Name     string `mapstructure:"name" json:"name" yaml:"name"`
	Location string `mapstructure:"location" json:"location" yaml:"location"`
	Seconds  bool   `mapstructure:"seconds" json:"seconds" yaml:"seconds"`
}

type Task struct {
	Name       string
	Spec       string
	RunOnStart bool
	Job        Job
}

type Scheduler struct {
	cfg     Config
	logger  Logger
	tasks   []Task
	cron    *cron.Cron
	ctx     context.Context
	cancel  context.CancelFunc
	mu      sync.Mutex
	started bool
}

func New(cfg Config, logger Logger) *Scheduler {
	if cfg.Name == "" {
		cfg.Name = "scheduler"
	}
	return &Scheduler{cfg: cfg, logger: logger}
}

func (s *Scheduler) Name() string {
	return s.cfg.Name
}

func (s *Scheduler) Add(task Task) error {
	if task.Name == "" {
		return errors.New("task name cannot be empty")
	}
	if task.Spec == "" {
		return errors.New("task spec cannot be empty")
	}
	if task.Job == nil {
		return errors.New("task job cannot be nil")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.started {
		return errors.New("cannot add task after scheduler started")
	}
	s.tasks = append(s.tasks, task)
	return nil
}

func (s *Scheduler) Start(ctx context.Context) error {
	s.mu.Lock()
	if s.started {
		s.mu.Unlock()
		return nil
	}
	s.ctx, s.cancel = context.WithCancel(ctx)
	c := cron.New(s.cronOptions()...)
	for _, task := range s.tasks {
		task := task
		if _, err := c.AddFunc(task.Spec, func() {
			safego.Go(func() {
				s.execute(task)
			}, safego.WithLogger(s.logger))
		}); err != nil {
			s.cancel()
			s.mu.Unlock()
			return fmt.Errorf("add cron task %s: %w", task.Name, err)
		}
		if task.RunOnStart {
			safego.Go(func() {
				s.execute(task)
			}, safego.WithLogger(s.logger))
		}
	}
	s.cron = c
	s.started = true
	s.mu.Unlock()

	c.Start()
	return nil
}

func (s *Scheduler) Stop(ctx context.Context) error {
	s.mu.Lock()
	if !s.started {
		s.mu.Unlock()
		return nil
	}
	c := s.cron
	s.cancel()
	s.started = false
	s.mu.Unlock()

	done := c.Stop().Done()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	}
}

func (s *Scheduler) execute(task Task) {
	if s.logger != nil {
		s.logger.Infow("scheduled task started", "task", task.Name, "spec", task.Spec)
	}
	if err := task.Job(s.ctx); err != nil && s.logger != nil {
		s.logger.Errorw("scheduled task failed", "task", task.Name, "spec", task.Spec, "error", err)
	}
}

func (s *Scheduler) cronOptions() []cron.Option {
	options := []cron.Option{
		cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)),
	}
	if s.cfg.Seconds {
		options = append(options, cron.WithSeconds())
	}
	if s.cfg.Location != "" {
		location, err := time.LoadLocation(s.cfg.Location)
		if err != nil {
			if s.logger != nil {
				s.logger.Errorw("load scheduler location failed", "location", s.cfg.Location, "error", err)
			}
			location = time.Local
		}
		options = append(options, cron.WithLocation(location))
	}
	return options
}

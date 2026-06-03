package workerpool

import (
	"context"
	"errors"
	"sync"

	"github.com/FiyZou/handygo/safego"
)

type Job func(context.Context) error

type Logger interface {
	Errorw(msg string, keysAndValues ...any)
}

type Config struct {
	Name    string `mapstructure:"name" json:"name" yaml:"name"`
	Workers int    `mapstructure:"workers" json:"workers" yaml:"workers"`
	Queue   int    `mapstructure:"queue" json:"queue" yaml:"queue"`
}

type Pool struct {
	cfg     Config
	logger  Logger
	jobs    chan Job
	ctx     context.Context
	cancel  context.CancelFunc
	wg      sync.WaitGroup
	mu      sync.Mutex
	started bool
}

func New(cfg Config, logger Logger) *Pool {
	if cfg.Name == "" {
		cfg.Name = "workerpool"
	}
	if cfg.Workers <= 0 {
		cfg.Workers = 4
	}
	if cfg.Queue <= 0 {
		cfg.Queue = cfg.Workers * 16
	}
	return &Pool{cfg: cfg, logger: logger, jobs: make(chan Job, cfg.Queue)}
}

func (p *Pool) Name() string {
	return p.cfg.Name
}

func (p *Pool) Start(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.started {
		return nil
	}
	p.ctx, p.cancel = context.WithCancel(ctx)
	p.started = true
	for i := 0; i < p.cfg.Workers; i++ {
		p.wg.Add(1)
		safego.Go(func() {
			defer p.wg.Done()
			p.worker()
		}, safego.WithLogger(p.logger))
	}
	return nil
}

func (p *Pool) Stop(ctx context.Context) error {
	p.mu.Lock()
	if !p.started {
		p.mu.Unlock()
		return nil
	}
	p.cancel()
	p.started = false
	p.mu.Unlock()

	done := make(chan struct{})
	go func() {
		p.wg.Wait()
		close(done)
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	}
}

func (p *Pool) Submit(ctx context.Context, job Job) error {
	if job == nil {
		return errors.New("job cannot be nil")
	}
	p.mu.Lock()
	started := p.started
	p.mu.Unlock()
	if !started {
		return errors.New("worker pool is not started")
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case p.jobs <- job:
		return nil
	}
}

func (p *Pool) worker() {
	for {
		select {
		case <-p.ctx.Done():
			return
		case job := <-p.jobs:
			if job == nil {
				continue
			}
			if err := job(p.ctx); err != nil && p.logger != nil {
				p.logger.Errorw("worker job failed", "pool", p.cfg.Name, "error", err)
			}
		}
	}
}

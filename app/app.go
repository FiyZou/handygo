package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Component interface {
	Name() string
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type Logger interface {
	Infow(msg string, keysAndValues ...any)
	Errorw(msg string, keysAndValues ...any)
}

type Option func(*App)

type App struct {
	name            string
	components      []Component
	logger          Logger
	shutdownTimeout time.Duration
	mu              sync.Mutex
	started         bool
}

func New(name string, opts ...Option) *App {
	a := &App{
		name:            name,
		shutdownTimeout: 15 * time.Second,
	}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

func WithLogger(logger Logger) Option {
	return func(a *App) {
		a.logger = logger
	}
}

func WithShutdownTimeout(timeout time.Duration) Option {
	return func(a *App) {
		if timeout > 0 {
			a.shutdownTimeout = timeout
		}
	}
}

func (a *App) Name() string {
	return a.name
}

func (a *App) Register(components ...Component) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.components = append(a.components, components...)
}

func (a *App) Start(ctx context.Context) error {
	a.mu.Lock()
	if a.started {
		a.mu.Unlock()
		return nil
	}
	components := append([]Component(nil), a.components...)
	a.started = true
	a.mu.Unlock()

	started := make([]Component, 0, len(components))
	for _, component := range components {
		a.info("starting component", "name", component.Name())
		if err := component.Start(ctx); err != nil {
			a.error("component start failed", "name", component.Name(), "error", err)
			stopErr := stopComponents(ctx, started)
			a.mu.Lock()
			a.started = false
			a.mu.Unlock()
			if stopErr != nil {
				return errors.Join(fmt.Errorf("start %s: %w", component.Name(), err), stopErr)
			}
			return fmt.Errorf("start %s: %w", component.Name(), err)
		}
		started = append(started, component)
	}
	return nil
}

func (a *App) Stop(ctx context.Context) error {
	a.mu.Lock()
	if !a.started {
		a.mu.Unlock()
		return nil
	}
	components := append([]Component(nil), a.components...)
	a.started = false
	a.mu.Unlock()

	return stopComponents(ctx, components)
}

func (a *App) Run(ctx context.Context) error {
	if err := a.Start(ctx); err != nil {
		return err
	}

	sigCtx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	<-sigCtx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), a.shutdownTimeout)
	defer cancel()
	return a.Stop(shutdownCtx)
}

func stopComponents(ctx context.Context, components []Component) error {
	var errs []error
	for i := len(components) - 1; i >= 0; i-- {
		component := components[i]
		if err := component.Stop(ctx); err != nil {
			errs = append(errs, fmt.Errorf("stop %s: %w", component.Name(), err))
		}
	}
	return errors.Join(errs...)
}

func (a *App) info(msg string, fields ...any) {
	if a.logger != nil {
		a.logger.Infow(msg, fields...)
	}
}

func (a *App) error(msg string, fields ...any) {
	if a.logger != nil {
		a.logger.Errorw(msg, fields...)
	}
}

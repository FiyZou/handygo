package bootstrap

import (
	"context"
	"fmt"

	"github.com/FiyZou/handygo/app"
	exampleconfig "github.com/FiyZou/handygo/examples/internal/config"
	"github.com/FiyZou/handygo/examples/internal/tasks"
	"github.com/FiyZou/handygo/queue"
	"github.com/FiyZou/handygo/scheduler"
	"github.com/FiyZou/handygo/workerpool"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

func newLocalWorkers(cfg exampleconfig.AppConfig, zapLogger *zap.Logger) (*workerpool.Pool, *scheduler.Scheduler, error) {
	localPool := workerpool.New(cfg.Worker.Pool, zapLogger.Sugar())
	localScheduler := scheduler.New(cfg.Worker.Scheduler, zapLogger.Sugar())
	if err := localScheduler.Add(scheduler.Task{
		Name:       "local-health-log",
		Spec:       "*/1 * * * *",
		RunOnStart: true,
		Job: func(ctx context.Context) error {
			return localPool.Submit(ctx, func(ctx context.Context) error {
				zapLogger.Info("local scheduled task executed")
				return nil
			})
		},
	}); err != nil {
		return nil, nil, fmt.Errorf("add local scheduled task: %w", err)
	}
	return localPool, localScheduler, nil
}

func registerAsynq(application *app.App, cfg exampleconfig.AppConfig, zapLogger *zap.Logger) error {
	client := queue.NewClient(cfg.Asynq.Client)
	server := queue.NewServer(cfg.Asynq.Server)
	tasks.Register(server, zapLogger)

	task, err := tasks.NewUserReportTask(1)
	if err != nil {
		return fmt.Errorf("new user report task: %w", err)
	}
	asynqScheduler, err := queue.NewScheduler(cfg.Asynq.Scheduler)
	if err != nil {
		return fmt.Errorf("new asynq scheduler: %w", err)
	}
	if _, err := asynqScheduler.Register("*/5 * * * *", task, asynq.Queue("default")); err != nil {
		return fmt.Errorf("register asynq task: %w", err)
	}

	application.Register(client, server, asynqScheduler)
	return nil
}

package scheduler

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestSchedulerRunsTaskOnStart(t *testing.T) {
	s := New(Config{Name: "test"}, nil)
	var count int64
	if err := s.Add(Task{
		Name:       "tick",
		Spec:       "0 0 1 1 *",
		RunOnStart: true,
		Job: func(ctx context.Context) error {
			atomic.AddInt64(&count, 1)
			return nil
		},
	}); err != nil {
		t.Fatalf("add task: %v", err)
	}

	if err := s.Start(context.Background()); err != nil {
		t.Fatalf("start scheduler: %v", err)
	}
	defer func() {
		_ = s.Stop(context.Background())
	}()

	deadline := time.After(time.Second)
	for atomic.LoadInt64(&count) != 1 {
		select {
		case <-deadline:
			t.Fatalf("count = %d", count)
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func TestSchedulerRejectsInvalidSpec(t *testing.T) {
	s := New(Config{Name: "test"}, nil)
	if err := s.Add(Task{
		Name: "bad",
		Spec: "invalid",
		Job:  func(ctx context.Context) error { return nil },
	}); err != nil {
		t.Fatalf("add task: %v", err)
	}
	if err := s.Start(context.Background()); err == nil {
		t.Fatal("expected invalid cron spec error")
	}
}

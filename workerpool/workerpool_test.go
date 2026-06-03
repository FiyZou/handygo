package workerpool

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestPoolRunsJobs(t *testing.T) {
	pool := New(Config{Workers: 2, Queue: 4}, nil)
	if err := pool.Start(context.Background()); err != nil {
		t.Fatalf("start pool: %v", err)
	}
	defer func() {
		_ = pool.Stop(context.Background())
	}()

	var count int64
	for i := 0; i < 3; i++ {
		if err := pool.Submit(context.Background(), func(ctx context.Context) error {
			atomic.AddInt64(&count, 1)
			return nil
		}); err != nil {
			t.Fatalf("submit job: %v", err)
		}
	}

	deadline := time.After(time.Second)
	for atomic.LoadInt64(&count) != 3 {
		select {
		case <-deadline:
			t.Fatalf("count = %d", count)
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}

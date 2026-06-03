package app

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

type testComponent struct {
	name    string
	events  *[]string
	start   error
	stop    error
	started bool
}

func (c *testComponent) Name() string {
	return c.name
}

func (c *testComponent) Start(ctx context.Context) error {
	*c.events = append(*c.events, "start:"+c.name)
	c.started = true
	return c.start
}

func (c *testComponent) Stop(ctx context.Context) error {
	*c.events = append(*c.events, "stop:"+c.name)
	return c.stop
}

func TestAppStartStopOrder(t *testing.T) {
	events := []string{}
	a := New("test")
	a.Register(&testComponent{name: "one", events: &events}, &testComponent{name: "two", events: &events})

	if err := a.Start(context.Background()); err != nil {
		t.Fatalf("start failed: %v", err)
	}
	if err := a.Stop(context.Background()); err != nil {
		t.Fatalf("stop failed: %v", err)
	}

	want := []string{"start:one", "start:two", "stop:two", "stop:one"}
	if !reflect.DeepEqual(events, want) {
		t.Fatalf("events = %v, want %v", events, want)
	}
}

func TestAppRollsBackStartedComponents(t *testing.T) {
	events := []string{}
	a := New("test")
	a.Register(
		&testComponent{name: "one", events: &events},
		&testComponent{name: "two", events: &events, start: errors.New("boom")},
	)

	if err := a.Start(context.Background()); err == nil {
		t.Fatal("expected start error")
	}

	want := []string{"start:one", "start:two", "stop:one"}
	if !reflect.DeepEqual(events, want) {
		t.Fatalf("events = %v, want %v", events, want)
	}
}

package logger

import (
	"testing"

	"go.uber.org/zap"
)

func TestNewLogger(t *testing.T) {
	l, err := New(Config{
		Level:             "debug",
		Mode:              "dev",
		Encoding:          "json",
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
		DisableCaller:     true,
		DisableStacktrace: true,
	})
	if err != nil {
		t.Fatalf("new logger: %v", err)
	}
	if l == nil {
		t.Fatal("logger is nil")
	}
}

func TestDefaultLogger(t *testing.T) {
	l := zap.NewNop()
	SetDefault(l)
	if L() != l {
		t.Fatal("default logger was not set")
	}
}

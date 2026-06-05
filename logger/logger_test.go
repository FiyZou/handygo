package logger

import (
	"os"
	"path/filepath"
	"strings"
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

func TestNewLoggerWithTimeFormatAndRotation(t *testing.T) {
	logPath := filepath.Join(t.TempDir(), "app.log")
	l, err := New(Config{
		Level:          "debug",
		Mode:           "prod",
		Encoding:       "json",
		TimeFormat:     "2006-01-02 15:04:05",
		FileOutputPath: logPath,
		Rotation: RotationConfig{
			Enabled:    true,
			MaxSizeMB:  1,
			MaxAgeDays: 1,
			MaxBackups: 1,
			LocalTime:  true,
			Compress:   false,
		},
		DisableCaller:     true,
		DisableStacktrace: true,
	})
	if err != nil {
		t.Fatalf("new logger: %v", err)
	}
	l.Info("rotation test")

	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("read log file: %v", err)
	}
	if !strings.Contains(string(data), "rotation test") {
		t.Fatalf("log file missing message:\n%s", string(data))
	}
}

package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	content := []byte("server:\n  addr: ':9090'\n  readTimeout: 2s\nlogger:\n  level: debug\n")
	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	var cfg struct {
		Server struct {
			Addr        string        `mapstructure:"addr"`
			ReadTimeout time.Duration `mapstructure:"readTimeout"`
		} `mapstructure:"server"`
		Logger struct {
			Level string `mapstructure:"level"`
		} `mapstructure:"logger"`
	}

	if err := Load(path, &cfg); err != nil {
		t.Fatalf("load config: %v", err)
	}
	if cfg.Server.Addr != ":9090" {
		t.Fatalf("addr = %q", cfg.Server.Addr)
	}
	if cfg.Server.ReadTimeout != 2*time.Second {
		t.Fatalf("read timeout = %s", cfg.Server.ReadTimeout)
	}
	if cfg.Logger.Level != "debug" {
		t.Fatalf("level = %q", cfg.Logger.Level)
	}
}

func TestLoadDataYAML(t *testing.T) {
	var cfg struct {
		Server struct {
			Addr        string        `mapstructure:"addr"`
			ReadTimeout time.Duration `mapstructure:"readTimeout"`
		} `mapstructure:"server"`
	}

	data := []byte("server:\n  addr: ':9090'\n  readTimeout: 2s\n")
	if err := LoadData("config", "yaml", data, &cfg); err != nil {
		t.Fatalf("load config data: %v", err)
	}
	if cfg.Server.Addr != ":9090" {
		t.Fatalf("addr = %q", cfg.Server.Addr)
	}
	if cfg.Server.ReadTimeout != 2*time.Second {
		t.Fatalf("read timeout = %s", cfg.Server.ReadTimeout)
	}
}

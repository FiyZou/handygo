package config

import (
	"testing"
	"time"
)

func TestNormalizeAndValidateRejectsMissingJWTSecret(t *testing.T) {
	cfg := validConfig()
	cfg.Auth.JWTSecret = ""

	if err := cfg.NormalizeAndValidate(); err == nil {
		t.Fatal("expected missing jwt secret error")
	}
}

func TestNormalizeAndValidateRejectsProductionPlaceholderJWTSecret(t *testing.T) {
	cfg := validConfig()
	cfg.Server.Mode = "release"
	cfg.Auth.JWTSecret = placeholderJWTSecret

	if err := cfg.NormalizeAndValidate(); err == nil {
		t.Fatal("expected production placeholder jwt secret error")
	}
}

func TestNormalizeAndValidateAllowsLocalPlaceholderJWTSecret(t *testing.T) {
	cfg := validConfig()
	cfg.App.Name = ""
	cfg.Server.Mode = "debug"
	cfg.Auth.JWTSecret = placeholderJWTSecret

	if err := cfg.NormalizeAndValidate(); err != nil {
		t.Fatalf("validate config: %v", err)
	}
	if cfg.App.Name != "handygo-web" {
		t.Fatalf("app name = %q", cfg.App.Name)
	}
}

func TestNormalizeAndValidateRejectsInvalidAuthSettings(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*AppConfig)
	}{
		{
			name: "token ttl",
			mutate: func(cfg *AppConfig) {
				cfg.Auth.TokenTTL = 0
			},
		},
		{
			name: "password cost",
			mutate: func(cfg *AppConfig) {
				cfg.Auth.PasswordCost = 0
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := validConfig()
			tt.mutate(&cfg)
			if err := cfg.NormalizeAndValidate(); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

func validConfig() AppConfig {
	return AppConfig{
		App: App{Name: "handygo-web"},
		Auth: Auth{
			JWTSecret:    "test-secret",
			TokenTTL:     time.Hour,
			PasswordCost: 4,
		},
	}
}

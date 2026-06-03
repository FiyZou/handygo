package config

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type Options struct {
	Path           string
	EnvPrefix      string
	EnvKeyReplacer *strings.Replacer
	Defaults       map[string]any
}

type DataOptions struct {
	Name           string
	Type           string
	EnvPrefix      string
	EnvKeyReplacer *strings.Replacer
	Defaults       map[string]any
}

func Load(path string, out any) error {
	return LoadWithOptions(Options{Path: path}, out)
}

func LoadWithOptions(opts Options, out any) error {
	if out == nil {
		return errors.New("config output cannot be nil")
	}
	if opts.Path == "" {
		return errors.New("config path cannot be empty")
	}

	v := viper.New()
	applyCommonOptions(v, opts.EnvPrefix, opts.EnvKeyReplacer, opts.Defaults)

	v.SetConfigFile(opts.Path)
	if configType := strings.TrimPrefix(filepath.Ext(opts.Path), "."); configType != "" {
		v.SetConfigType(configType)
	}

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("read config: %w", err)
	}
	if err := v.Unmarshal(out); err != nil {
		return fmt.Errorf("unmarshal config: %w", err)
	}
	return nil
}

func LoadData(name string, configType string, data []byte, out any) error {
	return LoadReaderWithOptions(DataOptions{Name: name, Type: configType}, bytes.NewReader(data), out)
}

func LoadReaderWithOptions(opts DataOptions, reader io.Reader, out any) error {
	if out == nil {
		return errors.New("config output cannot be nil")
	}
	if reader == nil {
		return errors.New("config reader cannot be nil")
	}
	if opts.Type == "" {
		return errors.New("config type cannot be empty")
	}

	v := viper.New()
	applyCommonOptions(v, opts.EnvPrefix, opts.EnvKeyReplacer, opts.Defaults)
	if opts.Name != "" {
		v.SetConfigName(opts.Name)
	}
	v.SetConfigType(opts.Type)

	if err := v.ReadConfig(reader); err != nil {
		return fmt.Errorf("read config: %w", err)
	}
	if err := v.Unmarshal(out); err != nil {
		return fmt.Errorf("unmarshal config: %w", err)
	}
	return nil
}

func applyCommonOptions(v *viper.Viper, envPrefix string, envKeyReplacer *strings.Replacer, defaults map[string]any) {
	for key, value := range defaults {
		v.SetDefault(key, value)
	}
	if envPrefix != "" {
		v.SetEnvPrefix(envPrefix)
	}
	if envKeyReplacer != nil {
		v.SetEnvKeyReplacer(envKeyReplacer)
	} else {
		v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	}
	v.AutomaticEnv()
}

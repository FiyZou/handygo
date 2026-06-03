package cache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Name         string        `mapstructure:"name" json:"name" yaml:"name"`
	Addr         string        `mapstructure:"addr" json:"addr" yaml:"addr"`
	Username     string        `mapstructure:"username" json:"username" yaml:"username"`
	Password     string        `mapstructure:"password" json:"password" yaml:"password"`
	DB           int           `mapstructure:"db" json:"db" yaml:"db"`
	DialTimeout  time.Duration `mapstructure:"dialTimeout" json:"dialTimeout" yaml:"dialTimeout"`
	ReadTimeout  time.Duration `mapstructure:"readTimeout" json:"readTimeout" yaml:"readTimeout"`
	WriteTimeout time.Duration `mapstructure:"writeTimeout" json:"writeTimeout" yaml:"writeTimeout"`
	PoolSize     int           `mapstructure:"poolSize" json:"poolSize" yaml:"poolSize"`
}

type Cache struct {
	cfg    Config
	client *redis.Client
}

func New(cfg Config) *Cache {
	if cfg.Name == "" {
		cfg.Name = "redis"
	}
	if cfg.Addr == "" {
		cfg.Addr = "127.0.0.1:6379"
	}
	return &Cache{
		cfg: cfg,
		client: redis.NewClient(&redis.Options{
			Addr:         cfg.Addr,
			Username:     cfg.Username,
			Password:     cfg.Password,
			DB:           cfg.DB,
			DialTimeout:  cfg.DialTimeout,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			PoolSize:     cfg.PoolSize,
		}),
	}
}

func (c *Cache) Name() string {
	if c == nil || c.cfg.Name == "" {
		return "redis"
	}
	return c.cfg.Name
}

func (c *Cache) Client() *redis.Client {
	if c == nil {
		return nil
	}
	return c.client
}

func (c *Cache) Start(ctx context.Context) error {
	if c == nil || c.client == nil {
		return errors.New("redis client is not initialized")
	}
	return c.client.Ping(ctx).Err()
}

func (c *Cache) Stop(ctx context.Context) error {
	if c == nil || c.client == nil {
		return nil
	}
	done := make(chan error, 1)
	go func() {
		done <- c.client.Close()
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}

func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

func (c *Cache) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return c.client.Set(ctx, key, value, expiration).Err()
}

func (c *Cache) Del(ctx context.Context, keys ...string) error {
	return c.client.Del(ctx, keys...).Err()
}

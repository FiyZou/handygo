package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	Name            string        `mapstructure:"name" json:"name" yaml:"name"`
	Driver          string        `mapstructure:"driver" json:"driver" yaml:"driver"`
	DSN             string        `mapstructure:"dsn" json:"dsn" yaml:"dsn"`
	MaxIdleConns    int           `mapstructure:"maxIdleConns" json:"maxIdleConns" yaml:"maxIdleConns"`
	MaxOpenConns    int           `mapstructure:"maxOpenConns" json:"maxOpenConns" yaml:"maxOpenConns"`
	ConnMaxLifetime time.Duration `mapstructure:"connMaxLifetime" json:"connMaxLifetime" yaml:"connMaxLifetime"`
	ConnMaxIdleTime time.Duration `mapstructure:"connMaxIdleTime" json:"connMaxIdleTime" yaml:"connMaxIdleTime"`
	SlowThreshold   time.Duration `mapstructure:"slowThreshold" json:"slowThreshold" yaml:"slowThreshold"`
	LogLevel        string        `mapstructure:"logLevel" json:"logLevel" yaml:"logLevel"`
}

type Database struct {
	cfg Config
	db  *gorm.DB
}

func New(cfg Config, opts ...gorm.Option) (*Database, error) {
	if cfg.Name == "" {
		cfg.Name = "database"
	}
	dialector, err := dialector(cfg)
	if err != nil {
		return nil, err
	}

	options := []gorm.Option{&gorm.Config{Logger: gormLogger(cfg)}}
	options = append(options, opts...)
	db, err := gorm.Open(dialector, options...)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql db: %w", err)
	}
	if cfg.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.ConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	}
	if cfg.ConnMaxIdleTime > 0 {
		sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
	}

	return &Database{cfg: cfg, db: db}, nil
}

func (d *Database) Name() string {
	if d == nil || d.cfg.Name == "" {
		return "database"
	}
	return d.cfg.Name
}

func (d *Database) DB() *gorm.DB {
	if d == nil {
		return nil
	}
	return d.db
}

func (d *Database) Start(ctx context.Context) error {
	if d == nil || d.db == nil {
		return errors.New("database is not initialized")
	}
	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}

func (d *Database) Stop(ctx context.Context) error {
	if d == nil || d.db == nil {
		return nil
	}
	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}
	done := make(chan error, 1)
	go func() {
		done <- sqlDB.Close()
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}

func (d *Database) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	if d == nil || d.db == nil {
		return errors.New("database is not initialized")
	}
	return d.db.WithContext(ctx).Transaction(fn)
}

func (d *Database) AutoMigrate(models ...any) error {
	if d == nil || d.db == nil {
		return errors.New("database is not initialized")
	}
	return d.db.AutoMigrate(models...)
}

func dialector(cfg Config) (gorm.Dialector, error) {
	if cfg.DSN == "" {
		return nil, errors.New("database dsn cannot be empty")
	}
	switch strings.ToLower(cfg.Driver) {
	case "mysql":
		return mysql.Open(cfg.DSN), nil
	case "postgres", "postgresql":
		return postgres.Open(cfg.DSN), nil
	case "sqlite", "sqlite3":
		return sqlite.Open(cfg.DSN), nil
	case "sqlserver", "mssql":
		return sqlserver.Open(cfg.DSN), nil
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}
}

func gormLogger(cfg Config) logger.Interface {
	level := logger.Warn
	switch strings.ToLower(cfg.LogLevel) {
	case "silent":
		level = logger.Silent
	case "error":
		level = logger.Error
	case "warn", "warning", "":
		level = logger.Warn
	case "info", "debug":
		level = logger.Info
	}
	gormCfg := logger.Config{
		SlowThreshold:             cfg.SlowThreshold,
		LogLevel:                  level,
		IgnoreRecordNotFoundError: true,
		Colorful:                  false,
	}
	if gormCfg.SlowThreshold <= 0 {
		gormCfg.SlowThreshold = 200 * time.Millisecond
	}
	return logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), gormCfg)
}

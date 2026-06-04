package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	handyconfig "github.com/FiyZou/handygo/config"
	"github.com/FiyZou/handygo/database"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

type ModelGenerateOptions struct {
	ConfigPath string
}

type demoUserStatus string

type demoUser struct {
	ID           uint           `gorm:"primaryKey"`
	Email        string         `gorm:"size:191;uniqueIndex;not null"`
	Name         string         `gorm:"size:64;not null"`
	PasswordHash string         `gorm:"size:255;not null"`
	Status       demoUserStatus `gorm:"size:16;not null;default:enabled"`
	Roles        []demoRole     `gorm:"many2many:user_roles;"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (*demoUser) TableName() string {
	return "users"
}

type demoRole struct {
	ID          uint             `gorm:"primaryKey"`
	Code        string           `gorm:"size:64;uniqueIndex;not null"`
	Name        string           `gorm:"size:64;not null"`
	Description string           `gorm:"size:255"`
	Permissions []demoPermission `gorm:"many2many:role_permissions;"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (*demoRole) TableName() string {
	return "roles"
}

type demoPermission struct {
	ID          uint   `gorm:"primaryKey"`
	Code        string `gorm:"size:128;uniqueIndex;not null"`
	Name        string `gorm:"size:64;not null"`
	Group       string `gorm:"size:64;not null"`
	Description string `gorm:"size:255"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (*demoPermission) TableName() string {
	return "permissions"
}

type demoUserRole struct {
	UserID uint `gorm:"primaryKey"`
	RoleID uint `gorm:"primaryKey"`
}

func (*demoUserRole) TableName() string {
	return "user_roles"
}

type demoRolePermission struct {
	RoleID       uint `gorm:"primaryKey"`
	PermissionID uint `gorm:"primaryKey"`
}

func (*demoRolePermission) TableName() string {
	return "role_permissions"
}

type modelGeneratorConfig struct {
	Database database.Config     `mapstructure:"database"`
	Generate modelGenerateConfig `mapstructure:"generate"`
}

type modelGenerateConfig struct {
	ModelPath       string `mapstructure:"modelPath"`
	AutoMigrateDemo bool   `mapstructure:"autoMigrateDemo"`
	Nullable        bool   `mapstructure:"nullable"`
	Coverable       bool   `mapstructure:"coverable"`
	Signable        bool   `mapstructure:"signable"`
	IndexTag        bool   `mapstructure:"indexTag"`
	TypeTag         bool   `mapstructure:"typeTag"`
	DefaultTag      bool   `mapstructure:"defaultTag"`
}

func GenerateModels(opts ModelGenerateOptions) error {
	configPath, err := resolveConfigPath(opts.ConfigPath)
	if err != nil {
		return err
	}

	projectRoot, configPath, err := resolveProjectRoot(configPath)
	if err != nil {
		return err
	}
	if err := os.Chdir(projectRoot); err != nil {
		return err
	}

	cfg, err := loadModelGeneratorConfig(configPath)
	if err != nil {
		return fmt.Errorf("load generator config: %w", err)
	}
	normalizeModelGeneratePaths(projectRoot, &cfg)

	db, err := openGeneratorDatabase(cfg.Database)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	if cfg.Generate.AutoMigrateDemo {
		if err := autoMigrateDemo(db); err != nil {
			return fmt.Errorf("auto migrate demo schema: %w", err)
		}
	}
	if err := cleanupGenerated(cfg.Generate.ModelPath); err != nil {
		return fmt.Errorf("cleanup generated files: %w", err)
	}

	tempQueryPath, err := os.MkdirTemp("", "handygo-gen-query-*")
	if err != nil {
		return fmt.Errorf("create temporary query path: %w", err)
	}
	defer os.RemoveAll(tempQueryPath)

	g := gen.NewGenerator(gen.Config{
		OutPath:             tempQueryPath,
		ModelPkgPath:        cfg.Generate.ModelPath,
		Incremental:         true,
		FieldNullable:       cfg.Generate.Nullable,
		FieldCoverable:      cfg.Generate.Coverable,
		FieldSignable:       cfg.Generate.Signable,
		FieldWithIndexTag:   cfg.Generate.IndexTag,
		FieldWithTypeTag:    cfg.Generate.TypeTag,
		FieldWithDefaultTag: cfg.Generate.DefaultTag,
	})
	g.UseDB(db)

	models, err := generateDatabaseModels(g, db)
	if err != nil {
		return err
	}
	if len(models) == 0 {
		return fmt.Errorf("no tables found in database; check %s database.dsn", configPath)
	}
	g.Execute()
	return nil
}

func resolveConfigPath(configPath string) (string, error) {
	if configPath != "" {
		return configPath, nil
	}
	for _, candidate := range []string{"manifest/gen.yaml", "examples/manifest/gen.yaml"} {
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("missing config file; use -c manifest/gen.yaml")
}

func resolveProjectRoot(configPath string) (string, string, error) {
	absConfig, err := filepath.Abs(configPath)
	if err != nil {
		return "", "", err
	}
	if _, err := os.Stat(absConfig); err != nil {
		return "", "", err
	}

	dir := filepath.Dir(absConfig)
	for {
		if filepath.Base(dir) == "manifest" {
			root := filepath.Dir(dir)
			return root, absConfig, nil
		}
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, absConfig, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", "", fmt.Errorf("project root not found for config: %s", configPath)
		}
		dir = parent
	}
}

func normalizeModelGeneratePaths(projectRoot string, cfg *modelGeneratorConfig) {
	cfg.Generate.ModelPath = absPath(projectRoot, cfg.Generate.ModelPath)
}

func absPath(root string, path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(root, path)
}

func cleanupGenerated(paths ...string) error {
	for _, path := range paths {
		if path == "" {
			continue
		}
		matches, err := filepath.Glob(filepath.Join(path, "*.gen.go"))
		if err != nil {
			return err
		}
		matches = append(matches, filepath.Join(path, ".genmanifest.json"))
		for _, match := range matches {
			if err := os.Remove(match); err != nil && !os.IsNotExist(err) {
				return err
			}
		}
	}
	return nil
}

func loadModelGeneratorConfig(path string) (modelGeneratorConfig, error) {
	cfg := modelGeneratorConfig{
		Generate: modelGenerateConfig{
			ModelPath:       "internal/model",
			AutoMigrateDemo: true,
			Signable:        true,
			IndexTag:        true,
			TypeTag:         true,
			DefaultTag:      true,
		},
	}
	if err := handyconfig.Load(path, &cfg); err != nil {
		return modelGeneratorConfig{}, err
	}
	if err := cfg.validate(); err != nil {
		return modelGeneratorConfig{}, err
	}
	return cfg, nil
}

func (cfg modelGeneratorConfig) validate() error {
	if strings.TrimSpace(cfg.Database.Driver) == "" {
		return fmt.Errorf("database.driver cannot be empty")
	}
	if strings.TrimSpace(cfg.Database.DSN) == "" {
		return fmt.Errorf("database.dsn cannot be empty")
	}
	if strings.TrimSpace(cfg.Generate.ModelPath) == "" {
		return fmt.Errorf("generate.modelPath cannot be empty")
	}
	return nil
}

func openGeneratorDatabase(cfg database.Config) (*gorm.DB, error) {
	var dialector gorm.Dialector
	switch strings.ToLower(strings.TrimSpace(cfg.Driver)) {
	case "mysql":
		dialector = mysql.Open(cfg.DSN)
	case "postgres", "postgresql":
		dialector = postgres.Open(cfg.DSN)
	case "sqlite", "sqlite3":
		dialector = sqlite.Open(cfg.DSN)
	case "sqlserver", "mssql":
		dialector = sqlserver.Open(cfg.DSN)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}
	return gorm.Open(dialector, &gorm.Config{})
}

func autoMigrateDemo(db *gorm.DB) error {
	return db.AutoMigrate(
		&demoUser{},
		&demoRole{},
		&demoPermission{},
		&demoUserRole{},
		&demoRolePermission{},
	)
}

func generateDatabaseModels(g *gen.Generator, db *gorm.DB) ([]interface{}, error) {
	tables, err := db.Migrator().GetTables()
	if err != nil {
		return nil, fmt.Errorf("get tables: %w", err)
	}
	tableSet := make(map[string]bool, len(tables))
	for _, table := range tables {
		if strings.HasPrefix(table, "sqlite_") {
			continue
		}
		tableSet[table] = true
	}

	models := make([]interface{}, 0, len(tableSet))
	generated := make(map[string]bool, len(tableSet))
	if tableSet["users"] {
		userOpts := append(commonModelOpts(),
			gen.FieldType("status", "UserStatus"),
			gen.FieldJSONTag("password_hash", "-"),
		)

		if tableSet["roles"] && tableSet["user_roles"] {
			roleOpts := commonModelOpts()
			if tableSet["permissions"] && tableSet["role_permissions"] {
				permission := g.GenerateModel("permissions", commonModelOpts()...)
				models = append(models, permission)
				generated["permissions"] = true

				roleOpts = append(roleOpts, gen.FieldRelate(field.Many2Many, "Permissions", permission, &field.RelateConfig{
					RelateSlice: true,
					JSONTag:     "permissions,omitempty",
					GORMTag:     field.GormTag{"many2many": []string{"role_permissions;"}},
				}))
			}

			role := g.GenerateModel("roles", roleOpts...)
			models = append(models, role)
			generated["roles"] = true

			userOpts = append(userOpts, gen.FieldRelate(field.Many2Many, "Roles", role, &field.RelateConfig{
				RelateSlice: true,
				JSONTag:     "roles,omitempty",
				GORMTag:     field.GormTag{"many2many": []string{"user_roles;"}},
			}))
		}

		user := g.GenerateModel("users", userOpts...)
		models = append(models, user)
		generated["users"] = true
	}

	for _, table := range tables {
		if strings.HasPrefix(table, "sqlite_") {
			continue
		}
		if generated[table] {
			continue
		}
		models = append(models, g.GenerateModel(table, commonModelOpts()...))
	}
	return models, nil
}

func commonModelOpts() []gen.ModelOpt {
	return []gen.ModelOpt{
		gen.FieldTypeReg("(^id$|_id$)", "uint"),
		gen.FieldJSONTagWithNS(func(columnName string) string {
			switch columnName {
			case "created_at":
				return "createdAt"
			case "updated_at":
				return "updatedAt"
			default:
				return columnName
			}
		}),
	}
}

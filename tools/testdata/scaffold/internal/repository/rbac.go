package repository

import (
	"context"

	"github.com/FiyZou/handygo/examples/internal/model"
	"gorm.io/gorm"
)

type RBACRepository struct {
	db *gorm.DB
}

func NewRBACRepository(db *gorm.DB) *RBACRepository {
	return &RBACRepository{db: db}
}

func (r *RBACRepository) WithDB(db *gorm.DB) *RBACRepository {
	if db == nil {
		return r
	}
	return &RBACRepository{db: db}
}

func (r *RBACRepository) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return r.db.WithContext(ctx).Transaction(fn)
}

func (r *RBACRepository) ListRoles(ctx context.Context) ([]model.Role, error) {
	var roles []model.Role
	err := r.db.WithContext(ctx).Preload("Permissions").Order("id asc").Find(&roles).Error
	return roles, err
}

func (r *RBACRepository) ListPermissions(ctx context.Context) ([]model.Permission, error) {
	var permissions []model.Permission
	err := r.db.WithContext(ctx).Order("id asc").Find(&permissions).Error
	return permissions, err
}

func (r *RBACRepository) FindRoleByCode(ctx context.Context, code string) (*model.Role, error) {
	var role model.Role
	err := r.db.WithContext(ctx).Preload("Permissions").Where("code = ?", code).First(&role).Error
	return &role, err
}

func (r *RBACRepository) FindRolesByIDs(ctx context.Context, ids []uint) ([]model.Role, error) {
	var roles []model.Role
	err := r.db.WithContext(ctx).Where("id in ?", ids).Find(&roles).Error
	return roles, err
}

func (r *RBACRepository) FindPermissionsByIDs(ctx context.Context, ids []uint) ([]model.Permission, error) {
	var permissions []model.Permission
	err := r.db.WithContext(ctx).Where("id in ?", ids).Find(&permissions).Error
	return permissions, err
}

func (r *RBACRepository) CreateRole(ctx context.Context, role *model.Role) error {
	return r.db.WithContext(ctx).Create(role).Error
}

func (r *RBACRepository) ReplaceRolePermissions(ctx context.Context, role *model.Role, permissions []model.Permission) error {
	return r.db.WithContext(ctx).Model(role).Association("Permissions").Replace(permissions)
}

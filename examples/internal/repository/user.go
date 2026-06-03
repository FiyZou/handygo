package repository

import (
	"context"

	"github.com/FiyZou/handygo/examples/internal/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *UserRepository) FindByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Preload("Roles").Preload("Roles.Permissions").First(&user, id).Error
	return &user, err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Preload("Roles").Preload("Roles.Permissions").Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *UserRepository) List(ctx context.Context, page, size int) ([]model.User, int64, error) {
	var users []model.User
	var total int64
	query := r.db.WithContext(ctx).Model(&model.User{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := query.Preload("Roles").Order("id desc").Limit(size).Offset((page - 1) * size).Find(&users).Error
	return users, total, err
}

func (r *UserRepository) Save(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *UserRepository) ReplaceRoles(ctx context.Context, user *model.User, roles []model.Role) error {
	return r.db.WithContext(ctx).Model(user).Association("Roles").Replace(roles)
}

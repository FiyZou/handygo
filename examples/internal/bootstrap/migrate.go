package bootstrap

import (
	"errors"

	"github.com/FiyZou/handygo/examples/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func migrateAndSeed(db *gorm.DB, passwordCost int) error {
	if err := db.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.Permission{},
		&model.UserRole{},
		&model.RolePermission{},
	); err != nil {
		return err
	}

	permissions, err := seedPermissions(db)
	if err != nil {
		return err
	}
	adminRole, err := seedAdminRole(db, permissions)
	if err != nil {
		return err
	}
	return seedAdminUser(db, adminRole, passwordCost)
}

func seedPermissions(db *gorm.DB) ([]model.Permission, error) {
	seeds := []model.Permission{
		{Code: "backend:access", Name: "后台访问", Group: "system", Description: "允许登录和访问后台接口"},
		{Code: "user:list", Name: "用户列表", Group: "user", Description: "查看后台用户列表"},
		{Code: "user:create", Name: "创建用户", Group: "user", Description: "创建后台用户"},
		{Code: "user:update", Name: "更新用户", Group: "user", Description: "更新后台用户和角色绑定"},
		{Code: "role:list", Name: "角色列表", Group: "rbac", Description: "查看角色列表"},
		{Code: "role:create", Name: "创建角色", Group: "rbac", Description: "创建角色并绑定权限"},
		{Code: "permission:list", Name: "权限列表", Group: "rbac", Description: "查看权限列表"},
	}
	for _, seed := range seeds {
		var permission model.Permission
		err := db.Where("code = ?", seed.Code).First(&permission).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := db.Create(&seed).Error; err != nil {
				return nil, err
			}
			continue
		}
		if err != nil {
			return nil, err
		}
	}

	var permissions []model.Permission
	if err := db.Order("id asc").Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

func seedAdminRole(db *gorm.DB, permissions []model.Permission) (*model.Role, error) {
	var role model.Role
	err := db.Where("code = ?", "admin").First(&role).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		role = model.Role{Code: "admin", Name: "管理员", Description: "系统默认管理员角色"}
		if err := db.Create(&role).Error; err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	if err := db.Model(&role).Association("Permissions").Replace(permissions); err != nil {
		return nil, err
	}
	return &role, nil
}

func seedAdminUser(db *gorm.DB, adminRole *model.Role, passwordCost int) error {
	if passwordCost <= 0 {
		passwordCost = bcrypt.DefaultCost
	}
	var user model.User
	err := db.Where("email = ?", "admin@example.com").First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		hash, err := bcrypt.GenerateFromPassword([]byte("admin123456"), passwordCost)
		if err != nil {
			return err
		}
		user = model.User{
			Email:        "admin@example.com",
			Name:         "Administrator",
			PasswordHash: string(hash),
			Status:       model.UserStatusEnabled,
		}
		if err := db.Create(&user).Error; err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return db.Model(&user).Association("Roles").Replace([]model.Role{*adminRole})
}

package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	rbacv1 "github.com/FiyZou/handygo/examples/api/rbac/v1"
	userv1 "github.com/FiyZou/handygo/examples/api/user/v1"
	"github.com/FiyZou/handygo/examples/internal/model"
	"github.com/FiyZou/handygo/examples/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestUserServiceCreateRollsBackWhenRoleMissing(t *testing.T) {
	db := newServiceTestDB(t)
	role := seedRole(t, db, "admin")
	svc := NewUserService(bcrypt.MinCost, repository.NewUserRepository(db), repository.NewRBACRepository(db))

	_, err := svc.Create(context.Background(), userv1.CreateReq{
		Email:    "new@example.com",
		Name:     "New User",
		Password: "password123",
		RoleIDs:  []uint{role.ID, role.ID + 999},
	})
	if !errors.Is(err, ErrRoleNotFound) {
		t.Fatalf("error = %v, want ErrRoleNotFound", err)
	}

	var total int64
	if err := db.Model(&model.User{}).Count(&total).Error; err != nil {
		t.Fatalf("count users: %v", err)
	}
	if total != 0 {
		t.Fatalf("users count = %d, want rollback to zero", total)
	}
}

func TestUserServiceUpdateRollsBackWhenRoleMissing(t *testing.T) {
	db := newServiceTestDB(t)
	oldRole := seedRole(t, db, "old")
	newRole := seedRole(t, db, "new")
	user := seedUser(t, db, "old@example.com", oldRole)
	svc := NewUserService(bcrypt.MinCost, repository.NewUserRepository(db), repository.NewRBACRepository(db))

	_, err := svc.Update(context.Background(), user.ID, userv1.UpdateReq{
		Name:    "Changed",
		Status:  string(model.UserStatusDisabled),
		RoleIDs: []uint{newRole.ID, newRole.ID + 999},
	})
	if !errors.Is(err, ErrRoleNotFound) {
		t.Fatalf("error = %v, want ErrRoleNotFound", err)
	}

	var got model.User
	if err := db.Preload("Roles").First(&got, user.ID).Error; err != nil {
		t.Fatalf("find user: %v", err)
	}
	if got.Name != user.Name || got.Status != user.Status {
		t.Fatalf("user = (%q, %q), want rollback to (%q, %q)", got.Name, got.Status, user.Name, user.Status)
	}
	if len(got.Roles) != 1 || got.Roles[0].ID != oldRole.ID {
		t.Fatalf("roles = %+v, want original role %d", got.Roles, oldRole.ID)
	}
}

func TestRBACServiceCreateRoleRollsBackWhenPermissionMissing(t *testing.T) {
	db := newServiceTestDB(t)
	permission := seedPermission(t, db, "backend:access")
	svc := NewRBACService(repository.NewRBACRepository(db))

	_, err := svc.CreateRole(context.Background(), rbacv1.CreateRoleReq{
		Code:          "operator",
		Name:          "Operator",
		PermissionIDs: []uint{permission.ID, permission.ID + 999},
	})
	if !errors.Is(err, ErrPermissionNotFound) {
		t.Fatalf("error = %v, want ErrPermissionNotFound", err)
	}

	var total int64
	if err := db.Model(&model.Role{}).Count(&total).Error; err != nil {
		t.Fatalf("count roles: %v", err)
	}
	if total != 0 {
		t.Fatalf("roles count = %d, want rollback to zero", total)
	}
}

func newServiceTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", strings.NewReplacer("/", "_", " ", "_").Replace(t.Name()))
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&model.User{}, &model.Role{}, &model.Permission{}, &model.UserRole{}, &model.RolePermission{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}
	return db
}

func seedRole(t *testing.T, db *gorm.DB, code string) model.Role {
	t.Helper()

	role := model.Role{Code: code, Name: code}
	if err := db.Create(&role).Error; err != nil {
		t.Fatalf("seed role: %v", err)
	}
	return role
}

func seedPermission(t *testing.T, db *gorm.DB, code string) model.Permission {
	t.Helper()

	permission := model.Permission{Code: code, Name: code, Group: "backend"}
	if err := db.Create(&permission).Error; err != nil {
		t.Fatalf("seed permission: %v", err)
	}
	return permission
}

func seedUser(t *testing.T, db *gorm.DB, email string, roles ...model.Role) model.User {
	t.Helper()

	hash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	user := model.User{
		Email:        email,
		Name:         "Original",
		PasswordHash: string(hash),
		Status:       model.UserStatusEnabled,
	}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}
	if len(roles) > 0 {
		if err := db.Model(&user).Association("Roles").Replace(roles); err != nil {
			t.Fatalf("seed user roles: %v", err)
		}
	}
	return user
}

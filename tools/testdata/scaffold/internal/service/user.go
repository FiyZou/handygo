package service

import (
	"context"
	"errors"
	"fmt"

	rbacv1 "github.com/FiyZou/handygo/examples/api/rbac/v1"
	userv1 "github.com/FiyZou/handygo/examples/api/user/v1"
	"github.com/FiyZou/handygo/examples/internal/model"
	"github.com/FiyZou/handygo/examples/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var ErrRoleNotFound = errors.New("one or more roles were not found")

type UserService struct {
	passwordCost int
	users        *repository.UserRepository
	rbac         *repository.RBACRepository
}

func NewUserService(passwordCost int, users *repository.UserRepository, rbac *repository.RBACRepository) *UserService {
	if passwordCost <= 0 {
		passwordCost = bcrypt.DefaultCost
	}
	return &UserService{passwordCost: passwordCost, users: users, rbac: rbac}
}

func (s *UserService) List(ctx context.Context, page, size int) (userv1.ListRes, int64, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 100 {
		size = 20
	}
	users, total, err := s.users.List(ctx, page, size)
	if err != nil {
		return nil, 0, err
	}
	items := make([]userv1.UserRes, 0, len(users))
	for _, user := range users {
		roleItems := make([]rbacv1.RoleRes, 0, len(user.Roles))
		for _, role := range user.Roles {
			permissionItems := make([]rbacv1.PermissionRes, 0, len(role.Permissions))
			for _, permission := range role.Permissions {
				permissionItems = append(permissionItems, rbacv1.PermissionRes{
					ID:          permission.ID,
					Code:        permission.Code,
					Name:        permission.Name,
					Group:       permission.Group,
					Description: permission.Description,
					CreatedAt:   permission.CreatedAt,
					UpdatedAt:   permission.UpdatedAt,
				})
			}
			roleItems = append(roleItems, rbacv1.RoleRes{
				ID:          role.ID,
				Code:        role.Code,
				Name:        role.Name,
				Description: role.Description,
				Permissions: permissionItems,
				CreatedAt:   role.CreatedAt,
				UpdatedAt:   role.UpdatedAt,
			})
		}
		items = append(items, userv1.UserRes{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			Status:    string(user.Status),
			Roles:     roleItems,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}
	return items, total, nil
}

func (s *UserService) Create(ctx context.Context, req userv1.CreateReq) (userv1.CreateRes, error) {
	var userID uint
	if err := s.users.Transaction(ctx, func(tx *gorm.DB) error {
		users := s.users.WithDB(tx)
		rbac := s.rbac.WithDB(tx)

		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), s.passwordCost)
		if err != nil {
			return err
		}
		status := model.UserStatusEnabled
		if req.Status != "" {
			status = model.UserStatus(req.Status)
		}
		user := &model.User{
			Email:        req.Email,
			Name:         req.Name,
			PasswordHash: string(hash),
			Status:       status,
		}
		if err := users.Create(ctx, user); err != nil {
			return err
		}
		if len(req.RoleIDs) > 0 {
			roles, err := rbac.FindRolesByIDs(ctx, req.RoleIDs)
			if err != nil {
				return err
			}
			if len(roles) != len(uniqueUint(req.RoleIDs)) {
				return fmt.Errorf("%w: roleIds", ErrRoleNotFound)
			}
			if err := users.ReplaceRoles(ctx, user, roles); err != nil {
				return err
			}
		}
		userID = user.ID
		return nil
	}); err != nil {
		return userv1.CreateRes{}, err
	}
	created, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return userv1.CreateRes{}, err
	}
	roleItems := userRoleItems(created)
	return userv1.CreateRes{
		ID:        created.ID,
		Email:     created.Email,
		Name:      created.Name,
		Status:    string(created.Status),
		Roles:     roleItems,
		CreatedAt: created.CreatedAt,
		UpdatedAt: created.UpdatedAt,
	}, nil
}

func (s *UserService) Update(ctx context.Context, id uint, req userv1.UpdateReq) (userv1.UpdateRes, error) {
	if err := s.users.Transaction(ctx, func(tx *gorm.DB) error {
		users := s.users.WithDB(tx)
		rbac := s.rbac.WithDB(tx)

		user, err := users.FindByID(ctx, id)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return gorm.ErrRecordNotFound
		}
		if err != nil {
			return err
		}
		user.Name = req.Name
		user.Status = model.UserStatus(req.Status)
		if err := users.Save(ctx, user); err != nil {
			return err
		}
		roles, err := rbac.FindRolesByIDs(ctx, req.RoleIDs)
		if err != nil {
			return err
		}
		if len(roles) != len(uniqueUint(req.RoleIDs)) {
			return fmt.Errorf("%w: roleIds", ErrRoleNotFound)
		}
		if err := users.ReplaceRoles(ctx, user, roles); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return userv1.UpdateRes{}, err
	}
	updated, err := s.users.FindByID(ctx, id)
	if err != nil {
		return userv1.UpdateRes{}, err
	}
	roleItems := userRoleItems(updated)
	return userv1.UpdateRes{
		ID:        updated.ID,
		Email:     updated.Email,
		Name:      updated.Name,
		Status:    string(updated.Status),
		Roles:     roleItems,
		CreatedAt: updated.CreatedAt,
		UpdatedAt: updated.UpdatedAt,
	}, nil
}

func userRoleItems(user *model.User) []rbacv1.RoleRes {
	roleItems := make([]rbacv1.RoleRes, 0, len(user.Roles))
	for _, role := range user.Roles {
		permissionItems := make([]rbacv1.PermissionRes, 0, len(role.Permissions))
		for _, permission := range role.Permissions {
			permissionItems = append(permissionItems, rbacv1.PermissionRes{
				ID:          permission.ID,
				Code:        permission.Code,
				Name:        permission.Name,
				Group:       permission.Group,
				Description: permission.Description,
				CreatedAt:   permission.CreatedAt,
				UpdatedAt:   permission.UpdatedAt,
			})
		}
		roleItems = append(roleItems, rbacv1.RoleRes{
			ID:          role.ID,
			Code:        role.Code,
			Name:        role.Name,
			Description: role.Description,
			Permissions: permissionItems,
			CreatedAt:   role.CreatedAt,
			UpdatedAt:   role.UpdatedAt,
		})
	}
	return roleItems
}

func uniqueUint(values []uint) map[uint]struct{} {
	seen := make(map[uint]struct{}, len(values))
	for _, value := range values {
		seen[value] = struct{}{}
	}
	return seen
}

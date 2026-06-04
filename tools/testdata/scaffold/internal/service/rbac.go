package service

import (
	"context"
	"errors"
	"fmt"

	rbacv1 "github.com/FiyZou/handygo/examples/api/rbac/v1"
	"github.com/FiyZou/handygo/examples/internal/model"
	"github.com/FiyZou/handygo/examples/internal/repository"
	"gorm.io/gorm"
)

var ErrPermissionNotFound = errors.New("one or more permissions were not found")

type RBACService struct {
	rbac *repository.RBACRepository
}

func NewRBACService(rbac *repository.RBACRepository) *RBACService {
	return &RBACService{rbac: rbac}
}

func (s *RBACService) ListRoles(ctx context.Context) (rbacv1.ListRolesRes, error) {
	roles, err := s.rbac.ListRoles(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]rbacv1.RoleRes, 0, len(roles))
	for _, role := range roles {
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
		items = append(items, rbacv1.RoleRes{
			ID:          role.ID,
			Code:        role.Code,
			Name:        role.Name,
			Description: role.Description,
			Permissions: permissionItems,
			CreatedAt:   role.CreatedAt,
			UpdatedAt:   role.UpdatedAt,
		})
	}
	return items, nil
}

func (s *RBACService) ListPermissions(ctx context.Context) (rbacv1.ListPermissionsRes, error) {
	permissions, err := s.rbac.ListPermissions(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]rbacv1.PermissionRes, 0, len(permissions))
	for _, permission := range permissions {
		items = append(items, rbacv1.PermissionRes{
			ID:          permission.ID,
			Code:        permission.Code,
			Name:        permission.Name,
			Group:       permission.Group,
			Description: permission.Description,
			CreatedAt:   permission.CreatedAt,
			UpdatedAt:   permission.UpdatedAt,
		})
	}
	return items, nil
}

func (s *RBACService) CreateRole(ctx context.Context, req rbacv1.CreateRoleReq) (rbacv1.CreateRoleRes, error) {
	var roleCode string
	if err := s.rbac.Transaction(ctx, func(tx *gorm.DB) error {
		rbac := s.rbac.WithDB(tx)
		role := &model.Role{
			Code:        req.Code,
			Name:        req.Name,
			Description: req.Description,
		}
		if err := rbac.CreateRole(ctx, role); err != nil {
			return err
		}
		if len(req.PermissionIDs) > 0 {
			permissions, err := rbac.FindPermissionsByIDs(ctx, req.PermissionIDs)
			if err != nil {
				return err
			}
			if len(permissions) != len(uniqueUint(req.PermissionIDs)) {
				return fmt.Errorf("%w: permissionIds", ErrPermissionNotFound)
			}
			if err := rbac.ReplaceRolePermissions(ctx, role, permissions); err != nil {
				return err
			}
		}
		roleCode = role.Code
		return nil
	}); err != nil {
		return rbacv1.CreateRoleRes{}, err
	}
	created, err := s.rbac.FindRoleByCode(ctx, roleCode)
	if err != nil {
		return rbacv1.CreateRoleRes{}, err
	}
	permissionItems := rolePermissionItems(created)
	return rbacv1.CreateRoleRes{
		ID:          created.ID,
		Code:        created.Code,
		Name:        created.Name,
		Description: created.Description,
		Permissions: permissionItems,
		CreatedAt:   created.CreatedAt,
		UpdatedAt:   created.UpdatedAt,
	}, nil
}

func rolePermissionItems(role *model.Role) []rbacv1.PermissionRes {
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
	return permissionItems
}

package v1

import "time"

type PermissionRes struct {
	ID          uint      `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Group       string    `json:"group"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type RoleRes struct {
	ID          uint            `json:"id"`
	Code        string          `json:"code"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Permissions []PermissionRes `json:"permissions,omitempty"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}

type CreateRoleReq struct {
	Code          string `json:"code" binding:"required" validate:"required"`
	Name          string `json:"name" binding:"required" validate:"required"`
	Description   string `json:"description"`
	PermissionIDs []uint `json:"permissionIds" binding:"omitempty" validate:"omitempty"`
}

type CreateRoleRes = RoleRes
type ListRolesRes = []RoleRes
type ListPermissionsRes = []PermissionRes

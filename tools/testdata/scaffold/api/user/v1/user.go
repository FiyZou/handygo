package v1

import (
	"time"

	rbacv1 "github.com/FiyZou/handygo/examples/api/rbac/v1"
)

type UserRes struct {
	ID        uint             `json:"id"`
	Email     string           `json:"email"`
	Name      string           `json:"name"`
	Status    string           `json:"status"`
	Roles     []rbacv1.RoleRes `json:"roles,omitempty"`
	CreatedAt time.Time        `json:"createdAt"`
	UpdatedAt time.Time        `json:"updatedAt"`
}

type CreateReq struct {
	Email    string `json:"email" binding:"required,email" validate:"required,email"`
	Name     string `json:"name" binding:"required" validate:"required"`
	Password string `json:"password" binding:"required,min=8" validate:"required,min=8"`
	Status   string `json:"status" binding:"omitempty,oneof=enabled disabled" validate:"omitempty,oneof=enabled disabled"`
	RoleIDs  []uint `json:"roleIds" binding:"omitempty" validate:"omitempty"`
}

type CreateRes = UserRes

type UpdateReq struct {
	Name    string `json:"name" binding:"required" validate:"required"`
	Status  string `json:"status" binding:"required,oneof=enabled disabled" validate:"required,oneof=enabled disabled"`
	RoleIDs []uint `json:"roleIds" binding:"omitempty" validate:"omitempty"`
}

type UpdateRes = UserRes

type MeRes = UserRes
type ListRes = []UserRes

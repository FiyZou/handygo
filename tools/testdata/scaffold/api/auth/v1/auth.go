package v1

import (
	"time"

	userv1 "github.com/FiyZou/handygo/examples/api/user/v1"
)

type RegisterReq struct {
	Email    string `json:"email" binding:"required,email" validate:"required,email"`
	Name     string `json:"name" binding:"required" validate:"required"`
	Password string `json:"password" binding:"required,min=8" validate:"required,min=8"`
}

type RegisterRes = userv1.UserRes

type LoginReq struct {
	Email    string `json:"email" binding:"required,email" validate:"required,email"`
	Password string `json:"password" binding:"required" validate:"required"`
}

type LoginRes struct {
	Token       string         `json:"token"`
	ExpiresAt   time.Time      `json:"expiresAt"`
	User        userv1.UserRes `json:"user"`
	Roles       []string       `json:"roles"`
	Permissions []string       `json:"permissions"`
}

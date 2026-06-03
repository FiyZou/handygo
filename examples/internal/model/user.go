package model

import "time"

type UserStatus string

const (
	UserStatusEnabled  UserStatus = "enabled"
	UserStatusDisabled UserStatus = "disabled"
)

type User struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	Email        string     `gorm:"size:191;uniqueIndex;not null" json:"email"`
	Name         string     `gorm:"size:64;not null" json:"name"`
	PasswordHash string     `gorm:"size:255;not null" json:"-"`
	Status       UserStatus `gorm:"size:16;not null;default:enabled" json:"status"`
	Roles        []Role     `gorm:"many2many:user_roles;" json:"roles,omitempty"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
}

type Role struct {
	ID          uint         `gorm:"primaryKey" json:"id"`
	Code        string       `gorm:"size:64;uniqueIndex;not null" json:"code"`
	Name        string       `gorm:"size:64;not null" json:"name"`
	Description string       `gorm:"size:255" json:"description"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
	CreatedAt   time.Time    `json:"createdAt"`
	UpdatedAt   time.Time    `json:"updatedAt"`
}

type Permission struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Code        string    `gorm:"size:128;uniqueIndex;not null" json:"code"`
	Name        string    `gorm:"size:64;not null" json:"name"`
	Group       string    `gorm:"size:64;not null" json:"group"`
	Description string    `gorm:"size:255" json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type UserRole struct {
	UserID uint `gorm:"primaryKey"`
	RoleID uint `gorm:"primaryKey"`
}

type RolePermission struct {
	RoleID       uint `gorm:"primaryKey"`
	PermissionID uint `gorm:"primaryKey"`
}

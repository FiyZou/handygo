package service

import (
	"context"
	"errors"
	"time"

	authv1 "github.com/FiyZou/handygo/examples/api/auth/v1"
	rbacv1 "github.com/FiyZou/handygo/examples/api/rbac/v1"
	userv1 "github.com/FiyZou/handygo/examples/api/user/v1"
	appconfig "github.com/FiyZou/handygo/examples/internal/config"
	"github.com/FiyZou/handygo/examples/internal/model"
	"github.com/FiyZou/handygo/examples/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrInvalidCredential = errors.New("invalid email or password")
	ErrUserDisabled      = errors.New("user is disabled")
	ErrForbidden         = errors.New("permission denied")
)

type AuthService struct {
	auth  appconfig.Auth
	users *repository.UserRepository
	rbac  *repository.RBACRepository
}

type TokenClaims struct {
	UserID uint `json:"userId"`
	jwt.RegisteredClaims
}

func NewAuthService(auth appconfig.Auth, users *repository.UserRepository, rbac *repository.RBACRepository) *AuthService {
	if auth.TokenTTL <= 0 {
		auth.TokenTTL = 24 * time.Hour
	}
	if auth.PasswordCost <= 0 {
		auth.PasswordCost = bcrypt.DefaultCost
	}
	return &AuthService{auth: auth, users: users, rbac: rbac}
}

func (s *AuthService) Register(ctx context.Context, req authv1.RegisterReq) (authv1.RegisterRes, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), s.auth.PasswordCost)
	if err != nil {
		return authv1.RegisterRes{}, err
	}
	user := &model.User{
		Email:        req.Email,
		Name:         req.Name,
		PasswordHash: string(hash),
		Status:       model.UserStatusEnabled,
	}
	if err := s.users.Create(ctx, user); err != nil {
		return authv1.RegisterRes{}, err
	}
	return authv1.RegisterRes{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Status:    string(user.Status),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req authv1.LoginReq, requireBackendAccess bool) (*authv1.LoginRes, error) {
	user, err := s.users.FindByEmail(ctx, req.Email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrInvalidCredential
	}
	if err != nil {
		return nil, err
	}
	if user.Status != model.UserStatusEnabled {
		return nil, ErrUserDisabled
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredential
	}
	permissions := userPermissionCodes(user)
	if requireBackendAccess && !contains(permissions, "backend:access") {
		return nil, ErrForbidden
	}
	expiresAt := time.Now().Add(s.auth.TokenTTL)
	token, err := s.signToken(user.ID, expiresAt)
	if err != nil {
		return nil, err
	}
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
	userRes := userv1.UserRes{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Status:    string(user.Status),
		Roles:     roleItems,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	return &authv1.LoginRes{
		Token:       token,
		ExpiresAt:   expiresAt,
		User:        userRes,
		Roles:       userRoleCodes(user),
		Permissions: permissions,
	}, nil
}

func (s *AuthService) ParseToken(token string) (*TokenClaims, error) {
	parsed, err := jwt.ParseWithClaims(token, &TokenClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(s.auth.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := parsed.Claims.(*TokenClaims)
	if !ok || !parsed.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

func (s *AuthService) Me(ctx context.Context, userID uint) (*model.User, error) {
	return s.users.FindByID(ctx, userID)
}

func (s *AuthService) MeRes(ctx context.Context, userID uint) (userv1.MeRes, error) {
	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return userv1.MeRes{}, err
	}
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
	return userv1.MeRes{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Status:    string(user.Status),
		Roles:     roleItems,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *AuthService) signToken(userID uint, expiresAt time.Time) (string, error) {
	claims := TokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.auth.JWTSecret))
}

func userRoleCodes(user *model.User) []string {
	codes := make([]string, 0, len(user.Roles))
	for _, role := range user.Roles {
		codes = append(codes, role.Code)
	}
	return codes
}

func userPermissionCodes(user *model.User) []string {
	seen := map[string]struct{}{}
	for _, role := range user.Roles {
		for _, permission := range role.Permissions {
			seen[permission.Code] = struct{}{}
		}
	}
	codes := make([]string, 0, len(seen))
	for code := range seen {
		codes = append(codes, code)
	}
	return codes
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

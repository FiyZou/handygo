package middleware

import (
	"net/http"
	"strings"

	"github.com/FiyZou/handygo/examples/internal/model"
	"github.com/FiyZou/handygo/examples/internal/service"
	"github.com/FiyZou/handygo/response"
	"github.com/gin-gonic/gin"
)

const CurrentUserKey = "currentUser"

func Auth(auth *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := bearerToken(c.GetHeader("Authorization"))
		if token == "" {
			response.Abort(c, http.StatusUnauthorized, response.CodeUnauthorized, "missing bearer token")
			return
		}
		claims, err := auth.ParseToken(token)
		if err != nil {
			response.Abort(c, http.StatusUnauthorized, response.CodeUnauthorized, "invalid token")
			return
		}
		user, err := auth.Me(c.Request.Context(), claims.UserID)
		if err != nil {
			response.Abort(c, http.StatusUnauthorized, response.CodeUnauthorized, "invalid user")
			return
		}
		if user.Status != model.UserStatusEnabled {
			response.Abort(c, http.StatusForbidden, response.CodeForbidden, "user is disabled")
			return
		}
		c.Set(CurrentUserKey, user)
		c.Next()
	}
}

func RequirePermission(code string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := CurrentUser(c)
		if !ok || !hasPermission(user, code) {
			response.Abort(c, http.StatusForbidden, response.CodeForbidden, "permission denied")
			return
		}
		c.Next()
	}
}

func CurrentUser(c *gin.Context) (*model.User, bool) {
	value, ok := c.Get(CurrentUserKey)
	if !ok {
		return nil, false
	}
	user, ok := value.(*model.User)
	return user, ok
}

func bearerToken(header string) string {
	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(header, prefix))
}

func hasPermission(user *model.User, code string) bool {
	for _, role := range user.Roles {
		for _, permission := range role.Permissions {
			if permission.Code == code {
				return true
			}
		}
	}
	return false
}

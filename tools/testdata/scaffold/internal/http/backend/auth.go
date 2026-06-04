package backend

import (
	"errors"
	"net/http"

	authv1 "github.com/FiyZou/handygo/examples/api/auth/v1"
	authmw "github.com/FiyZou/handygo/examples/internal/http/middleware"
	"github.com/FiyZou/handygo/examples/internal/service"
	"github.com/FiyZou/handygo/response"
	"github.com/FiyZou/handygo/validate"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	auth *service.AuthService
}

func NewAuthHandler(auth *service.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req authv1.LoginReq
	if err := validate.Bind(c, &req); err != nil {
		response.JSON(c, http.StatusBadRequest, response.CodeBadRequest, validate.ErrorMessage(err), validate.Format(err))
		return
	}
	res, err := h.auth.Login(c.Request.Context(), req, true)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, service.ErrInvalidCredential) {
			status = http.StatusUnauthorized
		}
		if errors.Is(err, service.ErrForbidden) {
			status = http.StatusForbidden
		}
		response.Fail(c, status, status, err.Error())
		return
	}
	response.OK(c, res)
}

func (h *AuthHandler) Me(c *gin.Context) {
	user, _ := authmw.CurrentUser(c)
	res, err := h.auth.MeRes(c.Request.Context(), user.ID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, response.CodeInternalError, err.Error())
		return
	}
	response.OK(c, res)
}

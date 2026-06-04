package backend

import (
	"net/http"

	rbacv1 "github.com/FiyZou/handygo/examples/api/rbac/v1"
	"github.com/FiyZou/handygo/examples/internal/service"
	"github.com/FiyZou/handygo/response"
	"github.com/FiyZou/handygo/validate"
	"github.com/gin-gonic/gin"
)

type RBACHandler struct {
	rbac *service.RBACService
}

func NewRBACHandler(rbac *service.RBACService) *RBACHandler {
	return &RBACHandler{rbac: rbac}
}

func (h *RBACHandler) Roles(c *gin.Context) {
	roles, err := h.rbac.ListRoles(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, response.CodeInternalError, err.Error())
		return
	}
	response.OK(c, roles)
}

func (h *RBACHandler) Permissions(c *gin.Context) {
	permissions, err := h.rbac.ListPermissions(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, response.CodeInternalError, err.Error())
		return
	}
	response.OK(c, permissions)
}

func (h *RBACHandler) CreateRole(c *gin.Context) {
	var req rbacv1.CreateRoleReq
	if err := validate.Bind(c, &req); err != nil {
		response.JSON(c, http.StatusBadRequest, response.CodeBadRequest, validate.ErrorMessage(err), validate.Format(err))
		return
	}
	res, err := h.rbac.CreateRole(c.Request.Context(), req)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}
	response.Created(c, res)
}

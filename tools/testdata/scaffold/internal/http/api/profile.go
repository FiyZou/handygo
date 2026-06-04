package api

import (
	"net/http"

	authmw "github.com/FiyZou/handygo/examples/internal/http/middleware"
	"github.com/FiyZou/handygo/examples/internal/service"
	"github.com/FiyZou/handygo/response"
	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	auth *service.AuthService
}

func NewProfileHandler(auth *service.AuthService) *ProfileHandler {
	return &ProfileHandler{auth: auth}
}

func (h *ProfileHandler) Me(c *gin.Context) {
	user, _ := authmw.CurrentUser(c)
	res, err := h.auth.MeRes(c.Request.Context(), user.ID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, response.CodeInternalError, err.Error())
		return
	}
	response.OK(c, res)
}

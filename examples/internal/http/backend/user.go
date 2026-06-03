package backend

import (
	"errors"
	"net/http"
	"strconv"

	userv1 "github.com/FiyZou/handygo/examples/api/user/v1"
	"github.com/FiyZou/handygo/examples/internal/service"
	"github.com/FiyZou/handygo/response"
	"github.com/FiyZou/handygo/validate"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	users *service.UserService
}

func NewUserHandler(users *service.UserService) *UserHandler {
	return &UserHandler{users: users}
}

func (h *UserHandler) List(c *gin.Context) {
	page := queryInt(c, "page", 1)
	size := queryInt(c, "size", 20)
	users, total, err := h.users.List(c.Request.Context(), page, size)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, response.CodeInternalError, err.Error())
		return
	}
	response.Page(c, users, total, page, size)
}

func (h *UserHandler) Create(c *gin.Context) {
	var req userv1.CreateReq
	if err := validate.Bind(c, &req); err != nil {
		response.JSON(c, http.StatusBadRequest, response.CodeBadRequest, validate.ErrorMessage(err), validate.Format(err))
		return
	}
	res, err := h.users.Create(c.Request.Context(), req)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}
	response.Created(c, res)
}

func (h *UserHandler) Update(c *gin.Context) {
	id, err := pathUint(c, "id")
	if err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeBadRequest, "invalid user id")
		return
	}
	var req userv1.UpdateReq
	if err := validate.Bind(c, &req); err != nil {
		response.JSON(c, http.StatusBadRequest, response.CodeBadRequest, validate.ErrorMessage(err), validate.Format(err))
		return
	}
	res, err := h.users.Update(c.Request.Context(), id, req)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(c, http.StatusNotFound, response.CodeNotFound, "user not found")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}
	response.OK(c, res)
}

func queryInt(c *gin.Context, key string, fallback int) int {
	value, err := strconv.Atoi(c.DefaultQuery(key, strconv.Itoa(fallback)))
	if err != nil {
		return fallback
	}
	return value
}

func pathUint(c *gin.Context, key string) (uint, error) {
	value, err := strconv.ParseUint(c.Param(key), 10, 64)
	return uint(value), err
}

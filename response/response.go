package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	CodeOK            = 0
	CodeBadRequest    = 400
	CodeUnauthorized  = 401
	CodeForbidden     = 403
	CodeNotFound      = 404
	CodeInternalError = 500
)

type Body struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type PageBody struct {
	List  any   `json:"list"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Size  int   `json:"size"`
}

func JSON(c *gin.Context, httpStatus int, code int, message string, data any) {
	c.JSON(httpStatus, Body{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func OK(c *gin.Context, data any) {
	JSON(c, http.StatusOK, CodeOK, "ok", data)
}

func Created(c *gin.Context, data any) {
	JSON(c, http.StatusCreated, CodeOK, "created", data)
}

func Fail(c *gin.Context, httpStatus int, code int, message string) {
	JSON(c, httpStatus, code, message, nil)
}

func Abort(c *gin.Context, httpStatus int, code int, message string) {
	Fail(c, httpStatus, code, message)
	c.Abort()
}

func Page(c *gin.Context, list any, total int64, page int, size int) {
	OK(c, PageBody{
		List:  list,
		Total: total,
		Page:  page,
		Size:  size,
	})
}

package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func TestServerRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)
	s := New(Config{Name: "test", Addr: ":0", Mode: gin.TestMode}, zap.NewNop())
	s.Register(func(r *gin.Engine) {
		r.GET("/ping", func(c *gin.Context) {
			c.String(http.StatusOK, "pong")
		})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	s.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
	if w.Body.String() != "pong" {
		t.Fatalf("body = %q", w.Body.String())
	}
}

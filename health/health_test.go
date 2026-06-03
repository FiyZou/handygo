package health

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestHealthHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := New(time.Second)
	h.Register(NewCheck("ok", func(ctx context.Context) error { return nil }))
	r := gin.New()
	r.GET("/health", h.Handler())

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestHealthUnhealthy(t *testing.T) {
	h := New(time.Second)
	h.Register(NewCheck("bad", func(ctx context.Context) error { return errors.New("down") }))
	report := h.Check(context.Background())
	if report.Healthy {
		t.Fatal("expected unhealthy report")
	}
}

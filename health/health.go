package health

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Checker interface {
	Name() string
	Check(ctx context.Context) error
}

type CheckFunc func(ctx context.Context) error

type NamedCheck struct {
	name string
	fn   CheckFunc
}

type Status struct {
	Name    string `json:"name"`
	Healthy bool   `json:"healthy"`
	Error   string `json:"error,omitempty"`
}

type Report struct {
	Healthy bool     `json:"healthy"`
	Checks  []Status `json:"checks"`
}

type Health struct {
	timeout time.Duration
	checks  []Checker
	mu      sync.RWMutex
}

func New(timeout time.Duration) *Health {
	if timeout <= 0 {
		timeout = 2 * time.Second
	}
	return &Health{timeout: timeout}
}

func NewCheck(name string, fn CheckFunc) Checker {
	return NamedCheck{name: name, fn: fn}
}

func (c NamedCheck) Name() string {
	return c.name
}

func (c NamedCheck) Check(ctx context.Context) error {
	return c.fn(ctx)
}

func (h *Health) Register(checks ...Checker) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.checks = append(h.checks, checks...)
}

func (h *Health) Check(ctx context.Context) Report {
	h.mu.RLock()
	checks := append([]Checker(nil), h.checks...)
	h.mu.RUnlock()

	report := Report{
		Healthy: true,
		Checks:  make([]Status, 0, len(checks)),
	}
	for _, check := range checks {
		checkCtx, cancel := context.WithTimeout(ctx, h.timeout)
		err := check.Check(checkCtx)
		cancel()

		status := Status{Name: check.Name(), Healthy: err == nil}
		if err != nil {
			status.Error = err.Error()
			report.Healthy = false
		}
		report.Checks = append(report.Checks, status)
	}
	return report
}

func (h *Health) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		report := h.Check(c.Request.Context())
		status := http.StatusOK
		if !report.Healthy {
			status = http.StatusServiceUnavailable
		}
		c.JSON(status, report)
	}
}

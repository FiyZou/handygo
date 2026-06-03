package middleware

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const RequestIDKey = "request_id"

type CORSConfig struct {
	AllowOrigins []string
	AllowMethods []string
	AllowHeaders []string
	MaxAge       time.Duration
}

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-Id")
		if requestID == "" {
			requestID = uuid.NewString()
		}
		c.Set(RequestIDKey, requestID)
		c.Header("X-Request-Id", requestID)
		c.Next()
	}
}

func AccessLog(logger *zap.Logger) gin.HandlerFunc {
	if logger == nil {
		logger = zap.NewNop()
	}
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		logger.Info("http request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.FullPath()),
			zap.String("raw_path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", time.Since(start)),
			zap.String("client_ip", c.ClientIP()),
			zap.String("request_id", requestID(c)),
		)
	}
}

func Recover(logger *zap.Logger) gin.HandlerFunc {
	if logger == nil {
		logger = zap.NewNop()
	}
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		logger.Error("panic recovered",
			zap.Any("panic", recovered),
			zap.String("path", c.Request.URL.Path),
			zap.String("request_id", requestID(c)),
		)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "internal server error",
		})
	})
}

func CORS(cfg CORSConfig) gin.HandlerFunc {
	origins := joinOrDefault(cfg.AllowOrigins, "*")
	methods := joinOrDefault(cfg.AllowMethods, "GET,POST,PUT,PATCH,DELETE,OPTIONS")
	headers := joinOrDefault(cfg.AllowHeaders, "Origin,Content-Type,Authorization,X-Request-Id")
	maxAge := "86400"
	if cfg.MaxAge > 0 {
		maxAge = formatSeconds(cfg.MaxAge)
	}
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", origins)
		c.Header("Access-Control-Allow-Methods", methods)
		c.Header("Access-Control-Allow-Headers", headers)
		c.Header("Access-Control-Max-Age", maxAge)
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func Timeout(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if timeout <= 0 {
			c.Next()
			return
		}
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func requestID(c *gin.Context) string {
	value, ok := c.Get(RequestIDKey)
	if !ok {
		return ""
	}
	requestID, _ := value.(string)
	return requestID
}

func joinOrDefault(values []string, fallback string) string {
	if len(values) == 0 {
		return fallback
	}
	result := values[0]
	for _, value := range values[1:] {
		result += "," + value
	}
	return result
}

func formatSeconds(duration time.Duration) string {
	return strconv.FormatInt(int64(duration.Seconds()), 10)
}

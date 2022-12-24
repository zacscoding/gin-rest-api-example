package middleware

import (
	"context"
	"gin-rest-api-example/pkg/logging"
	"gin-rest-api-example/pkg/trace"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	XRequestIdKey = "X-Request-ID" // request id header key
)

// RequestIDMiddleware attach request id and logger to context
// 1. extract request id from header if exist, otherwise generate
// 2. attach request id to logger and store it to context
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId := c.Request.Header.Get(XRequestIdKey)
		if requestId == "" {
			requestId = uuid.New().String()
		}

		c.Request = c.Request.WithContext(trace.WithRequestID(c, requestId))
		logger := logging.DefaultLogger().With("requestId", requestId)
		c.Request = c.Request.WithContext(logging.WithLogger(c, logger))
		c.Writer.Header().Set(XRequestIdKey, requestId)
	}
}

// TimeoutMiddleware attach deadline to gin.Request.Context
func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)

		defer func() {
			if ctx.Err() == context.DeadlineExceeded {
				c.AbortWithStatus(http.StatusGatewayTimeout)
			}
			cancel()
		}()
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// LoggingMiddleware use logging.DefaultLogger() i.e *zap.SugaredLogger with x-request-id
func LoggingMiddleware(skipPaths ...string) gin.HandlerFunc {
	skip := make(map[string]struct{}, len(skipPaths))
	for _, path := range skipPaths {
		skip[path] = struct{}{}
	}

	return func(c *gin.Context) {
		// skip logging
		if _, ok := skip[c.FullPath()]; ok {
			c.Next()
			return
		}

		start := time.Now()
		path := c.Request.URL.Path
		rawQuery := c.Request.URL.RawQuery

		// process request
		c.Next()

		logger := logging.FromContext(c.Request.Context())
		timestamp := time.Now()
		latency := timestamp.Sub(start)
		latencyValue := latency.String()
		clientIP := c.ClientIP()
		method := c.Request.Method
		status := c.Writer.Status()
		if rawQuery != "" {
			path = path + "?" + rawQuery
		}
		// append logger keys if not success or too slow latency.
		if status != http.StatusOK {
			logger = logger.With("status", status)
		}
		if latency > time.Second*3 {
			logger = logger.With("latency", latencyValue)
		}
		logger.Infof("[ARTICLE_API] %v | %3d | %s | %13v | %15s | %-7s %#v",
			timestamp.Format("2006/01/02 - 15:04:05"),
			status,
			latency,
			latencyValue,
			clientIP,
			method,
			path,
		)
	}
}

package middleware

import (
	"context"
	"gin-rest-api-example/pkg/logging"
	"gin-rest-api-example/pkg/trace"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
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

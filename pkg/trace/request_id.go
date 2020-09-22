package trace

import (
	"context"
	"github.com/gin-gonic/gin"
)

type contextKey = string

const requestIdKey = contextKey("requestId")

// WithRequestID creates a new context with the given request id attached.
func WithRequestID(ctx context.Context, requestId string) context.Context {
	if gCtx, ok := ctx.(*gin.Context); ok {
		ctx = gCtx.Request.Context()
	}
	return context.WithValue(ctx, requestIdKey, requestId)
}

// RequestIDFromContext returns a request id from given context if exist,
// otherwise returns empty string "".
func RequestIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if gCtx, ok := ctx.(*gin.Context); ok {
		ctx = gCtx.Request.Context()
	}
	if requestId, ok := ctx.Value(requestIdKey).(string); ok {
		return requestId
	}
	return ""
}

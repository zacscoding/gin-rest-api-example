package database

import (
	"context"
	"gorm.io/gorm"
)

type contextKey = string

const dbKey = contextKey("db")

// WithDB creates a new context with the provided db attached
func WithDB(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, dbKey, db)
}

// FromContext returns db stored in the context if exist, otherwise returns given db
func FromContext(ctx context.Context, db *gorm.DB) *gorm.DB {
	if ctx == nil {
		return db
	}
	if stored, ok := ctx.Value(dbKey).(*gorm.DB); ok {
		return stored
	}
	return db
}

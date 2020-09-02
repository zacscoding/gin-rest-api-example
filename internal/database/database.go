package database

import (
	"context"
	"fmt"
	"gin-rest-api-example/internal/config"
	"gorm.io/gorm"
)

type contextKey = string

const dbKey = contextKey("db")

// WithDB creates a new context with the provided db attached
func WithDB(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, db, db)
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

// NewDatabase creates a new database with given config
func NewDatabase(cfg *config.Config) (*gorm.DB, error) {
	fmt.Printf("database.NewDatabase() is called. config: %p\n", cfg)
	return nil, nil
}

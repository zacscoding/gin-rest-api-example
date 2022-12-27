package cache

import (
	"context"
	"errors"
	"fmt"
	"gin-rest-api-example/internal/config"
	"io"
)

var (
	ErrCacheMiss    = errors.New("key is missing")
	ErrInvalidKey   = errors.New("key is invalid")
	ErrInvalidValue = errors.New("value type is invalid")
)

type FetchFunc func() (interface{}, error)

//go:generate mockery --name Cacher --filename cache_mock.go
type Cacher interface {
	io.Closer
	// Fetch retrieves the item from the cache. If the item does not exist,
	// calls given FetchFunc to create a new item and sets to the cache.
	Fetch(ctx context.Context, key string, value interface{}, fetchFunc FetchFunc) error

	// Get gets an item for the given computeKey.
	Get(ctx context.Context, key string, value interface{}) error

	// Set adds an item to the cache.
	Set(ctx context.Context, key string, value interface{}) error

	// Exists returns a true if the given computeKey is exists, otherwise false.
	Exists(ctx context.Context, key string) (bool, error)

	// Delete removes an item from the cache.
	Delete(ctx context.Context, key string) error
}

func NewCacher(conf *config.Config) (Cacher, error) {
	if !conf.CacheConfig.Enabled {
		return nil, nil
	}
	switch conf.CacheConfig.Type {
	case "redis":
		return newRedisCacher(conf)
	default:
		return nil, fmt.Errorf("unknown cache type: %s", conf.CacheConfig.Type)
	}
}

type contextKey = string

const skipCacheKey = contextKey("skipCacheKey")

func IsCacheSkip(ctx context.Context) bool {
	if ctx == nil {
		return false
	}
	if skip, exists := ctx.Value(skipCacheKey).(bool); exists {
		return skip
	}
	return false
}

func WithCacheSkip(ctx context.Context, skipCache bool) context.Context {
	return context.WithValue(ctx, skipCacheKey, skipCache)
}

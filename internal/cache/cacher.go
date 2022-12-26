package cache

import (
	"context"
	"gin-rest-api-example/internal/config"
	"io"
)

type FetchFunc func() (interface{}, error)

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
	return nil, nil
}

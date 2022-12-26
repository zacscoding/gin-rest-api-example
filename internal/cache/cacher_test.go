package cache

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func testFetch(t *testing.T, cacher Cacher) {
	t.Run("Fetch from cache", func(t *testing.T) {
		key := uuid.New().String()
		value := "value1"
		assert.NoError(t, cacher.Set(context.TODO(), key, value))

		var find string
		cacher.Fetch(context.TODO(), key, &find, func() (interface{}, error) {
			t.Fail()
			return nil, nil
		})

		assert.EqualValues(t, value, find)
	})

	t.Run("Fetch from FetchFn", func(t *testing.T) {
		key := uuid.NewString()
		value := "value1"

		var (
			find   string
			called = false
		)
		err := cacher.Fetch(context.TODO(), key, &find, func() (interface{}, error) {
			called = true
			return value, nil
		})

		assert.NoError(t, err)
		assert.EqualValues(t, value, find)
		assert.True(t, called)

		exists, err := cacher.Exists(context.TODO(), key)
		assert.NoError(t, err)
		assert.True(t, exists)

		var find2 string
		assert.NoError(t, cacher.Get(context.TODO(), key, find2))
		assert.EqualValues(t, value, find2)
	})
}

package cache

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func testFetch(t *testing.T, cacher Cacher) {
	t.Run("Exist Item", func(t *testing.T) {
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

	t.Run("NotExist Item", func(t *testing.T) {
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
		assert.NoError(t, cacher.Get(context.TODO(), key, &find2))
		assert.EqualValues(t, value, find2)
	})

	t.Run("Invalid Key", func(t *testing.T) {
		err := cacher.Fetch(context.TODO(), "", "", nil)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), ErrInvalidKey.Error())
	})
}

func testGet(t *testing.T, cacher Cacher) {
	existKey := uuid.NewString()
	existValue := "value1"
	assert.NoError(t, cacher.Set(context.TODO(), existKey, existValue))

	t.Run("Exist Item", func(t *testing.T) {
		var find string
		err := cacher.Get(context.TODO(), existKey, &find)

		assert.NoError(t, err)
		assert.EqualValues(t, existValue, find)
	})

	t.Run("NotExist Item", func(t *testing.T) {
		key := uuid.NewString()

		var find string
		err := cacher.Get(context.TODO(), key, &find)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), ErrCacheMiss.Error())
	})

	t.Run("Invalid Key", func(t *testing.T) {
		var find string
		err := cacher.Get(context.TODO(), "", &find)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), ErrInvalidKey.Error())
	})

	t.Run("Invalid Value", func(t *testing.T) {
		key := uuid.NewString()
		value := "value1"
		assert.NoError(t, cacher.Set(context.TODO(), key, value))

		var find string
		err := cacher.Get(context.TODO(), key, find)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), ErrInvalidValue.Error())
	})
}

func testExists(t *testing.T, cacher Cacher) {
	existKey := uuid.NewString()
	existValue := "value1"
	assert.NoError(t, cacher.Set(context.TODO(), existKey, existValue))

	t.Run("Exist Item", func(t *testing.T) {
		ok, err := cacher.Exists(context.TODO(), existKey)

		assert.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("NotExist Item", func(t *testing.T) {
		ok, err := cacher.Exists(context.TODO(), uuid.NewString())

		assert.NoError(t, err)
		assert.False(t, ok)
	})

	t.Run("Invalid Key", func(t *testing.T) {
		ok, err := cacher.Exists(context.TODO(), uuid.NewString())

		assert.NoError(t, err)
		assert.False(t, ok)
	})
}

func testDelete(t *testing.T, cacher Cacher) {
	t.Run("Exist Item", func(t *testing.T) {
		key := uuid.NewString()
		value := "value"
		assert.NoError(t, cacher.Set(context.TODO(), key, value))

		err := cacher.Delete(context.TODO(), key)

		assert.NoError(t, err)
		ok, err := cacher.Exists(context.TODO(), key)
		assert.NoError(t, err)
		assert.False(t, ok)
	})

	t.Run("NotExist Item", func(t *testing.T) {
		err := cacher.Delete(context.TODO(), uuid.NewString())

		assert.NoError(t, err)
	})

	t.Run("With Invalid Key", func(t *testing.T) {
		err := cacher.Delete(context.TODO(), "")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), ErrInvalidKey.Error())
	})
}

func testSet(t *testing.T, cacher Cacher) {
	t.Run("Set valid", func(t *testing.T) {
		key := uuid.NewString()
		value := "value1"

		err := cacher.Set(context.TODO(), key, value)

		assert.NoError(t, err)
		ok, err := cacher.Exists(context.TODO(), key)
		assert.NoError(t, err)
		assert.True(t, ok)
	})
}

package database

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTemp(t *testing.T) {
	err := migrateDB("root:password@tcp(127.0.0.1:3306)/local_db?charset=utf8&parseTime=True")
	assert.NoError(t, err)
}

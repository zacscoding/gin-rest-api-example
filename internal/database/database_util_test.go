package database

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TODO : REMOVE (temporary for migration)
func TestConsole(t *testing.T) {
	t.Skip()
	dsn := "root:password@tcp(127.0.0.1:3306)/local_db?charset=utf8&parseTime=True&multiStatements=true"
	err := migrateDB(dsn)
	assert.NoError(t, err)
}

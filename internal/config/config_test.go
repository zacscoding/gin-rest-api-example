package config

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	cfg, err := Load("")

	assert.NoError(t, err)
	// server configs
	assert.Equal(t, defaultConfig["server.port"].(int), cfg.ServerConfig.Port)
	assert.Equal(t, defaultConfig["server.timeoutSecs"].(int), cfg.ServerConfig.TimeoutSecs)
	assert.Equal(t, defaultConfig["server.readTimeoutSecs"].(int), cfg.ServerConfig.ReadTimeoutSecs)
	assert.Equal(t, defaultConfig["server.writeTimeoutSecs"].(int), cfg.ServerConfig.WriteTimeoutSecs)
	// db configs
	assert.Equal(t, defaultConfig["db.driverName"].(string), cfg.DBConfig.DriverName)
	assert.Equal(t, defaultConfig["db.dataSourceName"].(string), cfg.DBConfig.DataSourceName)
	assert.Equal(t, defaultConfig["db.pool.maxOpen"].(int), cfg.DBConfig.Pool.MaxOpen)
	assert.Equal(t, defaultConfig["db.pool.maxIdle"].(int), cfg.DBConfig.Pool.MaxIdle)
	assert.Equal(t, defaultConfig["db.pool.maxLifetime"].(int), cfg.DBConfig.Pool.MaxLifetime)
	assert.Equal(t, defaultConfig["db.createTable"].(bool), cfg.DBConfig.CreateTable)
}

func TestLoadWithEnv(t *testing.T) {
	// given
	err := os.Setenv("ARTICLE_SERVER_SERVER_PORT", "4000")
	assert.NoError(t, err)

	// when
	cfg, err := Load("")

	// then
	assert.NoError(t, err)
	assert.Equal(t, 4000, cfg.ServerConfig.Port)
}

func TestLoadWithConfigFile(t *testing.T) {
	// given
	err := os.Setenv("ARTICLE_SERVER_SERVER_PORT", "4000")
	assert.NoError(t, err)

	config := `
server:
  port: 5000
`
	tempFile, err := ioutil.TempFile(os.TempDir(), "article-server-test")
	assert.NoError(t, err)
	fmt.Println("Create temp file::", tempFile.Name())
	defer os.Remove(tempFile.Name())

	_, err = tempFile.WriteString(config)
	assert.NoError(t, err)

	// when
	cfg, err := Load(tempFile.Name())

	// then
	assert.NoError(t, err)
	assert.Equal(t, 5000, cfg.ServerConfig.Port)
}

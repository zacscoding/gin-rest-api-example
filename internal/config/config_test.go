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
	// jwt configs
	assert.Equal(t, defaultConfig["jwt.secret"].(string), cfg.JwtConfig.Secret)
	assert.Equal(t, defaultConfig["jwt.sessionTime"].(int), cfg.JwtConfig.SessionTime)
	// db configs
	assert.Equal(t, defaultConfig["db.dataSourceName"].(string), cfg.DBConfig.DataSourceName)
	assert.Equal(t, defaultConfig["db.migrate.enable"].(bool), cfg.DBConfig.Migrate.Enable)
	assert.Equal(t, defaultConfig["db.migrate.dir"].(bool), cfg.DBConfig.Migrate.Dir)
	assert.Equal(t, defaultConfig["db.pool.maxOpen"].(int), cfg.DBConfig.Pool.MaxOpen)
	assert.Equal(t, defaultConfig["db.pool.maxIdle"].(int), cfg.DBConfig.Pool.MaxIdle)
	assert.Equal(t, defaultConfig["db.pool.maxLifetime"].(int), cfg.DBConfig.Pool.MaxLifetime)
	// metrics configs
	assert.Equal(t, defaultConfig["metrics.namespace"].(string), cfg.MetricsConfig.Namespace)
	assert.Equal(t, defaultConfig["metrics.subsystem"].(string), cfg.MetricsConfig.Subsystem)
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

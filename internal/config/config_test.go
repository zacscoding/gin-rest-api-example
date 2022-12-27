package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	cfg, err := Load("")
	assert.NoError(t, err)

	// server configs
	equal(t, 8080, defaultConfig["server.port"], cfg.ServerConfig.Port)
	equalDuration(t, 5*time.Second, defaultConfig["server.readTimeout"], cfg.ServerConfig.ReadTimeout)
	equalDuration(t, 10*time.Second, defaultConfig["server.writeTimeout"], cfg.ServerConfig.WriteTimeout)
	equalDuration(t, 30*time.Second, defaultConfig["server.gracefulShutdown"], cfg.ServerConfig.GracefulShutdown)
	// logging configs
	equal(t, -1, defaultConfig["logging.level"], cfg.LoggingConfig.Level)
	equal(t, "console", defaultConfig["logging.encoding"], cfg.LoggingConfig.Encoding)
	equal(t, true, defaultConfig["logging.development"], cfg.LoggingConfig.Development)
	// jwt configs
	equal(t, "secret-key", defaultConfig["jwt.secret"], cfg.JwtConfig.Secret)
	equalDuration(t, 864000*time.Second, defaultConfig["jwt.sessionTime"], cfg.JwtConfig.SessionTime)
	// db configs
	equal(t, "root:password@tcp(127.0.0.1:3306)/local_db?charset=utf8&parseTime=True&multiStatements=true", defaultConfig["db.dataSourceName"], cfg.DBConfig.DataSourceName)
	equal(t, 1, defaultConfig["db.logLevel"], cfg.DBConfig.LogLevel)
	equal(t, false, defaultConfig["db.migrate.enable"], cfg.DBConfig.Migrate.Enable)
	equal(t, "", defaultConfig["db.migrate.dir"], cfg.DBConfig.Migrate.Dir)
	equal(t, 10, defaultConfig["db.pool.maxOpen"], cfg.DBConfig.Pool.MaxOpen)
	equal(t, 5, defaultConfig["db.pool.maxIdle"], cfg.DBConfig.Pool.MaxIdle)
	equalDuration(t, 5*time.Minute, defaultConfig["db.pool.maxLifetime"], cfg.DBConfig.Pool.MaxLifetime)
	// cache configs
	equal(t, false, defaultConfig["cache.enabled"].(bool), cfg.CacheConfig.Enabled)
	equal(t, "article-", defaultConfig["cache.prefix"].(string), cfg.CacheConfig.Prefix)
	equal(t, "redis", defaultConfig["cache.type"].(string), cfg.CacheConfig.Type)
	equalDuration(t, 60*time.Second, defaultConfig["cache.ttl"], cfg.CacheConfig.TTL)
	equal(t, false, defaultConfig["cache.redis.cluster"].(bool), cfg.CacheConfig.RedisConfig.Cluster)
	equal(t, []string{"localhost:6379"}, defaultConfig["cache.redis.endpoints"].([]string), cfg.CacheConfig.RedisConfig.Endpoints)
	equalDuration(t, 3*time.Second, defaultConfig["cache.redis.readTimeout"], cfg.CacheConfig.RedisConfig.ReadTimeout)
	equalDuration(t, 3*time.Second, defaultConfig["cache.redis.writeTimeout"], cfg.CacheConfig.RedisConfig.WriteTimeout)
	equalDuration(t, 5*time.Second, defaultConfig["cache.redis.dialTimeout"], cfg.CacheConfig.RedisConfig.DialTimeout)
	equal(t, 10, defaultConfig["cache.redis.poolSize"].(int), cfg.CacheConfig.RedisConfig.PoolSize)
	equalDuration(t, 1*time.Minute, defaultConfig["cache.redis.poolTimeout"], cfg.CacheConfig.RedisConfig.PoolTimeout)
	equalDuration(t, 0, defaultConfig["cache.redis.maxConnAge"], cfg.CacheConfig.RedisConfig.MaxConnAge)
	equalDuration(t, 5*time.Minute, defaultConfig["cache.redis.idleTimeout"], cfg.CacheConfig.RedisConfig.IdleTimeout)
	// metrics configs
	equal(t, "article_server", defaultConfig["metrics.namespace"], cfg.MetricsConfig.Namespace)
	equal(t, "", defaultConfig["metrics.subsystem"], cfg.MetricsConfig.Subsystem)
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

func TestMarshalJSON(t *testing.T) {
	conf, err := Load("")
	assert.NoError(t, err)
	data, err := json.Marshal(conf)
	assert.NoError(t, err)

	var configMap map[string]interface{}
	assert.NoError(t, json.Unmarshal(data, &configMap))
	assert.True(t, strings.HasPrefix(configMap["db.dataSourceName"].(string), "root:****@tcp"))
	assert.Equal(t, "****", configMap["jwt.secret"])
}

func equal(t *testing.T, expected interface{}, values ...interface{}) {
	for _, v := range values {
		assert.EqualValues(t, expected, v)
	}
}

func equalDuration(t *testing.T, expected time.Duration, values ...interface{}) {
	for _, v := range values {
		if str, ok := v.(string); ok {
			d, err := time.ParseDuration(str)
			assert.NoError(t, err)
			assert.EqualValues(t, expected, d)
			continue
		}
		assert.EqualValues(t, expected, v)
	}
}

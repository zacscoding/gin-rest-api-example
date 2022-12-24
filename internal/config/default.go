package config

var defaultConfig = map[string]interface{}{
	"server.port":             8080,
	"server.readTimeout":      "5s",
	"server.writeTimeout":     "10s",
	"server.gracefulShutdown": "30s",

	"logging.level":       -1,
	"logging.encoding":    "console",
	"logging.development": true,

	"jwt.secret":      "secret-key",
	"jwt.sessionTime": "864000s",

	"db.dataSourceName":   "root:password@tcp(127.0.0.1:3306)/local_db?charset=utf8&parseTime=True&multiStatements=true",
	"db.logLevel":         1,
	"db.migrate.enable":   false,
	"db.migrate.dir":      "",
	"db.pool.maxOpen":     10,
	"db.pool.maxIdle":     5,
	"db.pool.maxLifetime": "5m",

	"metrics.namespace": "article_server",
	"metrics.subsystem": "",
}

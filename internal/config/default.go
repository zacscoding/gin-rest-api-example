package config

var defaultConfig = map[string]interface{}{
	"server.port":             3000,
	"server.timeoutSecs":      5,
	"server.readTimeoutSecs":  5,
	"server.writeTimeoutSecs": 30,

	"jwt.secret":      "secret-key",
	"jwt.sessionTime": 864000,

	"db.dataSourceName":   "root:password@tcp(127.0.0.1:3306)/local_db?charset=utf8&parseTime=True&multiStatements=true",
	"db.migrate.enable":   false,
	"db.migrate.dir":      "",
	"db.pool.maxOpen":     50,
	"db.pool.maxIdle":     5,
	"db.pool.maxLifetime": 5,
}

package database

import (
	"gin-rest-api-example/internal/config"
	"gin-rest-api-example/pkg/logging"
	"time"

	"go.uber.org/zap/zapcore"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// NewDatabase creates a new database with given config
func NewDatabase(cfg *config.Config) (*gorm.DB, error) {
	var (
		db     *gorm.DB
		err    error
		logger = NewLogger(time.Second, true, zapcore.Level(cfg.DBConfig.LogLevel))
	)

	for i := 0; i <= 30; i++ {
		db, err = gorm.Open(mysql.Open(cfg.DBConfig.DataSourceName), &gorm.Config{Logger: logger})
		if err == nil {
			break
		}
		logging.DefaultLogger().Warnf("failed to open database: %v", err)
		time.Sleep(500 * time.Millisecond)
	}
	if err != nil {
		return nil, err
	}

	rawDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	rawDB.SetMaxOpenConns(cfg.DBConfig.Pool.MaxOpen)
	rawDB.SetMaxIdleConns(cfg.DBConfig.Pool.MaxIdle)
	rawDB.SetConnMaxLifetime(cfg.DBConfig.Pool.MaxLifetime)

	if cfg.DBConfig.Migrate.Enable {
		err := migrateDB(cfg.DBConfig.DataSourceName, cfg.DBConfig.Migrate.Dir)
		if err != nil {
			return nil, err
		}
	}
	return db, nil
}

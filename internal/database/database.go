package database

import (
	"gin-rest-api-example/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

// NewDatabase creates a new database with given config
func NewDatabase(cfg *config.Config) (*gorm.DB, error) {
	var (
		db  *gorm.DB
		err error
	)

	for i := 0; i <= 30; i++ {
		db, err = gorm.Open(mysql.Open(cfg.DBConfig.DataSourceName), &gorm.Config{})
		if err != nil {
			time.Sleep(500 * time.Millisecond)
		}
	}
	if err != nil {
		return nil, err
	}

	origin, err := db.DB()
	if err != nil {
		return nil, err
	}
	origin.SetMaxOpenConns(cfg.DBConfig.Pool.MaxOpen)
	origin.SetMaxIdleConns(cfg.DBConfig.Pool.MaxIdle)
	origin.SetConnMaxLifetime(time.Duration(cfg.DBConfig.Pool.MaxLifetime) * time.Second)

	if cfg.DBConfig.Migrate {
		err := migrateDB(cfg.DBConfig.DataSourceName)
		if err != nil {
			return nil, err
		}
	}
	return db, nil
}

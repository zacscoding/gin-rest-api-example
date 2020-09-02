package database

import (
	"context"
	"fmt"
	"gin-rest-api-example/internal/account/model"
	"gin-rest-api-example/internal/database"
	"gin-rest-api-example/pkg/logging"
	"gorm.io/gorm"
)

type AccountDB interface {
	// Save saves a given account
	Save(ctx context.Context, account *model.Account) error
}

type accountDB struct {
	db *gorm.DB
}

func (a *accountDB) Save(ctx context.Context, account *model.Account) error {
	logger := logging.FromContext(ctx)
	db := database.FromContext(ctx, a.db)
	logger.Debugw("account.db.Save", "account", account)

	if err := db.Create(account).Error; err != nil {
		logger.Error("account.db.Save failed to save", "err", err)
		return err
	}
	return nil
}

// NewAccountDB create a new account db with given db
func NewAccountDB(db *gorm.DB) AccountDB {
	fmt.Printf("account.database.NewAccountDB() is called. db: %p\n", db)
	return &accountDB{
		db: db,
	}
}

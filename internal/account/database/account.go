package database

import (
	"context"
	"gin-rest-api-example/internal/account/model"
	"gin-rest-api-example/internal/cache"
	"gin-rest-api-example/internal/database"
	"gin-rest-api-example/internal/metric"
	"gin-rest-api-example/pkg/logging"

	"gorm.io/gorm"
)

//go:generate mockery --name AccountDB --filename account_mock.go
type AccountDB interface {
	// Save saves a given account
	Save(ctx context.Context, account *model.Account) error

	// Update updates a given account
	Update(ctx context.Context, email string, account *model.Account) error

	// FindByEmail returns an account with given email if exist
	FindByEmail(ctx context.Context, email string) (*model.Account, error)
}

// NewAccountDB creates a new account db with given db
func NewAccountDB(db *gorm.DB, cacher cache.Cacher, mp *metric.MetricsProvider) AccountDB {
	if cacher == nil {
		return &accountDB{db: db}
	}
	return newAccountCacheDB(cacher, mp, &accountDB{db: db})
}

type accountDB struct {
	db *gorm.DB
}

func (a *accountDB) Save(ctx context.Context, account *model.Account) error {
	logger := logging.FromContext(ctx)
	db := database.FromContext(ctx, a.db)
	logger.Debugw("account.db.Save", "account", account)

	if err := db.WithContext(ctx).Create(account).Error; err != nil {
		logger.Error("account.db.Save failed to save", "err", err)
		if database.IsKeyConflictErr(err) {
			return database.ErrKeyConflict
		}
		return err
	}
	return nil
}

func (a *accountDB) Update(ctx context.Context, email string, account *model.Account) error {
	logger := logging.FromContext(ctx)
	db := database.FromContext(ctx, a.db)
	logger.Debugw("account.db.Update", "account", account)

	fields := make(map[string]interface{})
	if account.Username != "" {
		fields["username"] = account.Username
	}
	if account.Password != "" {
		fields["password"] = account.Password
	}
	if account.Bio != "" {
		fields["bio"] = account.Bio
	}
	if account.Image != "" {
		fields["image"] = account.Image
	}

	chain := db.WithContext(ctx).
		Model(&model.Account{}).
		Where("email = ?", email).
		UpdateColumns(fields)
	if chain.Error != nil {
		logger.Error("account.db.Update failed to update", "err", chain.Error)
		return chain.Error
	}
	if chain.RowsAffected == 0 {
		return database.ErrNotFound
	}
	return nil
}

func (a *accountDB) FindByEmail(ctx context.Context, email string) (*model.Account, error) {
	logger := logging.FromContext(ctx)
	db := database.FromContext(ctx, a.db)
	logger.Debugw("account.db.FindByEmail", "email", email)

	var acc model.Account
	if err := db.WithContext(ctx).Where("email = ?", email).First(&acc).Error; err != nil {
		logger.Error("account.db.FindByEmail failed to find", "err", err)
		if database.IsRecordNotFoundErr(err) {
			return nil, database.ErrNotFound
		}
		return nil, err
	}
	return &acc, nil
}

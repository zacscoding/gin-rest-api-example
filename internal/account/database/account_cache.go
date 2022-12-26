package database

import (
	"context"
	"gin-rest-api-example/internal/account/model"
	"gin-rest-api-example/internal/metric"
)

var _ AccountDB = (*accountCachedDB)(nil)

type accountCachedDB struct {
	// cache *cache.Cacher
	mp       *metric.MetricsProvider
	delegate AccountDB
}

func (ac *accountCachedDB) Save(ctx context.Context, account *model.Account) error {
	if err := ac.delegate.Save(ctx, account); err != nil {
		return err
	}

	//TODO implement me
	panic("implement me")
}

func (ac *accountCachedDB) Update(ctx context.Context, email string, account *model.Account) error {
	//TODO implement me
	panic("implement me")
}

func (ac *accountCachedDB) FindByEmail(ctx context.Context, email string) (*model.Account, error) {
	//TODO implement me
	panic("implement me")
}

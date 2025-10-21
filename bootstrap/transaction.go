package bootstrap

import (
	"context"
	"gorm.io/gorm"
)

type TransactionKey struct{}
type DBClient struct {
	Db *gorm.DB
}

func (d *DBClient) Session(ctx context.Context) *gorm.DB {
	if client, ok := ctx.Value(TransactionKey{}).(*DBClient); ok {
		return client.Db
	}
	return d.Db.WithContext(ctx)
}

func (d *DBClient) Transaction(ctx context.Context, fn func(context.Context) error) error {
	return d.Db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(context.WithValue(ctx, TransactionKey{}, &DBClient{Db: tx}))
	})
}

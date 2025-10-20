package bootstrap

import (
	"context"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type TransactionKey struct{}
type DBClient struct {
	db *gorm.DB
}

func (d *DBClient) Session(ctx context.Context) *gorm.DB {
	client, ok := ctx.Value(TransactionKey{}).(*DBClient)
	db := lo.Ternary(ok, client.db, d.db.WithContext(ctx))
	return db
}

func (d *DBClient) Transaction(ctx context.Context, fn func(context.Context) error) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(context.WithValue(ctx, TransactionKey{}, &DBClient{db: tx}))
	})
}

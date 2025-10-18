package bootstrap

import (
	"context"
	"gorm.io/gorm"
)

type DBTransactionKey struct{}
type DBClient struct {
	db *gorm.DB
}

func (d *DBClient) DB(ctx context.Context) *gorm.DB {
	db, ok := ctx.Value(DBTransactionKey{}).(*gorm.DB)
	if ok {
		return db
	}
	// 新建一个db
	return d.db.WithContext(ctx)
}

func (d *DBClient) Transaction(ctx context.Context, fn func(context.Context) error) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(context.WithValue(ctx, DBTransactionKey{}, tx))
	})
}

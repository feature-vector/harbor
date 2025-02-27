package db

import (
	"context"
	"gorm.io/gorm"
)

type transactionContextKey struct{}

var (
	globalDb               *gorm.DB
	_transactionContextKey = transactionContextKey{}
)

func SetGlobalDb(db *gorm.DB) {
	globalDb = db
}

func BeginTx(ctx context.Context) context.Context {
	return context.WithValue(ctx, _transactionContextKey, globalDb.WithContext(ctx).Begin())
}

func GetTx(ctx context.Context) *gorm.DB {
	return ctx.Value(_transactionContextKey).(*gorm.DB)
}

func CommitOrRollbackTx(ctx context.Context) {
	tx := GetTx(ctx)
	if tx.Error != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}
}

func CommitTx(ctx context.Context) error {
	tx := GetTx(ctx)
	tx.Commit()
	return tx.Error
}

func RollbackTx(ctx context.Context) {
	tx := GetTx(ctx)
	tx.Rollback()
}

func FromContext(ctx context.Context) *gorm.DB {
	tx := ctx.Value(_transactionContextKey)
	if tx != nil {
		return tx.(*gorm.DB)
	}
	return globalDb.WithContext(ctx)
}

func Global() *gorm.DB {
	return globalDb
}

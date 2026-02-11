package gorm

import (
	"context"

	"github.com/mvcris/maya-guessr/backend/internal/core/transactions"
	"gorm.io/gorm"
)

type txContextKey struct{}

func injectTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txContextKey{}, tx)
}

func ExtractTx(ctx context.Context) (*gorm.DB, bool) {
	tx, ok := ctx.Value(txContextKey{}).(*gorm.DB)
	return tx, ok
}

type GormTransactionManager struct {
	db *gorm.DB
}

func NewGormTransactionManager(db *gorm.DB) transactions.TransactionManager {
	return &GormTransactionManager{db: db}
}

func (m *GormTransactionManager) RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	tx := m.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	txCtx := injectTx(ctx, tx)
	if err := fn(txCtx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

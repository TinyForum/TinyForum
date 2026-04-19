package transaction

import (
	"context"
	"gorm.io/gorm"
)

type TransactionManager interface {
	ExecuteInTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error
	BeginTx(ctx context.Context) *gorm.DB
}

type transactionManager struct {
	db *gorm.DB
}

func NewTransactionManager(db *gorm.DB) TransactionManager {
	return &transactionManager{db: db}
}

func (tm *transactionManager) ExecuteInTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return tm.db.WithContext(ctx).Transaction(fn)
}

func (tm *transactionManager) BeginTx(ctx context.Context) *gorm.DB {
	return tm.db.WithContext(ctx).Begin()
}

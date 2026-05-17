package transaction

import (
	"context"

	"gorm.io/gorm"
)

func (tm *transactionManager) ExecuteInTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return tm.db.WithContext(ctx).Transaction(fn)
}

func (tm *transactionManager) BeginTx(ctx context.Context) *gorm.DB {
	return tm.db.WithContext(ctx).Begin()
}

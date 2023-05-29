package trans

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/happyxhw/pkg/cx"
)

type Trans struct {
	db *gorm.DB
}

func NewTrans(db *gorm.DB) *Trans {
	return &Trans{
		db: db,
	}
}

func (t *Trans) Exec(ctx context.Context, fn func(ctx context.Context) error) error {
	if _, ok := cx.FromTx(ctx); ok {
		return fn(ctx)
	}

	return t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(cx.NewTx(ctx, tx))
	})
}

func DB(ctx context.Context, db *gorm.DB) *gorm.DB {
	if v, ok := cx.FromTx(ctx); ok && !cx.FromNoTx(ctx) {
		tx, ok2 := v.(*gorm.DB)
		if ok2 {
			if cx.FromTxLock(ctx) {
				tx = tx.Clauses(clause.Locking{Strength: "UPDATE"})
			}
			return tx
		}
	}

	return db.WithContext(ctx)
}

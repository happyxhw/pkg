package cx

import (
	"context"

	"gorm.io/gorm"
)

type (
	txCtx     struct{}
	noTxCtx   struct{}
	txLockCtx struct{}
	metricCtx struct{}
	traceCtx  struct{}
)

// NewTx wrap tx in context
func NewTx(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, txCtx{}, db)
}

func FromTx(ctx context.Context) (any, bool) {
	v := ctx.Value(txCtx{})
	return v, v != nil
}

// NewNoTx wrap no tx in context
func NewNoTx(ctx context.Context) context.Context {
	return context.WithValue(ctx, noTxCtx{}, true)
}

func FromNoTx(ctx context.Context) bool {
	v := ctx.Value(noTxCtx{})
	return v != nil && v.(bool)
}

// NewTxLock wrap lock tx in context
func NewTxLock(ctx context.Context) context.Context {
	return context.WithValue(ctx, txLockCtx{}, true)
}

func FromTxLock(ctx context.Context) bool {
	v := ctx.Value(txLockCtx{})
	return v != nil && v.(bool)
}

func NewMetricCtx(ctx context.Context) context.Context {
	return context.WithValue(ctx, metricCtx{}, true)
}

func FromMetricCtx(ctx context.Context) bool {
	v := ctx.Value(metricCtx{})
	return v != nil && v.(bool)
}

func NewTraceCtx(ctx context.Context, requestID interface{}) context.Context {
	return context.WithValue(ctx, traceCtx{}, requestID)
}

func RequestID(ctx context.Context) string {
	id, _ := ctx.Value(traceCtx{}).(string)
	return id
}

// Package database provides transaction-safe database helpers and context
// propagation for dotNext repositories and services.
package database

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type txCtxKey struct{}

// WithTx executes fn inside a MariaDB/GORM transaction. If ctx already carries an
// active transaction (from an outer WithTx call), fn participates in that outer
// transaction rather than starting a nested connection, ensuring atomic
// multi-repository domain operations without connection-pool deadlocks.
func WithTx(ctx context.Context, db *gorm.DB, fn func(tx *gorm.DB) error) error {
	if fn == nil {
		return nil
	}
	if existing, ok := TxFromContext(ctx); ok {
		return fn(existing)
	}
	if db == nil {
		return errors.New("database: nil *gorm.DB passed to WithTx")
	}
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}

// ContextWithTx binds tx onto ctx so repositories using DBFromContext(ctx, fallback)
// automatically enlist in the caller's active transaction.
func ContextWithTx(ctx context.Context, tx *gorm.DB) context.Context {
	if tx == nil {
		return ctx
	}
	return context.WithValue(ctx, txCtxKey{}, tx)
}

// TxFromContext extracts an active transaction from ctx, if present.
func TxFromContext(ctx context.Context) (*gorm.DB, bool) {
	if ctx == nil {
		return nil, false
	}
	tx, ok := ctx.Value(txCtxKey{}).(*gorm.DB)
	return tx, ok && tx != nil
}

// FromContext returns the active transaction on ctx if one exists, or fallback.WithContext(ctx).
func FromContext(ctx context.Context, fallback *gorm.DB) *gorm.DB {
	if tx, ok := TxFromContext(ctx); ok {
		return tx
	}
	if fallback == nil {
		return nil
	}
	return fallback.WithContext(ctx)
}

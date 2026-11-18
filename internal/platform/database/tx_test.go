package database

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func openMockGorm(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })

	gdb, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("gorm.Open: %v", err)
	}
	return gdb, mock
}

func TestWithTxCommitAndRollback(t *testing.T) {
	t.Parallel()

	gdb, mock := openMockGorm(t)

	mock.ExpectBegin()
	mock.ExpectCommit()
	if err := WithTx(context.Background(), gdb, func(tx *gorm.DB) error {
		ctx := ContextWithTx(context.Background(), tx)
		if FromContext(ctx, gdb) != tx {
			t.Fatal("expected FromContext to return active tx")
		}
		return nil
	}); err != nil {
		t.Fatalf("WithTx commit failed: %v", err)
	}

	mock.ExpectBegin()
	mock.ExpectRollback()
	boom := errors.New("boom")
	if err := WithTx(context.Background(), gdb, func(tx *gorm.DB) error {
		return boom
	}); !errors.Is(err, boom) {
		t.Fatalf("expected boom rollback error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestWithTxReusesOuterTxAndRejectsNilDB(t *testing.T) {
	t.Parallel()

	if err := WithTx(context.Background(), nil, func(*gorm.DB) error { return nil }); err == nil {
		t.Fatal("expected error on nil *gorm.DB")
	}

	gdb, mock := openMockGorm(t)
	mock.ExpectBegin()
	mock.ExpectCommit()

	err := WithTx(context.Background(), gdb, func(outer *gorm.DB) error {
		ctx := ContextWithTx(context.Background(), outer)
		return WithTx(ctx, gdb, func(inner *gorm.DB) error {
			if inner != outer {
				t.Fatal("expected nested WithTx to reuse outer transaction")
			}
			return nil
		})
	})
	if err != nil {
		t.Fatalf("nested WithTx failed: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

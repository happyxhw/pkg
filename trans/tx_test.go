package trans

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	"github.com/happyxhw/iself/model"
	"github.com/happyxhw/iself/pkg/cx"
	"github.com/happyxhw/iself/pkg/mymock"
)

func TestTx_ExecWithTx(t *testing.T) {
	gdb, gdbMock, _ := mymock.MockEqualDB()
	sql := `SELECT * FROM "user" WHERE id = $1 AND "user"."deleted_at" = $2`
	gdbMock.ExpectBegin()
	gdbMock.ExpectQuery(sql).
		WithArgs(1, 0).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "email"}).
				AddRow(1, "mock", "mock@mock.com"),
		)
	gdbMock.ExpectCommit()

	tx := Trans{db: gdb}

	err := tx.Exec(context.TODO(), func(ctx context.Context) error {
		db := DB(ctx, gdb)
		var u model.User
		return db.Where("id = ?", 1).Find(&u).Error
	})

	if err := gdbMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	require.NoError(t, err)
}

func TestTx_ExecWithoutTx(t *testing.T) {
	gdb, gdbMock, _ := mymock.MockEqualDB()
	sql := `SELECT * FROM "user" WHERE id = $1 AND "user"."deleted_at" = $2`
	gdbMock.ExpectQuery(sql).
		WithArgs(1, 0).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "email"}).
				AddRow(1, "mock", "mock@mock.com"),
		)

	tx := Trans{db: gdb}

	err := tx.Exec(cx.NewTx(context.TODO(), gdb), func(ctx context.Context) error {
		db := DB(ctx, gdb)
		var u model.User
		return db.Where("id = ?", 1).Find(&u).Error
	})

	if err := gdbMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	require.NoError(t, err)
}

func TestTx_ExecWithTxLock(t *testing.T) {
	gdb, gdbMock, _ := mymock.MockEqualDB()
	sql := `SELECT * FROM "user" WHERE id = $1 AND "user"."deleted_at" = $2 FOR UPDATE`
	gdbMock.ExpectBegin()
	gdbMock.ExpectQuery(sql).
		WithArgs(1, 0).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "email"}).
				AddRow(1, "mock", "mock@mock.com"),
		)
	gdbMock.ExpectCommit()

	tx := Trans{db: gdb}

	err := tx.Exec(cx.NewTxLock(context.TODO()), func(ctx context.Context) error {
		db := DB(ctx, gdb)
		var u model.User
		return db.Where("id = ?", 1).Find(&u).Error
	})

	if err := gdbMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	require.NoError(t, err)
}

func TestTx_ExecWithNoTx(t *testing.T) {
	gdb, gdbMock, _ := mymock.MockEqualDB()
	sql := `SELECT * FROM "user" WHERE id = $1 AND "user"."deleted_at" = $2`
	gdbMock.ExpectBegin()
	gdbMock.ExpectQuery(sql).
		WithArgs(1, 0).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "email"}).
				AddRow(1, "mock", "mock@mock.com"),
		)
	gdbMock.ExpectCommit()

	tx := Trans{db: gdb}

	err := tx.Exec(cx.NewNoTx(context.TODO()), func(ctx context.Context) error {
		db := DB(ctx, gdb)
		var u model.User
		return db.Where("id = ?", 1).Find(&u).Error
	})

	if err := gdbMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	require.NoError(t, err)
}

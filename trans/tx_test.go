package trans

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"gorm.io/plugin/soft_delete"

	"github.com/happyxhw/pkg/cx"
	"github.com/happyxhw/pkg/mymock"
)

type User struct {
	ID        int64  `gorm:"column:id;" json:"id"`
	Name      string `gorm:"column:name;" json:"name"`
	Email     string `gorm:"column:email;" json:"email"`
	Password  string `gorm:"column:password;" json:"password"`
	AvatarURL string `gorm:"avatar_url" json:"avatar_url"`
	Role      int    `gorm:"role" json:"role"`
	Source    string `gorm:"source" json:"source"`
	SourceID  int64  `gorm:"source_id" json:"source_id"`
	Status    int    `gorm:"status" json:"status"`

	CreatedAt time.Time             `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time             `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt soft_delete.DeletedAt `gorm:"column:deleted_at;default:0" json:"deleted_at"`
}

func (u User) TableName() string {
	return "user"
}

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
		var u User
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
		var u User
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
		var u User
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
		var u User
		return db.Where("id = ?", 1).Find(&u).Error
	})

	if err := gdbMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	require.NoError(t, err)
}

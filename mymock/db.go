package mymock

import (
	"database/sql"
	"os"
	"time"

	stdlog "log" //nolint:depguard

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func MockEqualDB() (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		return nil, nil, err
	}
	gdb := gormDB(db)

	return gdb, mock, err
}

func MockRegexDB() (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	gdb := gormDB(db)
	return gdb, mock, err
}

func gormDB(db *sql.DB) *gorm.DB {
	newLogger := logger.New(
		stdlog.New(os.Stdout, "\r\n", stdlog.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)
	gdb, _ := gorm.Open(postgres.New(postgres.Config{
		WithoutReturning:     true,
		Conn:                 db,
		DriverName:           "postgres",
		PreferSimpleProtocol: true,
	}), &gorm.Config{Logger: newLogger})

	gdb.SkipDefaultTransaction = true
	return gdb
}

package query

import (
	"errors"

	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

func IsErrDuplicatedKey(err error) bool {
	if err == nil {
		return false
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return true
	}

	return errors.Is(err, gorm.ErrDuplicatedKey)
}

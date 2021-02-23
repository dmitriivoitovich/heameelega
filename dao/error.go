package dao

import (
	"errors"

	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

func IsErrDuplicateKey(err error) bool {
	var pgError *pgconn.PgError

	return errors.As(err, &pgError) && pgError.Code == "23505"
}

func IsErrRecordNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

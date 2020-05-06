package utils

import (
	"context"
	"database/sql"

	"github.com/volatiletech/sqlboiler/v4/boil"
)

func Tx(ctx context.Context, fn func(tx *sql.Tx) error) (err error) {
	tx, err := boil.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if recv := recover(); recv != nil {
			tx.Rollback()
			panic(recv)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = fn(tx)
	return err
}

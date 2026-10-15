package postgres

import (
	"database/sql"
	"errors"
	"fmt"
)

var ErrEditData error = errors.New("failed to change the data")

func IFChangesInDB(result sql.Result) error {
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get affected rows: %w", err)
	}

	if rows == 0 {
		return ErrEditData
	}

	return nil
}

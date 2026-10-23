package postgres

import (
	"database/sql"
	"fmt"

	"github.com/lorsanstand/gontrol/internal/hub/store"
)

func ifChangesInDB(result sql.Result) error {
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get affected rows: %w", err)
	}

	if rows == 0 {
		return store.ErrChangeData
	}

	return nil
}

package postgres

import (
	"database/sql"
	"fmt"
)

func deleteDocsByPrimaryKeyPrefix(db *sql.DB, primaryKeyPrefix, tableName string) error {
	if _, err := db.Exec(
		// tableName comes from this code, this isn't unsafe
		fmt.Sprintf("DELETE FROM %s WHERE id LIKE $1", tableName),
		fmt.Sprintf("%s%%", primaryKeyPrefix),
	); err != nil {
		return fmt.Errorf("could not perform Postgres DELETE statement to delete docs in %s table with ID prefix %s: %v", tableName, primaryKeyPrefix, err)
	}

	return nil
}

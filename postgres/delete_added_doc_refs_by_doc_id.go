package postgres

import (
	"database/sql"
	"fmt"
)

func deleteAddedDocRefsByDocID(db *sql.DB, addedDocRefsTableName, docID string) error {
	if _, err := db.Exec(
		fmt.Sprintf("DELETE FROM %s WHERE id LIKE $1", addedDocRefsTableName),
		fmt.Sprintf("%s%%", docID),
	); err != nil {
		return fmt.Errorf("could not perform Postgres DELETE statement: %v", err)
	}

	return nil
}

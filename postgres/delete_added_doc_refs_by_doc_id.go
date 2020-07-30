package postgres

import (
	"database/sql"
	"fmt"
)

func NewDeleteAddedDocRefsByDocID(db *sql.DB, addedDocRefsTableName string) func(docID string) error {
	return func(docID string) error {
		if _, err := db.Exec(
			fmt.Sprintf("DELETE FROM %s WHERE id LIKE $1", addedDocRefsTableName),
			fmt.Sprintf("%s%%", docID),
		); err != nil {
			return fmt.Errorf("could not perform Postgres DELETE statement: %v", err)
		}

		return nil
	}
}

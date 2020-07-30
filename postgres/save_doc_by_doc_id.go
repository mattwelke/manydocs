package postgres

import (
	"database/sql"
	"fmt"
)

func NewSaveDocByDocID(db *sql.DB, docsByDocIDTableName string) func(docJSON, docID string) error {
	return func(docJSON, docID string) error {
		SQLStatement := fmt.Sprintf("INSERT INTO %s (id, value)	VALUES ($1, $2)", docsByDocIDTableName)

		if err := db.QueryRow(SQLStatement, docID, docJSON).Scan(); err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("could not save doc in docs_by_doc_id table in Postgres: %v", err)
		}
		return nil
	}
}

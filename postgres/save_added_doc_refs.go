package postgres

import (
	"database/sql"
	"fmt"
	"github.com/mattwelke/manydocs/saveddocrefs"
	"github.com/mattwelke/manydocs/utils"
)

func saveAddedDocRefs(db *sql.DB, addedDocRefsTableName string, docInsertPrimaryKeyEntries []saveddocrefs.DocInsertPrimaryKeyEntry) error {
	for _, entry := range docInsertPrimaryKeyEntries {
		finalDocInsertPrimaryKey := fmt.Sprintf("%s%s", entry.DocID, utils.NewID())

		SQLStatement := fmt.Sprintf("INSERT INTO %s (id, value, table_name) VALUES ($1, $2, $3)", addedDocRefsTableName)

		if err := db.QueryRow(
			SQLStatement,
			finalDocInsertPrimaryKey,
			entry.PrimaryKey,
			entry.TableName,
		).Scan(); err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("could not perform Postgres query: %v", err)
		}
	}

	return nil
}

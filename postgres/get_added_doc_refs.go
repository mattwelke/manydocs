package postgres

import (
	"database/sql"
	"fmt"
	"github.com/mattwelke/manydocs/saveddocrefs"
)

func getAddedDocRefs(db *sql.DB, addedDocRefsTableName, docID string) ([]saveddocrefs.DocInsertPrimaryKeyEntry, error) {
	SQLStatement := fmt.Sprintf("SELECT value, table_name FROM %s WHERE id LIKE $1", addedDocRefsTableName)

	rows, err := db.Query(SQLStatement, fmt.Sprintf("%s%%", docID))
	if err != nil {
		return []saveddocrefs.DocInsertPrimaryKeyEntry{}, fmt.Errorf("could not perform Postgres query: %v", err)
	}
	defer rows.Close()

	primaryKeyEntries := make([]saveddocrefs.DocInsertPrimaryKeyEntry, 0)

	for rows.Next() {
		entry := saveddocrefs.DocInsertPrimaryKeyEntry{
			DocID: docID,
		}

		if err := rows.Scan(&entry.PrimaryKey, &entry.TableName); err != nil {
			return []saveddocrefs.DocInsertPrimaryKeyEntry{}, fmt.Errorf("could not scan row for doc insert primary key entry from Postgres: %v", err)
		}

		primaryKeyEntries = append(primaryKeyEntries, entry)
	}

	if err := rows.Err(); err != nil {
		return []saveddocrefs.DocInsertPrimaryKeyEntry{}, fmt.Errorf("could not complete Postgres result iteration of doc insert primary key entries during delete doc operation: %v", err)
	}

	return primaryKeyEntries, nil
}

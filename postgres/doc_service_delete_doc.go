package postgres

import (
	"fmt"
	"github.com/mattwelke/manydocs/saveddocrefs"
)

// DeleteDoc deletes a doc by its doc ID. Returns true if delete operations were needed to do this, or an error.
func (service DocService) DeleteDoc(docID string) (bool, error) {
	// Get added doc refs
	SQLStatement := fmt.Sprintf("SELECT value, table_name FROM %s WHERE id LIKE $1", addedDocRefsTable)
	rows, err := service.db.Query(SQLStatement, fmt.Sprintf("%s%%", docID))
	if err != nil {
		return false, fmt.Errorf("could not perform Postgres SELECT statement to get added doc refs: %v", err)
	}
	defer rows.Close()

	primaryKeyEntries := make([]saveddocrefs.DocInsertPrimaryKeyEntry, 0)

	for rows.Next() {
		entry := saveddocrefs.DocInsertPrimaryKeyEntry{
			DocID: docID,
		}
		if err := rows.Scan(&entry.PrimaryKey, &entry.TableName); err != nil {
			return false, fmt.Errorf("could not scan row for added doc ref from Postgres: %v", err)
		}
		primaryKeyEntries = append(primaryKeyEntries, entry)
	}

	if err := rows.Err(); err != nil {
		return false, fmt.Errorf("could not complete Postgres result iteration of added doc refs to delete doc: %v", err)
	}

	if len(primaryKeyEntries) == 0 {
		// No deletes performed
		return false, nil
	}

	for _, entry := range primaryKeyEntries {
		if _, err := service.db.Exec(
			fmt.Sprintf("DELETE FROM %s WHERE id LIKE $1", addedDocRefsTable),
			fmt.Sprintf("%s%%", entry.PrimaryKey),
		); err != nil {
			return false, fmt.Errorf("could not perform Postgres DELETE statement in added doc refs table to delete doc: %v", err)
		}
	}

	if _, err := service.db.Exec(
		fmt.Sprintf("DELETE FROM %s WHERE id LIKE $1", addedDocRefsTable),
		fmt.Sprintf("%s%%", docID),
	); err != nil {
		return false, fmt.Errorf("could not perform Postgres DELETE statement in added doc refs table to delete doc: %v", err)
	}

	return true, nil
}

package pg

import (
	"database/sql"
	"encoding/json"
	"fmt"
	mdhttp "github.com/mattwelke/manydocs/http"
	"net/http"
)

// Operation used to save a new document. An ID will be generated for the new
// document.
type deleteDocOperation struct {
	DocID string `json:"docId"`
}

func getDocInsertPrimaryKeyEntries(db *sql.DB, docID string) ([]docInsertPrimaryKeyEntry, error) {
	SQLStatement := "SELECT value, table_name FROM doc_insert_primary_keys WHERE id LIKE $1"

	rows, err := db.Query(SQLStatement, fmt.Sprintf("%s%%", docID))
	if err != nil {
		return []docInsertPrimaryKeyEntry{}, fmt.Errorf("could not get doc insert primary keys from Postgres: %v", err)
	}
	defer rows.Close()

	primaryKeyEntries := make([]docInsertPrimaryKeyEntry, 0)

	for rows.Next() {
		entry := docInsertPrimaryKeyEntry{
			docID: docID,
		}

		if err := rows.Scan(&entry.primaryKey, &entry.tableName); err != nil {
			return []docInsertPrimaryKeyEntry{}, fmt.Errorf("could not scan row for doc insert primary key entry from Postgres: %v", err)
		}

		primaryKeyEntries = append(primaryKeyEntries, entry)
	}

	if err := rows.Err(); err != nil {
		return []docInsertPrimaryKeyEntry{}, fmt.Errorf("could not complete Postgres result iteration of doc insert primary key entries during delete doc operation: %v", err)
	}

	return primaryKeyEntries, nil
}

func deleteDocsByIDPrefix(db *sql.DB, IDPrefix, tableName string) error {
	if _, err := db.Exec(
		// tableName comes from this code, this isn't unsafe
		fmt.Sprintf("DELETE FROM %s WHERE id LIKE $1", tableName),
		fmt.Sprintf("%s%%", IDPrefix),
	); err != nil {
		return fmt.Errorf("could not delete docs in %s table with ID prefix %s: %v", tableName, IDPrefix, err)
	}

	return nil
}

func deleteDocInsertPrimaryKeyEntries(db *sql.DB, docID string) error {
	if _, err := db.Exec(
		"DELETE FROM doc_insert_primary_keys WHERE id LIKE $1",
		fmt.Sprintf("%s%%", docID),
	); err != nil {
		return fmt.Errorf("could not run DELETE Postgres query: %v", err)
	}

	return nil
}

func NewDeleteDocHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var op deleteDocOperation
		if err := json.NewDecoder(r.Body).Decode(&op); err != nil {
			mdhttp.WriteBadRequest(w, fmt.Sprintf("could not decode delete doc operation: %v", err))
			return
		}

		if op.DocID == "" {
			mdhttp.WriteBadRequest(w, fmt.Sprintf("missing param %q", "docId"))
			return
		}

		primaryKeyEntries, err := getDocInsertPrimaryKeyEntries(db, op.DocID)
		if err != nil {
			mdhttp.WriteError(w, fmt.Sprintf("could not get doc insert primary key entries: %v", err))
			return
		}

		if len(primaryKeyEntries) == 0 {
			mdhttp.WriteJSON(w, mdhttp.OperationResult{
				Operation: "delete doc",
				Success:   true,
				Data: map[string]interface{}{
					"noDeleteNeeded": true,
				},
			})
			return
		}

		for _, entry := range primaryKeyEntries {
			if err := deleteDocsByIDPrefix(db, entry.primaryKey, entry.tableName); err != nil {
				mdhttp.WriteError(w, fmt.Sprintf("could not delete docs for doc insert primary key entry: %v", err))
				return
			}
		}

		if err := deleteDocInsertPrimaryKeyEntries(db, op.DocID); err != nil {
			mdhttp.WriteError(w, fmt.Sprintf("could not delete doc insert primary key entries for doc ID %s: %v", op.DocID, err))
		}

		mdhttp.WriteJSON(w, mdhttp.OperationResult{
			Operation: "delete doc",
			Success:   true,
			Data: map[string]interface{}{
				"deletedDocId": op.DocID,
			},
		})
	}
}

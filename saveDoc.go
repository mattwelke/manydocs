package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

// Operation used to save a new document. An ID will be generated for the new
// document by hashing the document's contents.
type saveDocOperation struct {
	Doc       map[string]interface{} `json:"doc"`
	DocID     string                 `json:"docId"`
	QueryKeys []map[string]string    `json:"queryKeys"`
}

func saveDocInDocsByDocIDTable(docJSON, docID string, db *sql.DB) error {
	sqlStatement := `
		INSERT INTO docs_by_doc_id (id, value)
		VALUES ($1, $2)
	`

	if err := db.QueryRow(sqlStatement, docID, docJSON).Scan(); err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("could not save doc in docs_by_doc_id table in Postgres: %v", err)
	}
	return nil
}

// Inserts the document into the "docs by query key ID" table. Returns the primary key of
// the inserted document in the underlying data store table (so that it can be saved for
// deleting the document later) or an error.
func saveDocInQueryKeysTable(docJSON, docID, docQueryKeyID string, db *sql.DB) (string, error) {
	finalQueryKeyID := fmt.Sprintf("%s%s", docQueryKeyID, newID())

	sqlStatement := `
		INSERT INTO docs_by_query_key_id (id, value)
		VALUES ($1, $2)
	`

	if err := db.QueryRow(sqlStatement, finalQueryKeyID, docJSON).Scan(); err != nil && err != sql.ErrNoRows {
		return "", fmt.Errorf("could not insert doc in Postgres: %v", err)
	}
	return finalQueryKeyID, nil
}

func saveDocInsertPrimaryKeys(docInsertPrimaryKeyEntries []docInsertPrimaryKeyEntry, db *sql.DB) error {
	for _, entry := range docInsertPrimaryKeyEntries {
		finalDocInsertPrimaryKey := fmt.Sprintf("%s%s", entry.docID, newID())

		sqlStatement := `
			INSERT INTO doc_insert_primary_keys (id, value, table_name)
			VALUES ($1, $2, $3)
		`

		if err := db.QueryRow(
			sqlStatement,
			finalDocInsertPrimaryKey,
			entry.primaryKey,
			entry.tableName,
		).Scan(); err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("could not insert doc insert primary key into Postgres: %v", err)
		}
	}

	return nil
}

func newSaveDocHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var op saveDocOperation
		if err := json.NewDecoder(r.Body).Decode(&op); err != nil {
			fmt.Fprintf(w, "could not decode save doc operation: %v", err)
			return
		}

		newDocBytes, err := json.Marshal(op.Doc)
		if err != nil {
			fmt.Fprintf(w, "could not JSON encode document: %v", err)
		}

		newDoc := string(newDocBytes)
		newDocID := op.DocID
		if newDocID == "" {
			newDocID = newID()
		}

		docInsertPrimaryKeyEntries := make([]docInsertPrimaryKeyEntry, 0)

		if err := saveDocInDocsByDocIDTable(newDoc, newDocID, db); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "could not save doc in docs_by_doc_id table in Postgres: %v", err)
			return
		}

		// Save reference to PK for delete later
		docInsertPrimaryKeyEntries = append(docInsertPrimaryKeyEntries, docInsertPrimaryKeyEntry{
			docID:      newDocID,
			tableName:  "docs_by_doc_id",
			primaryKey: newDocID,
		})

		for _, queryKeys := range op.QueryKeys {
			newDocQueryID := docQueryID(queryKeys)

			insertedPrimaryKey, err := saveDocInQueryKeysTable(newDoc, newDocID, newDocQueryID, db)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "could not save doc in docs_by_query_key_id table in Postgres: %v", err)
				return
			}

			// Save reference to PK for delete later
			docInsertPrimaryKeyEntries = append(docInsertPrimaryKeyEntries, docInsertPrimaryKeyEntry{
				docID:      newDocID,
				tableName:  "docs_by_query_key_id",
				primaryKey: insertedPrimaryKey,
			})
		}

		if err := saveDocInsertPrimaryKeys(docInsertPrimaryKeyEntries, db); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "could not save doc insert primary keys in Postgres: %v", err)
			return
		}

		writeJSON(w, operationResult{
			Operation: "save doc",
			Success:   true,
			Data: map[string]interface{}{
				"newDocId": newDocID,
			},
		})
	}
}

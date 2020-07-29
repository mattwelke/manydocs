package pg

import (
	"database/sql"
	"encoding/json"
	"fmt"
	mdhttp "github.com/mattwelke/manydocs/http"
	"github.com/mattwelke/manydocs/utils"
	"net/http"
)

// Operation used to save a new document. An ID will be generated for the new
// document by hashing the document's contents.
type saveDocOperation struct {
	Doc           map[string]interface{} `json:"doc"`
	QueryPrefixes []string               `json:"queryPrefixes"`
}

func saveDocInDocsByDocIDTable(docJSON, docID string, db *sql.DB) error {
	SQLStatement := "INSERT INTO docs_by_doc_id (id, value)	VALUES ($1, $2)"

	if err := db.QueryRow(SQLStatement, docID, docJSON).Scan(); err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("could not save doc in docs_by_doc_id table in Postgres: %v", err)
	}
	return nil
}

// Inserts the document into the "docs by query key ID" table. Returns the primary key of
// the inserted document in the underlying data store table (so that it can be saved for
// deleting the document later) or an error.
func saveDocInQueryKeysTable(docJSON, docQueryKeyID string, db *sql.DB) (string, error) {
	finalQueryKeyID := fmt.Sprintf("%s%s", docQueryKeyID, utils.NewID())

	SQLStatement := "INSERT INTO docs_by_query_key_id (id, value) VALUES ($1, $2)"

	if err := db.QueryRow(SQLStatement, finalQueryKeyID, docJSON).Scan(); err != nil && err != sql.ErrNoRows {
		return "", fmt.Errorf("could not perform Postgres query: %v", err)
	}
	return finalQueryKeyID, nil
}

func saveDocInsertPrimaryKeys(docInsertPrimaryKeyEntries []docInsertPrimaryKeyEntry, db *sql.DB) error {
	for _, entry := range docInsertPrimaryKeyEntries {
		finalDocInsertPrimaryKey := fmt.Sprintf("%s%s", entry.docID, utils.NewID())

		SQLStatement := "INSERT INTO doc_insert_primary_keys (id, value, table_name) VALUES ($1, $2, $3)"

		if err := db.QueryRow(
			SQLStatement,
			finalDocInsertPrimaryKey,
			entry.primaryKey,
			entry.tableName,
		).Scan(); err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("could not perform Postgres query: %v", err)
		}
	}

	return nil
}

func NewSaveDocHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var op saveDocOperation
		if err := json.NewDecoder(r.Body).Decode(&op); err != nil {
			_, _ = fmt.Fprintf(w, "could not decode save doc operation: %v", err)
			return
		}

		newDocID := utils.NewID()

		op.Doc["_id"] = newDocID

		newDocBytes, err := json.Marshal(op.Doc)
		if err != nil {
			_, _ = fmt.Fprintf(w, "could not JSON encode document: %v", err)
		}
		newDoc := string(newDocBytes)

		docInsertPrimaryKeyEntries := make([]docInsertPrimaryKeyEntry, 0)

		if err := saveDocInDocsByDocIDTable(newDoc, newDocID, db); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, "could not save doc in docs_by_doc_id table in Postgres: %v", err)
			return
		}

		// Save reference to PK for delete later
		docInsertPrimaryKeyEntries = append(docInsertPrimaryKeyEntries, docInsertPrimaryKeyEntry{
			docID:      newDocID,
			tableName:  "docs_by_doc_id",
			primaryKey: newDocID,
		})

		for _, queryPrefix := range op.QueryPrefixes {
			insertedPrimaryKey, err := saveDocInQueryKeysTable(newDoc, queryPrefix, db)

			if err != nil {
				mdhttp.WriteError(w, fmt.Sprintf("could not save doc in docs by query prefix table: %v", err))
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
			mdhttp.WriteError(w, fmt.Sprintf("could not save doc insert primary keys: %v", err))
			return
		}

		fmt.Printf("Saved doc with doc ID %s\n.", newDocID)
		mdhttp.WriteJSON(w, mdhttp.OperationResult{
			Operation: "save doc",
			Success:   true,
			Data: map[string]interface{}{
				"newDocId": newDocID,
			},
		})
	}
}

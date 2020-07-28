package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

// Operation used to save a new document. An ID will be generated for the new
// document.
type getDocOperation struct {
	DocID string `json:"docId"`
}

func newGetDocHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var op getDocOperation
		if err := json.NewDecoder(r.Body).Decode(&op); err != nil {
			writeBadRequest(w, fmt.Sprintf("could not decode get doc operation: %v", err))
			return
		}

		if op.DocID == "" {
			writeBadRequest(w, fmt.Sprintf("missing param %q", "docId"))
			return
		}

		SQLStatement := `
SELECT value FROM docs_by_doc_id
WHERE id = $1
		`

		var docJSON string

		if err := db.QueryRow(SQLStatement, op.DocID).Scan(&docJSON); err != nil {
			if err == sql.ErrNoRows {
				writeNotFound(w)
				return
			}
			writeError(w, fmt.Sprintf("could not get doc from Postgres: %v", err))
			return
		}

		writeJSON(w, operationResult{
			Operation: "get doc",
			Success:   true,
			Data:      docJSON,
		})
	}
}

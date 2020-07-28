package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

// Operation used to save a new document. An ID will be generated for the new
// document.
type queryDocsOperation struct {
	QueryKeys map[string]string `json:"queryKeys"`
}

func newQueryDocsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var op queryDocsOperation
		if err := json.NewDecoder(r.Body).Decode(&op); err != nil {
			writeBadRequest(w, fmt.Sprintf("could not decode query docs operation: %v", err))
			return
		}

		SQLStatement := `
SELECT value FROM docs_by_query_key_id
WHERE id LIKE $1
		`

		queryID := docQueryID(op.QueryKeys)

		rows, err := db.Query(SQLStatement, fmt.Sprintf("%s%%", queryID))
		if err != nil {
			writeError(w, fmt.Sprintf("could not query docs from Postgres: %v", err))
			return
		}
		defer rows.Close()

		docs := make([]map[string]interface{}, 0)

		for rows.Next() {
			var docStr string
			if err := rows.Scan(&docStr); err != nil {
				writeError(w, fmt.Sprintf("could not scan row for queried doc from Postgres: %v", err))
				return
			}
			var doc map[string]interface{}
			if err := json.Unmarshal([]byte(docStr), &doc); err != nil {
				writeError(w, fmt.Sprintf("could not JSON decode queried doc from Postgres: %v", err))
				return
			}
			docs = append(docs, doc)
		}

		if err := rows.Err(); err != nil {
			writeError(w, fmt.Sprintf("could not complete Postgres result iteration during query docs operation: %v", err))
			return
		}

		writeJSON(w, operationResult{
			Operation: "query docs",
			Success:   true,
			Data:      docs,
		})
	}
}

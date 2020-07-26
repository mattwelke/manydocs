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

type queryDocsRow struct {
	Value string `sql:"value"`
}

func newQueryDocsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var op queryDocsOperation
		if err := json.NewDecoder(r.Body).Decode(&op); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "could not decode query docs operation: %v", err)
			return
		}

		sqlStatement := `
SELECT value FROM docs
WHERE query_id = $1
		`

		queryID := docQueryID(op.QueryKeys)

		rows, err := db.Query(sqlStatement, queryID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "could not query docs from Postgres: %v", err)
			return
		}
		defer rows.Close()

		docs := make([]map[string]interface{}, 0)

		for rows.Next() {
			var docStr string
			if err := rows.Scan(&docStr); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "could not scan row for queried doc from Postgres: %v", err)
				return
			}
			var doc map[string]interface{}
			if err := json.Unmarshal([]byte(docStr), &doc); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "could not JSON decode queried doc from Postgres: %v", err)
				return
			}
			docs = append(docs, doc)
		}

		if err := rows.Err(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "could not complete Postgres result iteration during query docs operation: %v", err)
			return
		}

		docsJSON, err := json.Marshal(docs)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "could not JSON encode queried docs: %v", err)
		}

		writeJSON(w, docsJSON)
	}
}

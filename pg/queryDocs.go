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
type queryDocsOperation struct {
	QueryPrefix string `json:"queryPrefix"`
}

func NewQueryDocsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var op queryDocsOperation
		if err := json.NewDecoder(r.Body).Decode(&op); err != nil {
			mdhttp.WriteBadRequest(w, fmt.Sprintf("could not decode query docs operation: %v", err))
			return
		}

		SQLStatement := "SELECT value FROM docs_by_query_key_id WHERE id LIKE $1"

		rows, err := db.Query(SQLStatement, fmt.Sprintf("%s%%", op.QueryPrefix))
		if err != nil {
			mdhttp.WriteError(w, fmt.Sprintf("could not perform Postgres query: %v", err))
			return
		}
		defer rows.Close()

		docs := make([]map[string]interface{}, 0)

		for rows.Next() {
			var docStr string
			if err := rows.Scan(&docStr); err != nil {
				mdhttp.WriteError(w, fmt.Sprintf("could not scan row for queried doc from Postgres: %v", err))
				return
			}
			var doc map[string]interface{}
			if err := json.Unmarshal([]byte(docStr), &doc); err != nil {
				mdhttp.WriteError(w, fmt.Sprintf("could not JSON decode queried doc from Postgres: %v", err))
				return
			}
			docs = append(docs, doc)
		}

		if err := rows.Err(); err != nil {
			mdhttp.WriteError(w, fmt.Sprintf("could not complete Postgres result iteration during query docs operation: %v", err))
			return
		}

		fmt.Printf("Queried docs using query prefix %s with %d matching docs.\n", op.QueryPrefix, len(docs))
		mdhttp.WriteJSON(w, mdhttp.OperationResult{
			Operation: "query docs",
			Success:   true,
			Data:      docs,
		})
	}
}

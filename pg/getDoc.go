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
type getDocOperation struct {
	DocID string `json:"docId"`
}

func newGetDocHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var op getDocOperation
		if err := json.NewDecoder(r.Body).Decode(&op); err != nil {
			mdhttp.WriteBadRequest(w, fmt.Sprintf("could not decode get doc operation: %v", err))
			return
		}

		if op.DocID == "" {
			mdhttp.WriteBadRequest(w, fmt.Sprintf("missing param %q", "docId"))
			return
		}

		var docJSON string

		if err := db.QueryRow(
			"SELECT value FROM docs_by_doc_id WHERE id = $1",
			op.DocID,
		).Scan(&docJSON); err != nil {
			if err == sql.ErrNoRows {
				mdhttp.WriteNotFound(w)
				return
			}
			mdhttp.WriteError(w, fmt.Sprintf("could not get doc from Postgres: %v", err))
			return
		}

		var doc map[string]interface{}
		if err := json.Unmarshal([]byte(docJSON), &doc); err != nil {
			mdhttp.WriteError(w, fmt.Sprintf("could not JSON decode retrieved document: %v", err))
			return
		}

		mdhttp.WriteJSON(w, mdhttp.OperationResult{
			Operation: "get doc",
			Success:   true,
			Data:      doc,
		})
	}
}

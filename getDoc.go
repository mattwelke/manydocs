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
	ID string `json:"id"`
}

type getDocRow struct {
	Value string `sql:"value"`
}

func newGetDocHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var op getDocOperation
		if err := json.NewDecoder(r.Body).Decode(&op); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "could not decode get doc operation: %v", err)
			return
		}

		sqlStatement := `
SELECT value FROM docs
WHERE id = $1
		`

		var row string

		if err := db.QueryRow(sqlStatement, op.ID).Scan(&row); err != nil {
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte{})
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "could not get doc from Postgres: %v", err)
			return
		}

		writeJSON(w, []byte(row))
	}
}

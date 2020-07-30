package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

func (service DocService) GetDoc(docID string) (map[string]interface{}, error) {
	var docJSON string

	if err := service.db.QueryRow(
		fmt.Sprintf("SELECT value FROM %s WHERE id = $1", docsByDocIDTable),
		docID,
	).Scan(&docJSON); err != nil {
		if err == sql.ErrNoRows {
			// None found - valid.
			return nil, nil
		}
		return nil, fmt.Errorf("could not perform SELECT statement in Postgres to get doc: %v", err)
	}

	var doc map[string]interface{}
	if err := json.Unmarshal([]byte(docJSON), &doc); err != nil {
		return nil, fmt.Errorf("could not JSON decode document retrieved from Postgres: %v", err)
	}
	return doc, nil
}

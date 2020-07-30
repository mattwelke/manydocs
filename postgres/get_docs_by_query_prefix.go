package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

func getDocsByQueryPrefix(db *sql.DB, docsByQueryPrefixTableName, queryPrefix string) ([]map[string]interface{}, error) {
	SQLStatement := fmt.Sprintf("SELECT value FROM %s WHERE id LIKE $1", docsByQueryPrefixTableName)

	rows, err := db.Query(SQLStatement, fmt.Sprintf("%s%%", queryPrefix))
	if err != nil {
		return nil, fmt.Errorf("could not perform Postgres query: %v", err)
	}
	defer rows.Close()

	docs := make([]map[string]interface{}, 0)

	for rows.Next() {
		var docStr string
		if err := rows.Scan(&docStr); err != nil {
			return nil, fmt.Errorf("could not scan row for queried doc from Postgres: %v", err)
		}
		var doc map[string]interface{}
		if err := json.Unmarshal([]byte(docStr), &doc); err != nil {
			return nil, fmt.Errorf("could not JSON decode queried doc from Postgres: %v", err)
		}
		docs = append(docs, doc)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("could not complete Postgres result iteration during query docs operation: %v", err)
	}
	return docs, nil
}

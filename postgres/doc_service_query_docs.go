package postgres

import (
	"encoding/json"
	"fmt"
)

func (service DocService) QueryDocs(queryPrefix string) ([]map[string]interface{}, error) {
	SQLStatement := fmt.Sprintf("SELECT value FROM %s WHERE id LIKE $1", docsByQueryPrefixTable)

	rows, err := service.db.Query(SQLStatement, fmt.Sprintf("%s%%", queryPrefix))
	if err != nil {
		return nil, fmt.Errorf("could not perform Postgres SELECT statement to query docs: %v", err)
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

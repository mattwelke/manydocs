package postgres

import (
	"database/sql"
	"fmt"
	"github.com/mattwelke/manydocs/utils"
)

func NewSaveDocByQueryPrefix(db *sql.DB, docsByQueryPrefixTableName string) func(docJSON, docQueryKeyID string) (string, error) {
	return func(docJSON, docQueryKeyID string) (string, error) {
		finalQueryKeyID := fmt.Sprintf("%s%s", docQueryKeyID, utils.NewID())

		SQLStatement := fmt.Sprintf("INSERT INTO %s (id, value) VALUES ($1, $2)", docsByQueryPrefixTableName)

		if err := db.QueryRow(SQLStatement, finalQueryKeyID, docJSON).Scan(); err != nil && err != sql.ErrNoRows {
			return "", fmt.Errorf("could not perform Postgres query: %v", err)
		}
		return finalQueryKeyID, nil
	}
}

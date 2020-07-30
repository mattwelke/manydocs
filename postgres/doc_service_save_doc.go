package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/mattwelke/manydocs/saveddocrefs"
	"github.com/mattwelke/manydocs/utils"
)

func (service DocService) SaveDoc(newDoc map[string]interface{}, newDocQueryPrefixes []string) (string, error) {
	// Generate and append needed data
	newDocID := utils.NewID()
	newDoc["_id"] = newDocID

	// JSON encode for inserts
	docJSON, err := json.Marshal(newDoc)
	if err != nil {
		return "", fmt.Errorf("could not JSON encode new doc: %v", err)
	}

	docInsertPrimaryKeyEntries := make([]saveddocrefs.DocInsertPrimaryKeyEntry, 0)

	// Insert into docs by doc ID table in Postgres
	SQLStatement := fmt.Sprintf("INSERT INTO %s (id, value)	VALUES ($1, $2)", docsByDocIDTable)
	if err := service.db.QueryRow(SQLStatement, newDocID, string(docJSON)).Scan(); err != nil && err != sql.ErrNoRows {
		return "", fmt.Errorf("could not perform Postgres INSERT INTO statement to save doc in docs by doc ID table: %v", err)
	}

	// Save ref to doc in docs by doc ID table for "delete doc" operation later
	docInsertPrimaryKeyEntries = append(docInsertPrimaryKeyEntries, saveddocrefs.DocInsertPrimaryKeyEntry{
		DocID:      newDocID,
		TableName:  docsByDocIDTable,
		PrimaryKey: newDocID,
	})

	for _, prefix := range newDocQueryPrefixes {
		finalQueryKeyID := fmt.Sprintf("%s%s", prefix, utils.NewID())

		SQLStatement := fmt.Sprintf("INSERT INTO %s (id, value) VALUES ($1, $2)", docsByQueryPrefixTable)

		if err := service.db.QueryRow(SQLStatement, finalQueryKeyID, docJSON).Scan(); err != nil && err != sql.ErrNoRows {
			return "", fmt.Errorf("could not perform Postgres INSERT INTO statement to save doc in docs by query prefix table: %v", err)
		}

		// Save reference to PK for delete later
		docInsertPrimaryKeyEntries = append(docInsertPrimaryKeyEntries, saveddocrefs.DocInsertPrimaryKeyEntry{
			DocID:      newDocID,
			TableName:  addedDocRefsTable,
			PrimaryKey: finalQueryKeyID,
		})
	}

	for _, entry := range docInsertPrimaryKeyEntries {
		finalDocInsertPrimaryKey := fmt.Sprintf("%s%s", entry.DocID, utils.NewID())

		SQLStatement := fmt.Sprintf("INSERT INTO %s (id, value, table_name) VALUES ($1, $2, $3)", addedDocRefsTable)

		if err := service.db.QueryRow(
			SQLStatement,
			finalDocInsertPrimaryKey,
			entry.PrimaryKey,
			entry.TableName,
		).Scan(); err != nil && err != sql.ErrNoRows {
			return "", fmt.Errorf("could not perform Postgres INSERT INTO statement to save added doc ref: %v", err)
		}
	}

	return newDocID, nil
}

package doc

import (
	"database/sql"
	"github.com/mattwelke/manydocs/postgres"
	"github.com/mattwelke/manydocs/saveddocrefs"
)

type Saver struct {
	saveDocByDocID       func(docJSON, docID string) error
	saveDocByQueryPrefix func(docJSON, docQueryPrefix string) (string, error)
	addSavedDocRefs      func(entries []saveddocrefs.DocInsertPrimaryKeyEntry) error
}

func NewPostgresSaver(db *sql.DB) Saver {
	return Saver{
		saveDocByDocID: func(docJSON, docID string) error {
			return postgres.SaveDocInDocsByDocIDTable(docJSON, docID, db)
		},
		saveDocByQueryPrefix: func(docJSON, docID string) (string, error) {
			return postgres.SaveDocInQueryKeysTable(docJSON, docID, db)
		},
		addSavedDocRefs: func(entries []saveddocrefs.DocInsertPrimaryKeyEntry) error {
			return postgres.SaveDocInsertPrimaryKeys(entries, db)
		},
	}
}

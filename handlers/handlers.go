package handlers

import "github.com/mattwelke/manydocs/saveddocrefs"

// DocService is able to perform operations for docs.
type DocService interface {
	// for "save doc"
	SaveDocByDocID(docJSON, docID string) error
	SaveDocByQueryPrefix(docJSON, docQueryPrefix string) (string, error)
	SaveAddedDocRefs(entries []saveddocrefs.DocInsertPrimaryKeyEntry) error

	// for "get doc"
	GetDocByDocID(docID string) (map[string]interface{}, error)

	// for "query docs"
	GetDocsByQueryPrefix(queryPrefix string) ([]map[string]interface{}, error)

	// for "delete doc"
	GetAddedDocRefs(docID string) ([]saveddocrefs.DocInsertPrimaryKeyEntry, error)
	DeleteDocsByPrimaryKeyPrefix(primaryKeyPrefix, tableName string) error
	DeleteAddedDocRefsByDocID(docID string) error
}

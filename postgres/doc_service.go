package postgres

import (
	"database/sql"
	"github.com/mattwelke/manydocs/saveddocrefs"
)

type TableNames struct {
	DocsByDocID       string
	DocsByQueryPrefix string
	AddedDocRefs      string
}

// DocService performs operations using Postgres as a storage engine.
type DocService struct {
	db         *sql.DB
	tableNames TableNames
}

// NewDocService creates a new DocService.
func NewDocService(db *sql.DB, tableNames TableNames) DocService {
	return DocService{
		db,
		tableNames,
	}
}

// for "save doc"
func (service DocService) SaveDocByDocID(docJSON, docID string) error {
	return saveDocByDocID(service.db, service.tableNames.DocsByDocID, docJSON, docID)
}
func (service DocService) SaveDocByQueryPrefix(docJSON, docQueryPrefix string) (string, error) {
	return saveDocByQueryPrefix(service.db, service.tableNames.DocsByQueryPrefix, docJSON, docQueryPrefix)
}
func (service DocService) SaveAddedDocRefs(entries []saveddocrefs.DocInsertPrimaryKeyEntry) error {
	return saveAddedDocRefs(service.db, service.tableNames.AddedDocRefs, entries)
}

// for "get doc"
func (service DocService) GetDocByDocID(docID string) (map[string]interface{}, error) {
	return getDocByDocID(service.db, service.tableNames.DocsByDocID, docID)
}

// for "query docs"
func (service DocService) GetDocsByQueryPrefix(queryPrefix string) ([]map[string]interface{}, error) {
	return getDocsByQueryPrefix(service.db, service.tableNames.DocsByQueryPrefix, queryPrefix)
}

// for "delete doc"
func (service DocService) GetAddedDocRefs(docID string) ([]saveddocrefs.DocInsertPrimaryKeyEntry, error) {
	return getAddedDocRefs(service.db, service.tableNames.AddedDocRefs, docID)
}
func (service DocService) DeleteDocsByPrimaryKeyPrefix(primaryKeyPrefix, tableName string) error {
	return deleteDocsByPrimaryKeyPrefix(service.db, primaryKeyPrefix, tableName)
}
func (service DocService) DeleteAddedDocRefsByDocID(docID string) error {
	return deleteAddedDocRefsByDocID(service.db, service.tableNames.AddedDocRefs, docID)
}

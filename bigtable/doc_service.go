package bigtable

import (
	"cloud.google.com/go/bigtable"
)

// DocService performs operations using Google Cloud Bigtable as a storage engine.
type DocService struct {
	docsByDocsIDTable      *bigtable.Table
	docsByQueryPrefixTable *bigtable.Table
	addedDocRefsTable      *bigtable.Table
}

// NewDocService creates a new DocService.
func NewDocService(docsByDocsIDTable, docsByQueryPrefixTable, addedDocRefsTable *bigtable.Table) DocService {
	return DocService{
		docsByDocsIDTable,
		docsByQueryPrefixTable,
		addedDocRefsTable,
	}
}

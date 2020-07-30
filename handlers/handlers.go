package handlers

// DocService is able to perform operations for docs.
type DocService interface {
	SaveDoc(newDoc map[string]interface{}, newDocQueryPrefixes []string) (string, error)
	GetDoc(docID string) (map[string]interface{}, error)
	QueryDocs(queryPrefix string) ([]map[string]interface{}, error)
	DeleteDoc(docID string) (bool, error)
}

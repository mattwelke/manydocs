package operations

// SaveDoc is the operation used to save a new document. An ID will be
// generated for the new document by hashing the document's contents.
type SaveDoc struct {
	Doc           map[string]interface{} `json:"doc"`
	QueryPrefixes []string               `json:"queryPrefixes"`
}

// GetDoc is the operation to get a document by its doc ID that was
// automatically-generation when the document was saved.
type GetDoc struct {
	DocID string `json:"docId"`
}

// QueryDocs is the operation used to get a set of documents that match a
// particular query prefix.
type QueryDocs struct {
	QueryPrefix string `json:"queryPrefix"`
}

// DeleteDoc is the operation used to delete a document by its doc ID that was
// automatically-generated when the document was saved.
type DeleteDoc struct {
	DocID string `json:"docId"`
}

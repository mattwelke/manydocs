package utils

import (
	"fmt"
	"github.com/google/uuid"
)

const (
	queryKeyDelim = "#"
)

func NewID() string {
	return uuid.New().String()
}

// Forms a new document's "query ID", which is the ID used for "queries", which
// are retrieval operations that retrieve multiple documents by query ID
// prefix.
func DocQueryID(docQueryKeys map[string]string) string {
	ID := ""
	for keyName, keyValue := range docQueryKeys {
		ID = fmt.Sprintf("%s%s%s=%s", ID, queryKeyDelim, keyName, keyValue)
	}
	return ID
}

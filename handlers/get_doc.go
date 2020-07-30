package handlers

import (
	"encoding/json"
	"fmt"
	mdhttp "github.com/mattwelke/manydocs/http"
	"net/http"
)

// Operation used to save a new document. An ID will be generated for the new
// document.
type getDocOperation struct {
	DocID string `json:"docId"`
}

func NewGetDocHandler(
	getDocByDocID func(docID string) (map[string]interface{}, error),
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var op getDocOperation
		if err := json.NewDecoder(r.Body).Decode(&op); err != nil {
			mdhttp.WriteBadRequest(w, fmt.Sprintf("could not decode get doc operation: %v", err))
			return
		}

		if op.DocID == "" {
			mdhttp.WriteBadRequest(w, fmt.Sprintf("missing param %q", "docId"))
			return
		}

		doc, err := getDocByDocID(op.DocID)
		if err != nil {
			mdhttp.WriteError(w, fmt.Sprintf("could not get document: %v", err))
			return
		}
		if doc == nil {
			mdhttp.WriteNotFound(w)
			return
		}

		fmt.Printf("Retrieved doc with doc ID %s\n.", op.DocID)
		mdhttp.WriteJSON(w, mdhttp.OperationResult{
			Operation: "get doc",
			Success:   true,
			Data:      doc,
		})
	}
}

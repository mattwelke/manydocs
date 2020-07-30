package handlers

import (
	"encoding/json"
	"fmt"
	mdhttp "github.com/mattwelke/manydocs/http"
	"github.com/mattwelke/manydocs/saveddocrefs"
	"github.com/mattwelke/manydocs/utils"
	"net/http"
)

// Operation used to save a new document. An ID will be generated for the new
// document by hashing the document's contents.
type saveDocOperation struct {
	Doc           map[string]interface{} `json:"doc"`
	QueryPrefixes []string               `json:"queryPrefixes"`
}

func NewSaveDocHandler(
	docService DocService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var op saveDocOperation
		if err := json.NewDecoder(r.Body).Decode(&op); err != nil {
			_, _ = fmt.Fprintf(w, "could not decode save doc operation: %v", err)
			return
		}

		newDocID := utils.NewID()

		op.Doc["_id"] = newDocID

		newDocBytes, err := json.Marshal(op.Doc)
		if err != nil {
			_, _ = fmt.Fprintf(w, "could not JSON encode document: %v", err)
		}
		newDoc := string(newDocBytes)

		docInsertPrimaryKeyEntries := make([]saveddocrefs.DocInsertPrimaryKeyEntry, 0)

		if err := docService.SaveDocByDocID(newDoc, newDocID); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, "could not save doc in docs_by_doc_id table in Postgres: %v", err)
			return
		}

		// Save reference to PK for delete later
		docInsertPrimaryKeyEntries = append(docInsertPrimaryKeyEntries, saveddocrefs.DocInsertPrimaryKeyEntry{
			DocID:      newDocID,
			TableName:  "docs_by_doc_id",
			PrimaryKey: newDocID,
		})

		for _, queryPrefix := range op.QueryPrefixes {
			insertedPrimaryKey, err := docService.SaveDocByQueryPrefix(newDoc, queryPrefix)

			if err != nil {
				mdhttp.WriteError(w, fmt.Sprintf("could not save doc in docs by query prefix table: %v", err))
				return
			}

			// Save reference to PK for delete later
			docInsertPrimaryKeyEntries = append(docInsertPrimaryKeyEntries, saveddocrefs.DocInsertPrimaryKeyEntry{
				DocID:      newDocID,
				TableName:  "docs_by_query_key_id",
				PrimaryKey: insertedPrimaryKey,
			})
		}

		if err := docService.SaveAddedDocRefs(docInsertPrimaryKeyEntries); err != nil {
			mdhttp.WriteError(w, fmt.Sprintf("could not save doc insert primary keys: %v", err))
			return
		}

		fmt.Printf("Saved doc with doc ID %s\n.", newDocID)
		mdhttp.WriteJSON(w, mdhttp.OperationResult{
			Operation: "save doc",
			Success:   true,
			Data: map[string]interface{}{
				"newDocId": newDocID,
			},
		})
	}
}

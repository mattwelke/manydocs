package handlers

import (
	"encoding/json"
	"fmt"
	mdhttp "github.com/mattwelke/manydocs/http"
	"net/http"
)

type deleteDocOperation struct {
	DocID string `json:"docId"`
}

func NewDeleteDocHandler(
	docService DocService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var op deleteDocOperation
		if err := json.NewDecoder(r.Body).Decode(&op); err != nil {
			mdhttp.WriteBadRequest(w, fmt.Sprintf("could not decode delete doc operation: %v", err))
			return
		}

		if op.DocID == "" {
			mdhttp.WriteBadRequest(w, fmt.Sprintf("missing param %q", "docId"))
			return
		}

		primaryKeyEntries, err := docService.GetAddedDocRefs(op.DocID)
		if err != nil {
			mdhttp.WriteError(w, fmt.Sprintf("could not get doc insert primary key entries: %v", err))
			return
		}

		if len(primaryKeyEntries) == 0 {
			mdhttp.WriteJSON(w, mdhttp.OperationResult{
				Operation: "delete doc",
				Success:   true,
				Data: map[string]interface{}{
					"noDeleteNeeded": true,
				},
			})
			return
		}

		for _, entry := range primaryKeyEntries {
			if err := docService.DeleteDocsByPrimaryKeyPrefix(entry.PrimaryKey, entry.TableName); err != nil {
				mdhttp.WriteError(w, fmt.Sprintf("could not delete docs for doc insert primary key entry: %v", err))
				return
			}
		}

		if err := docService.DeleteAddedDocRefsByDocID(op.DocID); err != nil {
			mdhttp.WriteError(w, fmt.Sprintf("could not delete doc insert primary key entries for doc ID %s: %v", op.DocID, err))
		}

		fmt.Printf("Deleted doc with doc ID %s\n.", op.DocID)
		mdhttp.WriteJSON(w, mdhttp.OperationResult{
			Operation: "delete doc",
			Success:   true,
			Data: map[string]interface{}{
				"deletedDocId": op.DocID,
			},
		})
	}
}

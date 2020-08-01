package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	mdhttp "github.com/mattwelke/manydocs/http"
	"github.com/mattwelke/manydocs/operations"
)

func NewSaveDocHandler(
	docService DocService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var op operations.SaveDoc
		if err := json.NewDecoder(r.Body).Decode(&op); err != nil {
			_, _ = fmt.Fprintf(w, "could not decode save doc operation: %v", err)
			return
		}

		newDocID, err := docService.SaveDoc(op.Doc, op.QueryPrefixes)
		if err != nil {
			fmt.Printf("could not perform save doc operation: %v\n", err)
			mdhttp.WriteError(w, "could not perform save doc operation")
			return
		}

		fmt.Printf("saved doc with doc ID %s\n", newDocID)
		mdhttp.WriteJSON(w, mdhttp.OperationResult{
			Operation: "save doc",
			Success:   true,
			Data: map[string]interface{}{
				"newDocId": newDocID,
			},
		}, 0)
	}
}

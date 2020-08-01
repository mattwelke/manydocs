package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	mdhttp "github.com/mattwelke/manydocs/http"
	"github.com/mattwelke/manydocs/operations"
)

func NewDeleteDocHandler(
	docService DocService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var op operations.DeleteDoc
		if err := json.NewDecoder(r.Body).Decode(&op); err != nil {
			mdhttp.WriteBadRequest(w, fmt.Sprintf("could not decode delete doc operation: %v", err))
			return
		}

		if op.DocID == "" {
			mdhttp.WriteBadRequest(w, fmt.Sprintf("missing param %q", "docId"))
			return
		}

		deletesPerformed, err := docService.DeleteDoc(op.DocID)
		if err != nil {
			fmt.Printf("could not perform delete doc operation: %v\n", err)
			mdhttp.WriteError(w, "could not perform delete doc operation")
			return
		}

		var deletesPerformedText string
		if deletesPerformed {
			deletesPerformedText = "deletes"
		} else {
			deletesPerformedText = "no deletes"
		}

		fmt.Printf("%s performed in storage engine for delete doc operation for doc ID %s\n", deletesPerformedText, op.DocID)
		mdhttp.WriteJSON(w, mdhttp.OperationResult{
			Operation: "delete doc",
			Success:   true,
			Data: map[string]interface{}{
				"deletedDocId":                    op.DocID,
				"deletesPerformedInStorageEngine": deletesPerformed,
			},
		}, 0)
	}
}

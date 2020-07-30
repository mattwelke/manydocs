package handlers

import (
	"encoding/json"
	"fmt"
	mdhttp "github.com/mattwelke/manydocs/http"
	"github.com/mattwelke/manydocs/operations"
	"net/http"
)

func NewGetDocHandler(
	docService DocService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var op operations.GetDoc
		if err := json.NewDecoder(r.Body).Decode(&op); err != nil {
			mdhttp.WriteBadRequest(w, fmt.Sprintf("could not decode get doc operation: %v", err))
			return
		}

		if op.DocID == "" {
			mdhttp.WriteBadRequest(w, fmt.Sprintf("missing param %q", "docId"))
			return
		}

		doc, err := docService.GetDoc(op.DocID)
		if err != nil {
			fmt.Printf("could not perform get doc operation: %v\n", err)
			mdhttp.WriteError(w, "could not perform get doc operation")
			return
		}
		if doc == nil {
			fmt.Printf("doc with doc ID %s not found during get doc operation\n", op.DocID)
			mdhttp.WriteNotFound(w)
			return
		}

		fmt.Printf("doc with doc ID %s found during get doc operation\n", op.DocID)
		mdhttp.WriteJSON(w, mdhttp.OperationResult{
			Operation: "get doc",
			Success:   true,
			Data:      doc,
		})
	}
}

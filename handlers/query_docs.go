package handlers

import (
	"encoding/json"
	"fmt"
	mdhttp "github.com/mattwelke/manydocs/http"
	"net/http"
)

type queryDocsOperation struct {
	QueryPrefix string `json:"queryPrefix"`
}

func NewQueryDocsHandler(
	docService DocService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var op queryDocsOperation
		if err := json.NewDecoder(r.Body).Decode(&op); err != nil {
			mdhttp.WriteBadRequest(w, fmt.Sprintf("could not decode query docs operation: %v", err))
			return
		}

		if op.QueryPrefix == "" {
			mdhttp.WriteBadRequest(w, fmt.Sprintf("missing param %q", "queryPrefix"))
			return
		}

		docs, err := docService.GetDocsByQueryPrefix(op.QueryPrefix)
		if err != nil {
			mdhttp.WriteError(w, fmt.Sprintf("could not get docs by query prefix: %v", err))
			return
		}

		fmt.Printf("Queried docs using query prefix %s with %d matching docs.\n", op.QueryPrefix, len(docs))

		mdhttp.WriteJSON(w, mdhttp.OperationResult{
			Operation: "query docs",
			Success:   true,
			Data:      docs,
		})
	}
}

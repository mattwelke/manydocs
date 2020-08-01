package handlers

import (
	"fmt"
	"net/http"

	mdhttp "github.com/mattwelke/manydocs/http"
	"github.com/mattwelke/manydocs/operations"
)

func NewQueryDocsHandler(
	docService DocService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		op := operations.QueryDocs{
			QueryPrefix: r.URL.Query().Get("p"),
		}

		if op.QueryPrefix == "" {
			mdhttp.WriteBadRequest(w, fmt.Sprintf("missing param %q", "queryPrefix"))
			return
		}

		docs, err := docService.QueryDocs(op.QueryPrefix)
		if err != nil {
			fmt.Printf("could not perform query docs operation: %v", err)
			mdhttp.WriteError(w, "could not perform query docs operation")
			return
		}

		fmt.Printf("%d docs found during query docs operation with query prefix %s\n", len(docs), op.QueryPrefix)
		mdhttp.WriteJSON(w, mdhttp.OperationResult{
			Operation: "query docs",
			Success:   true,
			Data:      docs,
		}, 60)
	}
}

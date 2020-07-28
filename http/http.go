package http

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	poweredBy = "manydocs"
)

type OperationResult struct {
	Operation string      `json:"operation"`
	Success   bool        `json:"success"`
	Data      interface{} `json:"data,omitempty"`
}

func WriteJSON(w http.ResponseWriter, result operationResult) {
	resultJSON, err := json.Marshal(result)
	if err != nil {
		WriteError(w, "could not marshal operation result: %v")
		return
	}
	w.Header().Add("X-Powered-By", poweredBy)
	w.Header().Add("Content-Type", "application/json")
	_, _ = w.Write(resultJSON)
}

func WriteNotFound(w http.ResponseWriter) {
	w.Header().Add("X-Powered-By", poweredBy)
	w.WriteHeader(http.StatusNotFound)
	_, _ = w.Write([]byte{})
}

func WriteError(w http.ResponseWriter, msg string) {
	w.Header().Add("X-Powered-By", poweredBy)
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = fmt.Fprintf(w, msg)
}

func WriteBadRequest(w http.ResponseWriter, msg string) {
	w.Header().Add("X-Powered-By", poweredBy)
	w.WriteHeader(http.StatusBadRequest)
	_, _ = fmt.Fprintf(w, msg)
}

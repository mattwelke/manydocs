package main

// Helpers for receiving requests and sending responses.

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	poweredBy = "manydocs"
)

type operationResult struct {
	Operation string      `json:"operation"`
	Success   bool        `json:"success"`
	Data      interface{} `json:"data,omitempty"`
}

func writeJSON(w http.ResponseWriter, result operationResult) {
	resultJSON, err := json.Marshal(result)
	if err != nil {
		writeError(w, "could not marshal operation result: %v")
		return
	}
	w.Header().Add("X-Powered-By", poweredBy)
	w.Header().Add("Content-Type", "application/json")
	w.Write(resultJSON)
}

func writeNotFound(w http.ResponseWriter) {
	w.Header().Add("X-Powered-By", poweredBy)
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte{})
}

func writeError(w http.ResponseWriter, msg string) {
	w.Header().Add("X-Powered-By", poweredBy)
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, msg)
}

func writeBadRequest(w http.ResponseWriter, msg string) {
	w.Header().Add("X-Powered-By", poweredBy)
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, msg)
}

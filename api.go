package main

// Helpers for receiving requests and sending responses.

import "net/http"

const (
	poweredBy = "manydocs"
)

func writeJSON(w http.ResponseWriter, res []byte) {
	w.Header().Add("X-Powered-By", poweredBy)
	w.Header().Add("Content-Type", "application/json")
	w.Write(res)
}

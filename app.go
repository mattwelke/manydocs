package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	_ "github.com/lib/pq"
)

const (
	host     = "hanno.db.elephantsql.com"
	port     = 5432
	user     = "ylbpptjl"
	password = "SwPrRoue2wlI7_nHuPQg2PtrZWqSa2Wb"
	dbname   = "ylbpptjl"

	queryKeyDelim = "#"
)

func newID() string {
	return uuid.New().String()
}

// Forms a new document's "query ID", which is the ID used for "queries", which
// are retrieval operations that retrieve multiple documents by query ID
// prefix.
func docQueryID(docQueryKeys map[string]string) string {
	ID := ""
	for keyName, keyValue := range docQueryKeys {
		ID = fmt.Sprintf("%s%s%s=%s", ID, queryKeyDelim, keyName, keyValue)
	}
	return ID
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/save", newSaveDocHandler(db))
	http.HandleFunc("/get", newGetDocHandler(db))
	http.HandleFunc("/query", newQueryDocsHandler(db))

	http.ListenAndServe("localhost:8080", nil)
}

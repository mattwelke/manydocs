package main

import (
	"database/sql"
	"fmt"
	"github.com/mattwelke/manydocs/pg"
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
)





func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/save", pg.newSaveDocHandler(db))
	http.HandleFunc("/get", pg.newGetDocHandler(db))
	http.HandleFunc("/query", pg.newQueryDocsHandler(db))
	http.HandleFunc("/delete", pg.newDeleteDocHandler(db))

	_ = http.ListenAndServe("localhost:8080", nil)
}

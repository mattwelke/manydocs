package main

import (
	"cloud.google.com/go/bigtable"
	"context"
	"database/sql"
	"fmt"
	mdbigtable "github.com/mattwelke/manydocs/bigtable"
	"github.com/mattwelke/manydocs/handlers"
	"github.com/mattwelke/manydocs/postgres"
	"net/http"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

const (
	host     = "otto.db.elephantsql.com"
	port     = 5432
	user     = "bsmeuqtn"
	password = "rHBDPgnK9ndrsOXMYs3mthmnYRNhnytA"
	dbname   = "bsmeuqtn"
)

func main() {
	serverPort, _ := strconv.ParseInt(os.Getenv("PORT"), 10, 32)
	storageEngine := os.Getenv("STORAGE_ENGINE")

	var docService handlers.DocService
	if storageEngine == "POSTGRES" {
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)
		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			panic(err)
		}

		db.SetMaxOpenConns(1)
		db.SetMaxIdleConns(1)

		docService = postgres.NewDocService(db)
	} else if storageEngine == "BIGTABLE" {
		client, err := bigtable.NewClient(
			context.Background(),
			"matt-welke-manydocs",
			"manydocs",
		)
		if err != nil {
			panic(fmt.Sprintf("could not create Bigtable client: %v", err))
		}
		docService = mdbigtable.NewDocService(
			client.Open("docs_by_doc_id"),
			client.Open("docs_by_query_prefix"),
			client.Open("added_doc_refs"),
		)
	} else {
		panic("storage engine must be POSTGRES or BIGTABLE")
	}

	http.HandleFunc("/save", handlers.NewSaveDocHandler(docService))
	http.HandleFunc("/get", handlers.NewGetDocHandler(docService))
	http.HandleFunc("/query", handlers.NewQueryDocsHandler(docService))
	http.HandleFunc("/delete", handlers.NewDeleteDocHandler(docService))

	_ = http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", serverPort), nil)
}

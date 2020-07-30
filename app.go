package main

import (
	"database/sql"
	"fmt"
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

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	postgresDocService := postgres.NewDocService(db)

	http.HandleFunc("/save", handlers.NewSaveDocHandler(postgresDocService))
	http.HandleFunc("/get", handlers.NewGetDocHandler(postgresDocService))
	http.HandleFunc("/query", handlers.NewQueryDocsHandler(postgresDocService))
	http.HandleFunc("/delete", handlers.NewDeleteDocHandler(postgresDocService))

	_ = http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", serverPort), nil)
}

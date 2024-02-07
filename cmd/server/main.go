package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"tn-rest/internal/db"

	_ "github.com/mattn/go-sqlite3"
)

func GetConnection() (context.Context, *db.Queries) {
	ctx := context.Background()
	dbx, err := sql.Open("sqlite3", "national_park.db")

	if err != nil {
		panic(err)
	}

	queries := db.New(dbx)

	return ctx, queries
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Welcome to Indonesia National Park API")
	})

	log.Fatal(http.ListenAndServe(":8080", mux))
}

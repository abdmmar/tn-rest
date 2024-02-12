package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"tn-rest/cmd/server/handler"
	"tn-rest/cmd/server/service"
	"tn-rest/pkg/router"

	_ "github.com/mattn/go-sqlite3"
)

func teardown(db *sql.DB) {
	if db != nil {
		if err := db.Close(); err != nil {
			log.Fatal("failed to close database connection", err.Error())
		} else {
			log.Fatal("successfully closed database connection")
		}
	}
}

func main() {
	ctx := context.Background()
	dbx, err := sql.Open("sqlite3", "national_park.db")

	defer teardown(dbx)
	if err != nil {
		panic(err)
	}

	route := router.NewRouter()

	nationalPark := handler.NationalParkHandler{&service.NationalParkService{ctx, dbx}}
	route.GET("/api/national_park", nationalPark.GetAll)
	route.POST("/api/national_park", nationalPark.Create)

	log.Printf("🚀 server listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", route))
}

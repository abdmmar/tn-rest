package main

import (
	"context"
	"database/sql"
	"fmt"
	"tn-rest/internal/db"
	"tn-rest/internal/sqlc"

	_ "github.com/mattn/go-sqlite3"
)

func prepareDatabase(dbx db.DBTX) {
	ctx := context.Background()
	if _, err := dbx.ExecContext(ctx, sqlc.DDL); err != nil {
		fmt.Println("err", err)
	}
}

func main(){
	dbx, err := sql.Open("sqlite3", "national_park.db")
	if err != nil {
		panic(err)
	}
	prepareDatabase(dbx)
}
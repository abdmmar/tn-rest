package main

import (
	"context"
	"database/sql"
	"net/http"
	"regexp"
	"tn-rest/internal/db"
)

type Route struct {
	method  string
	regex   *regexp.Regexp
	handler http.HandlerFunc
}

func route(method, pattern string, handler http.HandlerFunc) Route {
	return Route{method, regexp.MustCompile("^" + pattern + "$"), handler}
}

type Router struct {
	routes []*Route
}

func NewRouter(ctx context.Context, queries *db.Queries, db *sql.DB) []Route {
	nationalParkHandler := NationalParkHandler{Ctx: ctx, Queries: queries, DB: db}

	var routes = []Route{
		route("GET", "/api/national_park", nationalParkHandler.Create),
	}

	return routes
}

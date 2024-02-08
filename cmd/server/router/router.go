package router

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"tn-rest/pkg/validate"
)

type Route struct {
	method     string
	pattern    *regexp.Regexp
	handler    http.HandlerFunc
	paramKeys  []string
	middleware []http.HandlerFunc
}

type Router struct {
	routes []*Route
}

func NewRouter() *Router {
	return &Router{routes: []*Route{}}
}

func (r *Router) AddRoute(method, path string, handler http.HandlerFunc, middleware ...http.HandlerFunc) *Route {
	pathParamMatcher := regexp.MustCompile(":([a-zA-Z]+)")
	matches := pathParamMatcher.FindAllStringSubmatch(path, -1)
	paramKeys := []string{}
	pattern := path

	if len(matches) > 0 {
		pattern = pathParamMatcher.ReplaceAllLiteralString(path, "([^/]+)")
		for i, _ := range matches {
			paramKeys = append(paramKeys, matches[i][1])
		}
	}

	if !validate.Include(paramKeys) {
		panic(fmt.Sprint("duplicate parameters in path definition: ", path))
	}

	// check duplicate route with same method and path pattern
	regex := regexp.MustCompile("^" + pattern + "$")
	for _, route := range r.routes {
		if route.method == method && route.pattern.String() == regex.String() {
			panic(fmt.Sprintf("route already exists: %s %s", method, path))
		}
	}

	newRoute := &Route{
		method,
		regex,
		handler,
		paramKeys,
		middleware,
	}
	r.routes = append(r.routes, newRoute)

	return newRoute
}

func (r *Router) GET(pattern string, handler http.HandlerFunc, middleware ...http.HandlerFunc) *Route {
	return r.AddRoute(http.MethodGet, pattern, handler, middleware...)
}

func (r *Router) POST(pattern string, handler http.HandlerFunc, middleware ...http.HandlerFunc) *Route {
	return r.AddRoute(http.MethodPost, pattern, handler, middleware...)
}

func (r *Router) PUT(pattern string, handler http.HandlerFunc, middleware ...http.HandlerFunc) *Route {
	return r.AddRoute(http.MethodPut, pattern, handler, middleware...)
}

func (r *Router) PATCH(pattern string, handler http.HandlerFunc, middleware ...http.HandlerFunc) *Route {
	return r.AddRoute(http.MethodPatch, pattern, handler, middleware...)
}

func (r *Router) DELETE(pattern string, handler http.HandlerFunc, middleware ...http.HandlerFunc) *Route {
	return r.AddRoute(http.MethodDelete, pattern, handler, middleware...)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var allow []string

	for _, route := range r.routes {
		matches := route.pattern.FindStringSubmatch(req.URL.Path)
		if len(matches) > 0 {
			if req.Method != route.method {
				allow = append(allow, route.method)
				continue
			}

			values := matches[1:]
			if len(values) != len(route.paramKeys) {
				errorResponse := "unexpected number of path parameters in request"
				fmt.Sprintf("%s (%s)", errorResponse, req.URL.Path)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(errorResponse))
				return
			}

			route.handler(w, buildContext(req, route.paramKeys, values))
			return
		}
	}

	if len(allow) > 0 {
		w.Header().Set("Allow", strings.Join(allow, ", "))
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	http.NotFound(w, req)
}

type ContextKey string

func buildContext(req *http.Request, paramKeys, paramValues []string) *http.Request {
	ctx := req.Context()
	for i, _ := range paramKeys {
		ctx = context.WithValue(ctx, ContextKey(paramKeys[i]), paramValues[i])
	}
	return req.WithContext(ctx)
}

// func NewRouter(ctx context.Context, queries *db.Queries, db *sql.DB) []Route {
// 	nationalParkHandler := NationalParkHandler{Ctx: ctx, Queries: queries, DB: db}

// 	var routes = []Route{
// 		route("GET", "/api/national_park", nationalParkHandler.Create),
// 	}

// 	return routes
// }

package httplayer

import (
	"net/http"
)

type Routing interface {
	Routes(with *RoutingDefinition) []Route
}

// A Route consists of an HTTP method, a URL path, and a handler function to be executed when the route is matched.
//
//	  Example:
//
//		 GET /api/v1/users func(http.ResponseWriter, *http.Request)
type Route struct {
	method      []string
	path        string
	handlerFunc http.HandlerFunc
}

// Handler describes a `func(http.ResponseWriter, *http.Request)`.
//
// The final Handler become the result of merged middlewares plus the handler itself.
//
// Example:
//
//	NewBuilder("GET").Path("/").Handler(myHandler).Middleware(mid1, mid2, midN...)
//
// Result execution stack:
//
//	mid1 -> mid2 -> midN -> myHandler
func (a Route) Handler() http.HandlerFunc {
	return a.handlerFunc
}

// Path describe the api path
func (a Route) Path() string {
	return a.path
}

// Methods describes the api method (GET, POST, PUT...)
func (a Route) Methods() []string {
	return a.method
}

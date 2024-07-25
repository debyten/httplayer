package httplayer

import (
	"net/http"
)

// RoutingDefinition represents a slice of Route
type RoutingDefinition struct {
	commonMiddleware []Middleware
	routes           []Route
}

// NewDefinition returns a new *RoutingDefinition instance.
// The provided middleware functions are injected into the routes added after this call.
// The middleware functions are executed sequentially in the order they are provided.
//
//	Usage:
//	 - start with NewDefinition() and pass some middleware e.g. firstMiddleware, secondMiddleware, etc.
//	 - add routes using Add method.
//	 - invoke Build to return a slice of Route.
//
//	Example:
//	NewDefinition(firstMiddleware, secondMiddleware).
//	  Add(http.MethodGet, "/path", handler1).
//	  Add(http.MethodPost, "/path2", handler2, anotherMiddleware)
//
// Resulting stack for routes:
//   - (GET /path): firstMiddleware => secondMiddleware => handler1
//   - (POST /path2): firstMiddleware => secondMiddleware => anotherMiddleware => handler2
func NewDefinition(m ...Middleware) *RoutingDefinition {
	commonMiddleware := make([]Middleware, 0)
	commonMiddleware = append(commonMiddleware, m...)
	return &RoutingDefinition{routes: make([]Route, 0), commonMiddleware: commonMiddleware}
}

// Add saves a Route into *RoutingDefinition instance
//
// Example:
//
//	def := NewDefinition()
//	def.Route("GET", "/api/v1/users", func..., m1, m2, m3)
func (m *RoutingDefinition) Add(method string, path string, h http.HandlerFunc, mid ...Middleware) *RoutingDefinition {
	return m.AddMany([]string{method}, path, h, mid...)
}

// AddMany is like Add but this function can address more than one http method with the same handler.
func (m *RoutingDefinition) AddMany(methods []string, path string, h http.HandlerFunc, mid ...Middleware) *RoutingDefinition {
	builder := NewBuilder(methods...).Path(path).Handler(h)
	allMws := m.commonMiddleware
	if len(mid) > 0 {
		allMws = append(allMws, mid...)
	}
	builder.Middleware(allMws...)
	m.routes = append(m.routes, builder.Build())
	return m
}

// Detach creates a new RoutingDefinition by concatenating the current middlewares from m into
// the newly created RoutingDefinition
func (m *RoutingDefinition) Detach(mid ...Middleware) *RoutingDefinition {
	allMiddleware := make([]Middleware, 0)
	allMiddleware = append(allMiddleware, m.commonMiddleware...)
	allMiddleware = append(allMiddleware, mid...)
	return &RoutingDefinition{
		commonMiddleware: allMiddleware,
		routes:           make([]Route, 0),
	}
}

// Done returns the slice of Route.
func (m *RoutingDefinition) Done() []Route {
	return m.routes
}

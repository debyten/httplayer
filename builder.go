package httplayer

import (
	"net/http"
	"slices"
)

// Builder represents a helper struct to build a single Route
type Builder struct {
	method      []string
	path        string
	handlerFunc http.HandlerFunc
	middleware  []Middleware
}

// NewBuilder start building a Route specifying the method parameters (POST, PUT, PATCH...)
func NewBuilder(method ...string) *Builder {
	return &Builder{method: method}
}

// Path describe the builder path, e.g. /api/v1/test
func (rb *Builder) Path(p string) *Builder {
	rb.path = p
	return rb
}

// Handler is a simple http.HandlerFunc
func (rb *Builder) Handler(h http.HandlerFunc) *Builder {
	rb.handlerFunc = h
	return rb
}

// Middleware sets the middlewares to be injected to the http.HandlerFunc specified with Handler function.
// The slice of input `middleware` will be reversed to respect the sequentiality of the call stack.
//
// Example:
//
//	Middleware(m1, m2, m3) => m1 -> m2 -> m3
//	Middleware(m7, m5, m9) => m7 -> m5 -> m9
func (rb *Builder) Middleware(middleware ...Middleware) *Builder {
	if len(rb.middleware) == 0 {
		rb.middleware = make([]Middleware, 0)
	}
	rb.middleware = append(rb.middleware, middleware...)
	return rb
}

// Build finalize the api build process applying a reverse function on middleware slice (to preserve the order) and
// returns the Route with the handler func concatenated with the middlewares
func (rb *Builder) Build() Route {
	rb.handlerFunc = concat(rb.handlerFunc, rb.middleware)
	return Route{rb.method, rb.path, rb.handlerFunc}
}

func concat(h http.HandlerFunc, middleware []Middleware) http.HandlerFunc {
	slices.Reverse(middleware)
	for _, m := range middleware {
		h = m(h)
	}
	return h
}

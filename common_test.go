package httplayer

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

type ctxTestKey int

const ctxTestValue ctxTestKey = iota

func testResponse(fn func(w http.ResponseWriter, r *http.Request)) (*httptest.ResponseRecorder, error) {
	req, err := http.NewRequest("GET", "/testMiddleware", nil)
	if err != nil {
		return nil, err
	}
	rr := httptest.NewRecorder()
	fn(rr, req)
	return rr, nil
}

// simpleHandler verifies if ctxTestValue is in the request context and writes the result to `w`.
// Replies with 500 InternalServerError otherwise.
func simpleHandler(w http.ResponseWriter, r *http.Request) {
	stack, ok := r.Context().Value(ctxTestValue).([]string)
	if !ok {
		w.WriteHeader(500)
		return
	}
	_ = json.NewEncoder(w).Encode(&stack)
}

func middlewareWithCtxValue(val string) Middleware {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			stack, ok := r.Context().Value(ctxTestValue).([]string)
			if !ok {
				ctx := context.WithValue(r.Context(), ctxTestValue, []string{val})
				h(w, r.WithContext(ctx))
				return
			}
			stack = append(stack, val)
			ctx := context.WithValue(r.Context(), ctxTestValue, stack)
			h(w, r.WithContext(ctx))
		}
	}
}

func buildMiddlewares(values ...string) []Middleware {
	mws := make([]Middleware, 0)
	for _, value := range values {
		mws = append(mws, middlewareWithCtxValue(value))
	}
	return mws
}

func definitionWithCtxValues(values ...string) *RoutingDefinition {
	mws := buildMiddlewares(values...)
	return NewDefinition(mws...)
}

// newTestRoute build a RoutingDefinition with a set of parent middlewares (see definitionWithCtxValues)
// and the simpleHandler.
//
// The provided parentValues will be used to build the parent middlewares.
func newTestRoute(values ...string) Route {
	return definitionWithCtxValues(values...).
		Add("GET", "/testMiddleware", simpleHandler).
		Done()[0]
}

func newTestRouteWithMws(parentValues []string, routeValues []string) Route {
	childMWs := buildMiddlewares(routeValues...)
	return definitionWithCtxValues(parentValues...).
		Add("GET", "/testMiddleware", simpleHandler, childMWs...).
		Done()[0]
}

func newHandlerWithMiddleware(values ...string) Route {
	mws := buildMiddlewares("2", "3")
	return definitionWithCtxValues("0", "1").
		Add("GET", "/testMiddleware", simpleHandler, mws...).Done()[0]
}

// dummy middlewares

var preflightRequestMW Middleware = preflightRequest
var basicAuthMW Middleware = basicAuth

func preflightRequest(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			return
		}
		h(w, r)
	}
}

func basicAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, p, ok := r.BasicAuth()
		if !ok {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		if u != "admin" || p != "admin" {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		h(w, r)
	}
}

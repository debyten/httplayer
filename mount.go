package httplayer

import (
	"context"
	"fmt"
	"net/http"
)

func MountServices(services ...*ServiceBuilder) *http.ServeMux {
	mux := http.NewServeMux()
	for _, service := range services {
		service.MountTo(mux)
	}
	return mux
}

// Mount the provided routing definitions to `http.ServeMux` and returns it.
func Mount(specs ...*RoutingDefinition) *http.ServeMux {
	return mount(specs...)
}

// DynamicMount uses definitions to register the input routes by using the Mount function.
// The channel is used to dynamically register additional routes in a separate
// go routine.
//
// The goroutine will halt when context is done.
func DynamicMount(ctx context.Context, ch <-chan Route, definitions ...*RoutingDefinition) *http.ServeMux {
	mux := mount(definitions...)
	go listenDynamicMounts(ctx, mux, ch)
	return mux
}

func listenDynamicMounts(ctx context.Context, mux *http.ServeMux, ch <-chan Route) {
	for {
		select {
		case <-ctx.Done():
			return
		case a := <-ch:
			mountApi(mux, a)
		}
	}
}

func mount(definitions ...*RoutingDefinition) *http.ServeMux {
	mux := http.NewServeMux()
	for _, def := range definitions {
		for _, route := range def.Done() {
			mountApi(mux, route)
		}
	}
	return mux
}

func mountApi(mux *http.ServeMux, r Route) {
	for _, method := range r.Methods() {
		pattern := fmt.Sprintf("%s %s", method, r.Path())
		mux.HandleFunc(pattern, r.Handler())
	}
}

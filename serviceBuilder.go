package httplayer

import (
	"net/http"
)

func NewServiceBuilder(commonMiddlewares ...Middleware) *ServiceBuilder {
	return &ServiceBuilder{
		routing: make([]Routing, 0),
		mws:     commonMiddlewares,
	}
}

type ServiceBuilder struct {
	routing []Routing
	mws     []Middleware
}

func (s *ServiceBuilder) Add(routing ...Routing) *ServiceBuilder {
	s.routing = append(s.routing, routing...)
	return s
}

// MW appends new middleware to the already known `middlewares` specified in NewServiceBuilder.
func (s *ServiceBuilder) MW(mws ...Middleware) *ServiceBuilder {
	s.mws = append(s.mws, mws...)
	return s
}

func (s *ServiceBuilder) MountTo(mux *http.ServeMux) {
	for _, service := range s.routing {
		definition := NewDefinition(s.mws...)
		for _, a := range service.Routes(definition) {
			MountRoute(mux, a)
		}
	}
}

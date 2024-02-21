package goji

import (
	"github.com/debyten/apibuilder"
	"github.com/debyten/apibuilder/handlers"
	gojiMux "goji.io"
	"goji.io/pat"
)

func NewGojiHandler(m *gojiMux.Mux) handlers.Handler {
	return &gojiHandler{
		m: m,
	}
}

type gojiHandler struct {
	m *gojiMux.Mux
}

func (g *gojiHandler) Handle(apis []apibuilder.API) {
	for _, api := range apis {
		path := pat.NewWithMethods(api.Path(), api.Methods()...)
		g.m.HandleFunc(path, api.Handler())
	}
}

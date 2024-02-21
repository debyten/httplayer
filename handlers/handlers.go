package handlers

import "github.com/debyten/apibuilder"

type Handler interface {
	Handle(apis []httplayer.API)
}

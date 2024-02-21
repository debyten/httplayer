package gin

import (
	"github.com/debyten/apibuilder"
	"github.com/debyten/apibuilder/handlers"
	"github.com/gin-gonic/gin"
)

func NewGinHandler(engine *gin.Engine) handlers.Handler {
	return &ginHandler{
		engine: engine,
	}
}

type ginHandler struct {
	engine *gin.Engine
}

func (g *ginHandler) Handle(apis []apibuilder.API) {
	for _, api := range apis {
		h := gin.WrapF(api.Handler())
		for _, method := range api.Methods() {
			g.engine.Handle(method, api.Path(), h)
		}
	}
}

package routing

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type BaseApi struct {
	route      string
	engine     *gin.Engine
	middleware func(c *gin.Context) bool
}

func InitialiseRoutes(e *gin.Engine, apiFuncs ...func(api BaseApi)) {
	api := BaseApi{
		engine:     e,
		route:      "",
		middleware: func(context *gin.Context) bool { return true },
	}

	for i := range apiFuncs {
		apiFuncs[i](api)
	}
}

func (api BaseApi) Handle(method string, route string, handler func(c *gin.Context)) {
	path := api.route + "/" + route
	var handlerWithMiddleware = func(c *gin.Context) {
		if api.middleware(c) {
			handler(c)
		}
	}
	api.engine.Handle(method, path, handlerWithMiddleware)
}

func (api BaseApi) Get(route string, handler func(c *gin.Context)) {
	api.Handle(http.MethodGet, route, handler)
}

func (api BaseApi) Post(route string, handler func(c *gin.Context)) {
	api.Handle(http.MethodPost, route, handler)
}

func (api BaseApi) Patch(route string, handler func(c *gin.Context)) {
	api.Handle(http.MethodPatch, route, handler)
}

func (api BaseApi) Delete(route string, handler func(c *gin.Context)) {
	api.Handle(http.MethodDelete, route, handler)
}

func (api BaseApi) Group(route string, middleware func(c *gin.Context) bool, group func(api BaseApi)) {
	if route != "" {
		group(BaseApi{
			engine: api.engine,
			route:  api.route + "/" + route,
			middleware: func(c *gin.Context) bool {
				if api.middleware(c) {
					return middleware(c)
				}

				return false
			},
		})
	} else {
		group(BaseApi{
			engine: api.engine,
			route:  api.route,
			middleware: func(c *gin.Context) bool {
				if api.middleware(c) {
					return middleware(c)
				}

				return false
			},
		})
	}
}

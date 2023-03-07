package routing

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitialiseRoutes(e *gin.Engine, apiFuncs ...func(api BaseApi)) {
	api := BaseApi{
		engine:  e,
		route:   "",
		options: ApiOptions{},
	}

	for i := range apiFuncs {
		apiFuncs[i](api)
	}
}

func (api *BaseApi) Handle(method string, route string, handler func(c *gin.Context)) {
	path := api.route + "/" + route
	var handlerWithMiddleware = func(c *gin.Context) {
		if api.runMiddleware(c) {
			handler(c)
		}
	}
	api.engine.Handle(method, path, handlerWithMiddleware)
}

func (api *BaseApi) Get(route string, handler func(c *gin.Context)) {
	api.Handle(http.MethodGet, route, handler)
}

func (api *BaseApi) Post(route string, handler func(c *gin.Context)) {
	api.Handle(http.MethodPost, route, handler)
}

func (api *BaseApi) Patch(route string, handler func(c *gin.Context)) {
	api.Handle(http.MethodPatch, route, handler)
}

func (api *BaseApi) Delete(route string, handler func(c *gin.Context)) {
	api.Handle(http.MethodDelete, route, handler)
}

func (api *BaseApi) Group(route string, options ApiOptions, group func(api BaseApi)) {
	if route != "" {
		group(BaseApi{
			engine:  api.engine,
			route:   api.route + "/" + route,
			options: api.options.mergeOptions(options),
		})
	} else {
		group(BaseApi{
			engine:  api.engine,
			route:   api.route,
			options: api.options.mergeOptions(options),
		})
	}
}

package routing

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
)

var (
	defaultConfig = Config{
		Port:              80,
		DefaultStatusCode: 500,
	}
)

type Response struct {
	// The status code for the response. Defaults to "500 Internal Server Error" unless specified in Setup Config
	Status int
	// The response body.
	Body any
	// See JSONResponse, HTMLResponse, XMLResponse
	Type int
}

// Config provides the required config options to the Setup function.
type Config struct {
	// The port to listen on. If not provided, will default to :80
	Port int
	// Override for Response.Status
	DefaultStatusCode int
}

func Setup(conf ...Config) (*Router, error) {
	config := parseConfig(conf)

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Port))
	if err != nil {
		return nil, err
	}

	router := &Router{
		routes:   make(path),
		listener: ln,
		config:   config,
	}

	return router, nil
}

func parseConfig(conf []Config) Config {
	config := defaultConfig

	providedConfig := Config{}
	if len(conf) > 0 {
		providedConfig = conf[0]
	}

	if providedConfig.Port != 0 {
		config.Port = 0
	}

	if providedConfig.DefaultStatusCode != 0 {
		config.DefaultStatusCode = providedConfig.DefaultStatusCode
	}

	return config
}

// InitialiseRoutes initialises all routes defined within the provided functions.
// Requires a *gin.Engine to be passed. Refer to gin documentation for basic setup.
func InitialiseRoutes(e *gin.Engine, apiFuncs ...func(api BaseApi)) {
	api := BaseApi{
		router:  e,
		route:   "",
		options: ApiOptions{},
	}

	for i := range apiFuncs {
		apiFuncs[i](api)
	}
}

// Handle adds a new api endpoint with the provided method, route, and handler.
func (api *BaseApi) Handle(method string, route string, handler func(c *gin.Context)) {
	path := api.route + "/" + route
	var handlerWithMiddleware = func(c *gin.Context) {
		if api.runMiddleware(c) {
			handler(c)
		}
	}
	http.Handle()
	api.router.Handle(method, path, handlerWithMiddleware)
}

// Get is a shorthand functions for Handle. It has the same functionality as if the user ran `api.Handle(http.MethodGet, ...)`.
// Get adds a new GET endpoint with the provided route and handler.
func (api *BaseApi) Get(route string, handler func(c *gin.Context)) {
	api.Handle(http.MethodGet, route, handler)
}

// Post is a shorthand functions for Handle. It has the same functionality as if the user ran `api.Handle(http.MethodPost, ...)`.
// Post adds a new GET endpoint with the provided route and handler.
func (api *BaseApi) Post(route string, handler func(c *gin.Context)) {
	api.Handle(http.MethodPost, route, handler)
}

// Patch is a shorthand functions for Handle. It has the same functionality as if the user ran `api.Handle(http.MethodPatch, ...)`.
// Patch adds a new GET endpoint with the provided route and handler.
func (api *BaseApi) Patch(route string, handler func(c *gin.Context)) {
	api.Handle(http.MethodPatch, route, handler)
}

// Delete is a shorthand functions for Handle. It has the same functionality as if the user ran `api.Handle(http.MethodDelete, ...)`.
// Delete adds a new GET endpoint with the provided route and handler.
func (api *BaseApi) Delete(route string, handler func(c *gin.Context)) {
	api.Handle(http.MethodDelete, route, handler)
}

// Group defines a new route group with the provided route, and options.
// Use the function within the define sub routes for the Group.
func (api *BaseApi) Group(route string, options ApiOptions, group func(api BaseApi)) {
	if route != "" {
		group(BaseApi{
			router:  api.router,
			route:   api.route + "/" + route,
			options: options.mergeOptions(api.options),
		})
	} else {
		group(BaseApi{
			router:  api.router,
			route:   api.route,
			options: api.options.mergeOptions(options),
		})
	}
}

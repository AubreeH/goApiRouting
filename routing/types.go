package routing

import (
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
)

// Middleware is a function that runs prior to the main function defined for an endpoint.
// Use the provided *gin.Context to run all required checks.
// Return true to run the next function.
// If false is returned, response must be set within the middleware.
type Middleware = func(c *gin.Context) bool

// ApiOptions defines options to use for a certain endpoint group.
type ApiOptions struct {
	Middleware []Middleware
}

// BaseApi is the base struct for all goApiRouting functions.
type BaseApi struct {
	route   string
	router  *Router
	options ApiOptions
}

// NoMiddleware returns an empty ApiOptions struct.
func (_ *BaseApi) NoMiddleware(_ *gin.Context) ApiOptions { return ApiOptions{} }

type Router struct {
	routes   path
	listener net.Listener
	config   Config
}

type path = map[string]func(r *http.Request) (Response, error)

type Context struct {
	request *http.Request
	store   map[string]any
}

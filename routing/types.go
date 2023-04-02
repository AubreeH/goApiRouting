package routing

import (
	"net/http"
)

// Middleware is a function that runs prior to the main function defined for an endpoint.
// Use the provided *Context to run all required checks.
// Return true to run the next function.
// If false is returned, response must be set within the middleware.
type Middleware = func(c *Context, respond func(Response)) bool

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
func (_ *BaseApi) NoMiddleware() ApiOptions {
	return ApiOptions{}
}

type Router struct {
	routes pathMap
	config Config
}

type Context struct {
	Request *http.Request
	writer  http.ResponseWriter
	Store   map[string]interface{}
}

type pathMap = map[string]endpointMap
type endpointMap = map[string]endpointFunc
type endpointFunc = func(*Context, func(Response))

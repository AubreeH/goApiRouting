package routing

import (
	"net/http"
	"regexp"
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
func (*BaseApi) NoMiddleware() ApiOptions {
	return ApiOptions{}
}

func NoMiddleware() ApiOptions {
	return ApiOptions{}
}

type Router struct {
	endpoints endpoints
	routes    pathMap
	config    Config
}

type MethodNotSupportedError error
type NotFoundError error

type Context struct {
	Request *http.Request
	writer  http.ResponseWriter
	Store   map[string]interface{}
}

type pathMap map[string]endpointMap
type endpointMap map[string]endpointFunc
type endpointFunc func(*Context, func(Response))

type endpoints map[string]endpointGroup

type endpointGroup struct {
	endpoints endpoints
	functions endpointMap
	rawRegex  string
	regex     *regexp.Regexp
}

type Response struct {
	// The status code for the response. Defaults to "500 Internal Server Error" unless specified in Setup Config
	Status int
	// The response body.
	Body interface{}
	// See [JSONResponse], [HTMLResponse], [XMLResponse], [FileResponse]
	Type ResponseType
	// The headers to add to the response.
	Headers map[string]string
}

// Config provides the required config options to the Setup function.
type Config struct {
	// The port to listen on. If not provided, will default to :80
	Port int
	// Override for Response.Status
	DefaultStatusCode int
	// OPTION Request Handler
	Option func(*Context) Response
}

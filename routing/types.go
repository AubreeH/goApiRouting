package routing

import (
	"net/http"
	"net/url"
	"regexp"
	"sync"
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
	endpoints *endpointGroup
	routes    pathMap
	config    Config
	mux       *http.ServeMux
	rwMutex   sync.RWMutex
}

type MethodNotSupportedError error
type NotFoundError error
type InternalServerError error

type Context struct {
	// The request object.
	Request *http.Request
	// The response writer object.
	// Using this may lead to unexpected behaviour.
	Writer http.ResponseWriter
	// The store is a map that can be used to store data between middlewares and functions.
	Store *Store
}

type Store struct {
	pathParameters map[string]string
	query          url.Values
	body           []byte
	bodyMap        map[string]interface{}
	store          map[string]interface{}
	mux            sync.RWMutex
}

type pathMap map[string]endpointMap
type endpointMap map[string]endpointFunc
type endpointFunc struct {
	function func(*Context, func(Response)) `json:"-"`
	Path     string
}

type endpoints map[string]*endpointGroup

type endpointGroup struct {
	GroupName        string
	Endpoints        endpoints
	Functions        endpointMap
	RawRegex         string
	CanMatchRawRegex bool
	Regex            *regexp.Regexp
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
	// BaseResponseHeaders are the headers to apply to every response.
	//
	// These values can be overwriten in [Response.Headers]
	//
	// [Response.Headers]: https://pkg.go.dev/github.com/AubreeH/goApiRouting/routing#Response.Headers
	BaseResponseHeaders map[string]string
	// MaxContentLength is the maximum size of the request body in bytes. Set to 0 to disable.
	//
	// Disabled by default
	MaxContentLength uint64
}

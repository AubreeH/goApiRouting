package routing

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
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
	Body interface{}
	// See JSONResponse, HTMLResponse, XMLResponse, FileResponse
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
	// HEAD Request Handler
	Head func(*Context) Response
}

func Setup(conf ...Config) *Router {
	config := parseConfig(conf)

	router := &Router{
		routes: make(pathMap),
		config: config,
	}

	return router
}

func parseConfig(conf []Config) Config {
	config := defaultConfig

	providedConfig := Config{}
	if len(conf) > 0 {
		providedConfig = conf[0]
	}

	if providedConfig.Port != 0 {
		config.Port = providedConfig.Port
	}

	if providedConfig.DefaultStatusCode != 0 {
		config.DefaultStatusCode = providedConfig.DefaultStatusCode
	}

	return config
}

// InitialiseRoutes initialises all routes defined within the provided functions.
// Requires a *gin.Engine to be passed. Refer to gin documentation for basic setup.
func (r *Router) InitialiseRoutes(apiFuncs ...func(api BaseApi)) {
	api := BaseApi{
		router:  r,
		route:   "",
		options: ApiOptions{},
	}

	for i := range apiFuncs {
		apiFuncs[i](api)
	}

	err := r.setupEndpointGroups()
	if err != nil {
		panic(err)
	}

	r.Handle()
}

func (r *Router) Start() error {
	return http.ListenAndServe(fmt.Sprintf(":%v", r.config.Port), nil)
}

func writeResponse(writer http.ResponseWriter, request *http.Request, response Response) {
	var body []byte
	var err error
	switch response.Type {
	case JSONResponse:
		body, err = json.Marshal(response.Body)
		writer.Header().Add("Content-Type", "application/json")
		break
	case HTMLResponse:
		body, err = func() (body []byte, err error) {
			defer func() {
				if e := recover(); e != nil {
					body = []byte{}
					err = errors.New("an error occurred whilst converting body to []byte")
				}
			}()
			return []byte(response.Body.(string)), nil
		}()
		writer.Header().Add("Content-Type", "text/html")
		break
	case XMLResponse:
		body, err = xml.Marshal(response.Body)
		writer.Header().Add("Content-Type", "application/xml")
		break
	case PlainTextResponse:
		body, err = func() (body []byte, err error) {
			defer func() {
				if e := recover(); e != nil {
					body = []byte{}
					err = errors.New("an error occurred whilst converting body to []byte")
				}
			}()
			return []byte(response.Body.(string)), nil
		}()
		writer.Header().Add("Content-Type", "text/plain")
		break
	case FileResponse:
		http.ServeFile(writer, request, response.Body.(string))
		break
	}
	if err != nil {
		writer.WriteHeader(500)
		panic(err)
	}
	for key, value := range response.Headers {
		writer.Header().Add(key, value)
	}
	writer.WriteHeader(response.Status)
	_, err = writer.Write(body)
	if err != nil {
		panic(err)
	}
}

func (r *Router) Handle() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		write := func(response Response) {
			writeResponse(writer, request, response)
		}

		store := make(map[string]interface{})

		fmt.Println(request.Method, request.URL.Path, r.config.Head != nil)

		if request.Method == http.MethodHead {
			if r.config.Head != nil {
				headResponse := r.config.Head(&Context{
					Request: request,
					writer:  writer,
					Store:   store,
				})
				fmt.Println(headResponse)
				write(headResponse)
				return
			}
		}

		handler, err := r.getFunc(request.Method, request.URL.Path, store)
		if err != nil {
			switch err.Error() {
			case "method not supported":
				writeResponse(writer, request, Response{
					Body:   "Method not supported",
					Type:   PlainTextResponse,
					Status: 400,
				})
			case "not found":
				writeResponse(writer, request, Response{
					Body:   "Not Found",
					Type:   PlainTextResponse,
					Status: 404,
				})
			default:
				writeResponse(writer, request, Response{
					Body:   "Unexpected Error",
					Type:   PlainTextResponse,
					Status: 500,
				})
			}

		} else {
			context := Context{
				writer:  writer,
				Request: request,
				Store:   store,
			}
			handler(&context, write)
		}
	})
}

// Handle adds a new api endpoint with the provided method, route, and handler.
func (api *BaseApi) Handle(method string, route string, handler func(*Context) Response) {
	path := api.route + "/" + route
	var handlerWithMiddleware = func(c *Context, respond func(Response)) {
		if api.runMiddleware(c, respond) {
			respond(handler(c))
		}
	}
	if _, ok := api.router.routes[path]; ok {
		api.router.routes[path][method] = handlerWithMiddleware
	} else {
		api.router.routes[path] = make(endpointMap)
		api.router.routes[path][method] = handlerWithMiddleware
	}

}

// Get is a shorthand functions for Handle. It has the same functionality as if the user ran `api.Handle(http.MethodGet, ...)`.
// Get adds a new GET endpoint with the provided route and handler.
func (api *BaseApi) Get(route string, handler func(c *Context) Response) {
	api.Handle(http.MethodGet, route, handler)
}

// Post is a shorthand functions for Handle. It has the same functionality as if the user ran `api.Handle(http.MethodPost, ...)`.
// Post adds a new GET endpoint with the provided route and handler.
func (api *BaseApi) Post(route string, handler func(c *Context) Response) {
	api.Handle(http.MethodPost, route, handler)
}

// Patch is a shorthand functions for Handle. It has the same functionality as if the user ran `api.Handle(http.MethodPatch, ...)`.
// Patch adds a new GET endpoint with the provided route and handler.
func (api *BaseApi) Patch(route string, handler func(c *Context) Response) {
	api.Handle(http.MethodPatch, route, handler)
}

// Delete is a shorthand functions for Handle. It has the same functionality as if the user ran `api.Handle(http.MethodDelete, ...)`.
// Delete adds a new GET endpoint with the provided route and handler.
func (api *BaseApi) Delete(route string, handler func(c *Context) Response) {
	api.Handle(http.MethodDelete, route, handler)
}

func (api *BaseApi) Head(route string, handler func(c *Context) Response) {
	api.Handle(http.MethodHead, route, handler)
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
			options: options.mergeOptions(api.options),
		})
	}
}

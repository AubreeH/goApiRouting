package routing

import (
	"fmt"
	"net/http"
)

var (
	defaultConfig = Config{
		Port:              80,
		DefaultStatusCode: 500,
	}
)

// Setup creates a new router instance with the provided config.
// Only the first config provided will be used.
// If no config is provided, the default config will be used.
func Setup(conf ...Config) *Router {
	config := defaultConfig

	if len(conf) > 0 {
		config = conf[0]
	}

	router := &Router{
		routes: make(pathMap),
		config: config,
	}

	return router
}

// NewRouter creates a new router instance with the provided config.
// Only the first config provided will be used.
// If no config is provided, the default config will be used.
// Alias for Setup.
func NewRouter(conf ...Config) *Router {
	return Setup(conf...)
}

// InitialiseRoutes initialises all routes defined within the provided functions.
func (r *Router) InitialiseRoutes(apiFuncs ...func(api BaseApi)) {
	api := BaseApi{
		router:  r,
		route:   "",
		options: ApiOptions{},
	}

	for i := range apiFuncs {
		apiFuncs[i](api)
	}
	r.setupEndpointGroups()
	r.setupHandler()
}

func (r *Router) Start() error {
	return http.ListenAndServe(fmt.Sprintf(":%v", r.config.Port), r.mux)
}

func (r *Router) AddRoute(path, method string, handler func(*Context, func(Response))) {
	r.rwMutex.Lock()
	defer r.rwMutex.Unlock()

	if r.routes == nil {
		r.routes[path] = make(endpointMap)
	} else if _, ok := r.routes[path]; !ok {
		r.routes[path] = make(endpointMap)
	} else if _, ok := r.routes[path][method]; ok {
		panic(fmt.Sprintf("route %v %v already exists", method, path))
	}

	r.routes[path][method] = endpointFunc{
		function: handler,
		Path:     path,
	}
}

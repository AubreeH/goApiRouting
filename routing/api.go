package routing

import (
	"net/http"
)

// Handle adds a new api endpoint with the provided method, route, and handler.
func (api *BaseApi) Handle(method string, route string, handler func(*Context) Response) {
	path := api.route + "/" + route
	var handlerWithMiddleware = func(c *Context, respond func(Response)) {
		if api.runMiddleware(c, respond) {
			respond(handler(c))
		}
	}
	api.router.AddRoute(path, method, handlerWithMiddleware)
}

// Get is a shorthand function for Handle. It has the same functionality as if the user ran `api.Handle(http.MethodGet, ...)`.
// Get adds a new GET endpoint with the provided route and handler.
func (api *BaseApi) Get(route string, handler func(c *Context) Response) {
	api.Handle(http.MethodGet, route, handler)
}

// Post is a shorthand function for Handle. It has the same functionality as if the user ran `api.Handle(http.MethodPost, ...)`.
// Post adds a new POST endpoint with the provided route and handler.
func (api *BaseApi) Post(route string, handler func(c *Context) Response) {
	api.Handle(http.MethodPost, route, handler)
}

// PUT is a shorthand function for Handle. It has the same functionality as if the user ran `api.Handle(http.MethodPut, ...)`.
// PUT adds a new PUT endpoint with the provided route and handler.
func (api *BaseApi) Put(route string, handler func(c *Context) Response) {
	api.Handle(http.MethodPut, route, handler)
}

// Patch is a shorthand function for Handle. It has the same functionality as if the user ran `api.Handle(http.MethodPatch, ...)`.
// Patch adds a new PATCH endpoint with the provided route and handler.
func (api *BaseApi) Patch(route string, handler func(c *Context) Response) {
	api.Handle(http.MethodPatch, route, handler)
}

// Delete is a shorthand function for Handle. It has the same functionality as if the user ran `api.Handle(http.MethodDelete, ...)`.
// Delete adds a new DELETE endpoint with the provided route and handler.
func (api *BaseApi) Delete(route string, handler func(c *Context) Response) {
	api.Handle(http.MethodDelete, route, handler)
}

// HEAD is a shorthand function for Handle. It has the same functionality as if the user ran `api.Handle(http.MethodHead, ...)`.
// HEAD adds a new HEAD endpoint with the provided route and handler.
func (api *BaseApi) Head(route string, handler func(c *Context) Response) {
	api.Handle(http.MethodHead, route, handler)
}

// OPTIONS is a shorthand function for Handle. It has the same functionality as if the user ran `api.Handle(http.MethodOptions, ...)`.
// OPTIONS adds a new OPTIONS endpoint with the provided route and handler.
func (api *BaseApi) Options(route string, handler func(c *Context) Response) {
	api.Handle(http.MethodOptions, route, handler)
}

// Any is a shorthand function for Handle. It has the same functionality as if the user ran `api.Handle("*", ...)`.
// Any adds a new endpoint with the provided route and handler for all methods.
func (api *BaseApi) Any(route string, handler func(c *Context) Response) {
	api.Handle("*", route, handler)
}

// Group defines a new route group with the provided route, and options.
// Use the function within the define sub routes for the Group.
func (api *BaseApi) Group(route string, options ApiOptions, group func(api BaseApi)) {
	if route == "*" {
		panic("group cannot have wildcard group route")
	}

	if route != "" {
		group(BaseApi{
			router:  api.router,
			route:   api.route + "/" + route,
			options: api.options.mergeOptions(options),
		})
	} else {
		group(BaseApi{
			router:  api.router,
			route:   api.route,
			options: api.options.mergeOptions(options),
		})
	}
}

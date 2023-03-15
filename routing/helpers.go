package routing

import "github.com/gin-gonic/gin"

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
	engine  *gin.Engine
	options ApiOptions
}

// NoMiddleware returns an empty ApiOptions struct.
func (_ *BaseApi) NoMiddleware(_ *gin.Context) ApiOptions { return ApiOptions{} }

// runMiddleware runs all middlewares defined within the ApiOptions instance within the targeted BaseApi.
func (api *BaseApi) runMiddleware(c *gin.Context) bool {
	if api.options.Middleware != nil {
		for _, m := range api.options.Middleware {
			if !m(c) {
				return false
			}
		}
	}

	return true
}

// mergeOptions creates a new ApiOptions instance with all middlewares within both the targeted and supplied ApiOptions instances.
func (options *ApiOptions) mergeOptions(newOptions ApiOptions) ApiOptions {
	options.Middleware = append(options.Middleware, newOptions.Middleware...)

	return *options
}

// WithMiddleware creates a new ApiOptions instance with all provided middlewares.
func WithMiddleware(middleware ...Middleware) ApiOptions {
	return ApiOptions{
		Middleware: middleware,
	}
}

package routing

import "github.com/gin-gonic/gin"

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

package routing

import "github.com/gin-gonic/gin"

type Middleware = func(c *gin.Context) bool

type ApiOptions struct {
	Middleware []Middleware
}

type BaseApi struct {
	route   string
	engine  *gin.Engine
	options ApiOptions
}

func (_ *BaseApi) NoMiddleware(_ *gin.Context) ApiOptions { return ApiOptions{} }

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

func (options *ApiOptions) mergeOptions(newOptions ApiOptions) ApiOptions {
	options.Middleware = append(options.Middleware, newOptions.Middleware...)

	return *options
}

func WithMiddleware(middleware ...Middleware) ApiOptions {
	return ApiOptions{
		Middleware: middleware,
	}
}

package routing

// runMiddleware runs all middlewares defined within the ApiOptions instance within the targeted BaseApi.
func (api *BaseApi) runMiddleware(c *Context, respond func(Response)) bool {
	if api.options.Middleware != nil {
		cont := true
		for _, m := range api.options.Middleware {
			cont = false

			next := func() {
				cont = true
			}

			if !m(c, next, respond) {
				return false
			}

			if !cont {
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

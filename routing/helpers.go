package routing

import (
	"errors"
)

// runMiddleware runs all middlewares defined within the ApiOptions instance within the targeted BaseApi.
func (api *BaseApi) runMiddleware(c *Context, respond func(Response)) bool {
	if api.options.Middleware != nil {
		for _, m := range api.options.Middleware {
			if !m(c, respond) {
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

func (r *Router) getFunc(method, path string) (func(*Context, func(Response)), map[string]string, error) {
	pathParameters := make(map[string]string)

	groups := groupSplitRegex.FindAllStringSubmatch(path, -1)
	currentEndpointGroup := r.endpoints
	var closestWildcardEndpointGroup *endpointGroup

	// prettyPrint(r.endpoints)

	for _, group := range groups {
		if eg, ok := currentEndpointGroup.Endpoints["*"]; ok {
			closestWildcardEndpointGroup = eg
		}

		if eg := currentEndpointGroup.Endpoints.getGroup(group[1]); eg != nil {
			if !eg.CanMatchRawRegex {
				params := eg.extractPathParameters(group[1])
				for k, v := range params {
					pathParameters[k] = v
				}
			}
			currentEndpointGroup = eg
			continue
		}

		currentEndpointGroup = closestWildcardEndpointGroup
		break
	}

	if currentEndpointGroup == nil {
		return nil, nil, errors.New("not found").(NotFoundError)
	}

	if function, ok := currentEndpointGroup.Functions[method]; ok && function.function == nil {
		return nil, nil, errors.New("internal server error").(InternalServerError)
	} else if !ok {
		if function, ok := currentEndpointGroup.Functions["*"]; ok && function.function == nil {
			return nil, nil, errors.New("internal server error").(InternalServerError)
		} else if !ok {
			return nil, nil, errors.New("method not supported").(MethodNotSupportedError)
		} else {
			return function.function, pathParameters, nil
		}
	} else {
		return function.function, pathParameters, nil
	}
}

// WithMiddleware creates a new ApiOptions instance with all provided middlewares.
func WithMiddleware(middleware ...Middleware) ApiOptions {
	return ApiOptions{
		Middleware: middleware,
	}
}

func (context *Context) ServeFile(path string) Response {
	return Response{}
}

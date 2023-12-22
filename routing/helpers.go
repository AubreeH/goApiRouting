package routing

import (
	"errors"
	"strings"
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
	groups := groupPartRegex.Split(strings.Trim(path, "/"), -1)
	if len(groups) == 0 || (len(groups) == 1 && groups[0] == "") {
		if function, err := r.endpoints.getFunc(method); err != nil {
			return nil, nil, err
		} else {
			return function, nil, nil
		}
	}

	pathParameters := make(map[string]string)
	currentEndpointGroup := r.endpoints
	var closestWildcardEndpointGroup *endpointGroup

	for _, group := range groups {
		if eg, ok := currentEndpointGroup.Endpoints["*"]; ok {
			if _, ok := eg.Functions[method]; ok {
				closestWildcardEndpointGroup = eg
			} else if _, ok := eg.Functions["*"]; ok {
				closestWildcardEndpointGroup = eg
			}
		}

		if eg := currentEndpointGroup.Endpoints.getGroup(group); eg != nil {
			if !eg.CanMatchRawRegex {
				params := eg.extractPathParameters(group)
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

	if function, err := currentEndpointGroup.getFunc(method); err != nil {
		return nil, nil, err
	} else {
		return function, pathParameters, nil
	}
}

func (e *endpointGroup) getFunc(method string) (func(*Context, func(Response)), error) {
	if function, ok := e.Functions[method]; ok && function.function == nil {
		return nil, errors.New("internal server error").(InternalServerError)
	} else if ok {
		return function.function, nil
	} else {
	}

	if function, ok := e.Functions["*"]; ok && function.function == nil {
		return nil, errors.New("internal server error").(InternalServerError)
	} else if ok {
		return function.function, nil
	}

	if wildcardEndpointGroup, ok := e.Endpoints["*"]; ok {
		return wildcardEndpointGroup.getFunc(method)
	}

	return nil, errors.New("method not supported").(MethodNotSupportedError)
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

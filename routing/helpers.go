package routing

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	groupSplitRegex      = regexp.MustCompile(`/+([^/]+)`)
	groupConditionsRegex = regexp.MustCompile(`\${(?P<name>[^=]*?)(?:="(?P<condition>.*?)")?}`)
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

func (r *Router) setupEndpointGroups() error {
	if r.endpoints == nil {
		r.endpoints = make(endpoints)
	}

	for route, functions := range r.routes {
		groups := groupSplitRegex.FindAllStringSubmatch(route, -1)
		parents := make([]string, 0)
		for i, group := range groups {
			if len(group) == 2 {
				groupName, regex := getGroupNameAndRegexCondition(group[1])

				endpointGroupMap, err := r.getEndpoints(parents)
				if err != nil {
					return err
				}

				if _, ok := endpointGroupMap[groupName]; !ok {
					var newEndpointGroupFunctions endpointMap
					if i == len(groups)-1 {
						newEndpointGroupFunctions = functions
					} else {
						newEndpointGroupFunctions = make(endpointMap)
					}

					compiledRegex, err := regexp.Compile(regex)
					if err != nil {
						return fmt.Errorf("error whilst compiling regex (%s) for route node with name %s: %v", regex, groupName, err)
					}

					endpointGroupMap[groupName] = endpointGroup{
						endpoints: make(endpoints),
						functions: newEndpointGroupFunctions,
						rawRegex:  regex,
						regex:     compiledRegex,
					}
				}

				parents = append(parents, groupName)
			}
		}
	}

	return nil
}

func getGroupNameAndRegexCondition(group string) (string, string) {
	result := groupConditionsRegex.FindStringSubmatch(group)
	if len(result) == 3 {
		return result[1], result[2]
	} else if len(result) == 2 {
		return result[1], `.*`
	}
	return group, group
}

func (r *Router) getEndpoints(parents []string) (endpoints, error) {
	if parents == nil {
		return nil, errors.New("nil groups array provided in getEndpointGroup")
	}

	var endpointGroupMap = r.endpoints

	for i, groupName := range parents {
		if group, ok := endpointGroupMap[groupName]; ok {
			endpointGroupMap = group.endpoints
		} else {
			return nil, fmt.Errorf("endpoint group with name %s does not exist in tree `/%s`", groupName, strings.Join(parents[:i], "/"))
		}
	}

	return endpointGroupMap, nil
}

func (r *Router) getEndpointGroup(path string, store map[string]interface{}) (endpointGroup, error) {
	endpointGroupMap := r.endpoints

	groups := groupSplitRegex.FindAllStringSubmatch(path, -1)

	for i, group := range groups {
		if len(group) == 2 {
			if matchedEndpointGroup, err := endpointGroupMap.getGroup(group[1], store); err == nil {
				if len(groups)-1 == i {
					return matchedEndpointGroup, nil
				} else {
					endpointGroupMap = matchedEndpointGroup.endpoints
				}
			} else {
				return endpointGroup{}, err
			}
		}
	}

	return endpointGroup{}, errors.New("not found").(NotFoundError)
}

func (r *Router) getFunc(method, path string, store map[string]interface{}) (func(*Context, func(Response)), error) {
	group, err := r.getEndpointGroup(path, store)
	if err != nil {
		return nil, err
	}

	if function, ok := group.functions[method]; ok && function != nil {
		return function, nil
	} else {
		return nil, errors.New("not found").(NotFoundError)
	}
}

func (e endpoints) getGroup(value string, store map[string]interface{}) (endpointGroup, error) {
	var regexEndpoints endpoints

	for groupName, group := range e {
		if group.rawRegex != groupName {
			if regexEndpoints == nil {
				regexEndpoints = make(endpoints)
			}

			regexEndpoints[groupName] = group
		} else {
			if groupName == value {
				return group, nil
			}
		}
	}

	for groupName, group := range regexEndpoints {
		if group.regex.MatchString(value) {
			store[groupName] = value
			return group, nil
		}
	}

	return endpointGroup{}, errors.New("not found")
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

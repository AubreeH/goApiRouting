package routing

// Initialises the endpoint groups based on the specified routes.
func (r *Router) setupEndpointGroups() {
	r.endpoints = &endpointGroup{
		Endpoints: make(endpoints),
		GroupName: "",
	}
	for route, m := range r.routes {
		// fmt.Printf("route: '%s'\n", route)
		// fmt.Printf("endpoints before: %+v\n", r.endpoints)
		r.endpoints.addRoute(route, m)
		// fmt.Printf("endpoints after: %+v\n", r.endpoints)

	}
}

// Adds a route to the endpoint group map.
func (eg *endpointGroup) addRoute(route string, m endpointMap) {
	groups := groupSplitRegex.FindAllStringSubmatch(route, -1)

	if len(groups) == 0 {
		eg.addFunctions(m)
		return
	}

	currentEndpointGroup := eg
	for i, group := range groups {
		if i == len(groups)-1 {
			currentEndpointGroup.addGroup(group[1]).addFunctions(m)
		} else {
			currentEndpointGroup = currentEndpointGroup.addGroup(group[1])
		}
	}
}

// Creates a new endpoint group if one doesn't exist. Returns the existing/created endpoint group.
func (eg *endpointGroup) addGroup(groupName string) *endpointGroup {
	if g, ok := eg.Endpoints[groupName]; ok {
		if g.Endpoints == nil {
			g.Endpoints = make(endpoints)
		}

		return g
	}

	return eg.newEndpointGroup(groupName).setupRegex()
}

// Creates a new endpoint group.
func (eg *endpointGroup) newEndpointGroup(groupName string) *endpointGroup {
	newEndpointGroup := &endpointGroup{
		Endpoints: make(endpoints),
		GroupName: groupName,
	}
	eg.Endpoints[groupName] = newEndpointGroup
	return newEndpointGroup
}

func (e endpoints) getGroup(value string) *endpointGroup {
	for _, group := range e {
		if group.GroupName == "*" {
			continue
		}
		if group.CanMatchRawRegex && group.RawRegex == value {
			return group
		}
	}

	for _, group := range e {
		if group.GroupName == "*" {
			continue
		}
		if group.Regex.MatchString(value) {
			return group
		}
	}

	return nil
}

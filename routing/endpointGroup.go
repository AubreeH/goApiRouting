package routing

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	groupSplitRegex      = regexp.MustCompile(`/+([^/]+)`)
	groupConditionsRegex = regexp.MustCompile(`(?:\${(?P<name>[^=]*?)(?:="(?P<condition>.*?)")?})`)
)

func (g *endpointGroup) addFunctions(handlers endpointMap) {
	for method, handler := range handlers {
		if g.Functions == nil {
			g.Functions = make(endpointMap)
		}

		if _, ok := g.Functions[method]; ok {
			panic(fmt.Errorf("duplicate method %s for endpoint group %s", method, g.GroupName))
		}

		g.Functions[method] = handler
	}
}

func (g *endpointGroup) setupRegex() *endpointGroup {
	var formattedRegex = g.GroupName
	if conditions := groupConditionsRegex.FindAllStringSubmatch(formattedRegex, -1); conditions != nil {
		for _, condition := range conditions {
			g.CanMatchRawRegex = false
			str := condition[0]
			conditionName := condition[1]
			conditionValue := condition[2]
			if conditionValue == "" {
				conditionValue = ".+?"
			}
			formattedRegex = strings.Replace(formattedRegex, str, fmt.Sprintf("(?P<%s>%s)", conditionName, conditionValue), 1)
		}
	}
	formattedRegex = "^" + formattedRegex + "$"

	g.RawRegex = formattedRegex
	g.Regex = regexp.MustCompile(formattedRegex)

	return g
}

func (eg *endpointGroup) extractPathParameters(value string) map[string]string {
	matches := eg.Regex.FindStringSubmatch(value)
	if matches == nil {
		return nil
	} else if len(matches) == 0 {
		return make(map[string]string)
	}

	params := make(map[string]string)
	for i, name := range eg.Regex.SubexpNames() {
		if i != 0 && name != "" {
			params[name] = matches[i]
		}
	}

	return params
}

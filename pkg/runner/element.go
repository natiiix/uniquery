package runner

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/natiiix/uniquery/pkg/filters"
	"github.com/natiiix/uniquery/pkg/parser"
)

type Element struct {
	Value  interface{}
	Parent *Element
}

func (e Element) GetChildren() []Element {
	switch t := e.Value.(type) {
	case map[string]interface{}:
		children := []Element{}
		for _, v := range t {
			children = append(children, NewElement(v, &e))
		}
		return children

	case map[interface{}]interface{}:
		children := []Element{}
		for _, v := range t {
			children = append(children, NewElement(v, &e))
		}
		return children

	case []interface{}:
		children := []Element{}
		for _, v := range t {
			children = append(children, NewElement(v, &e))
		}
		return children

	default:
		return []Element{}
	}
}

func (e Element) GetChildrenRecursive() []Element {
	// NOTE: Includes the element itself

	children := []Element{e}

	for _, c := range e.GetChildren() {
		children = append(children, c.GetChildrenRecursive()...)
	}

	return children
}

func (e Element) MatchesFilters(valueFilters []filters.Filter) bool {
	for _, f := range valueFilters {
		if !f.IsMatch(e.Value) {
			return false
		}
	}

	return true
}

func NewElement(value interface{}, parent *Element) Element {
	return Element{
		Value:  value,
		Parent: parent,
	}
}

func NewElementRoot(value interface{}) Element {
	return NewElement(value, nil)
}

func compareKey(key interface{}, specifier string) bool {
	switch t := key.(type) {
	case string:
		return t == specifier

	case bool:
		specLower := strings.ToLower(specifier)
		if t {
			return specLower == "true" || specLower == "on" || specLower == "enabled" || specLower == "enable"
		} else {
			return specLower == "false" || specLower == "off" || specLower == "disabled" || specLower == "disable"
		}

	default:
		return fmt.Sprintf("%v", t) == specifier
	}
}

func (e Element) Query(parts []parser.QueryPart) []Element {
	if len(parts) == 0 {
		return []Element{e}
	}

	part := parts[0]
	subquery := parts[1:]

	selected := []Element{}
	spec := part.Specifier

	if spec == "" && e.Parent != nil {
		selected = []Element{*e.Parent}
	} else if spec == "*" {
		selected = e.GetChildren()
	} else if spec == "**" {
		selected = e.GetChildrenRecursive()
	} else {
		switch t := e.Value.(type) {
		case map[string]interface{}:
			if child, exists := t[spec]; exists {
				selected = []Element{NewElement(child, &e)}
			}

		case map[interface{}]interface{}:
			for k, v := range t {
				if compareKey(k, spec) {
					selected = append(selected, NewElement(v, &e))
				}
			}

		case []interface{}:
			if index, err := strconv.Atoi(spec); err == nil && (index >= 0 && index < len(t)) {
				selected = []Element{NewElement(t[index], &e)}
			}
		}
	}

	results := []Element{}
	for _, e := range selected {
		if e.MatchesFilters(part.Filters) {
			results = append(results, e.Query(subquery)...)
		}
	}
	return results
}

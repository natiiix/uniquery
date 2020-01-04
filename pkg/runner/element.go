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
	Key    interface{}
}

func (e Element) GetChildren() map[string]Element {
	switch t := e.Value.(type) {
	case map[string]interface{}:
		children := map[string]Element{}
		for k, v := range t {
			elem := NewElement(v, &e, k)
			children[elem.GetFullPath()] = elem
		}
		return children

	case map[interface{}]interface{}:
		children := map[string]Element{}
		for k, v := range t {
			elem := NewElement(v, &e, k)
			children[elem.GetFullPath()] = elem
		}
		return children

	case []interface{}:
		children := map[string]Element{}
		for k, v := range t {
			elem := NewElement(v, &e, k)
			children[elem.GetFullPath()] = elem
		}
		return children

	default:
		return map[string]Element{}
	}
}

func (e Element) GetChildrenRecursive() map[string]Element {
	// NOTE: Includes the element itself

	children := e.ToMap()

	for k1, v1 := range e.GetChildren() {
		children[k1] = v1

		for k2, v2 := range v1.GetChildrenRecursive() {
			children[k2] = v2
		}
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

func (e Element) GetFullPath() string {
	keyStr := fmt.Sprintf("%#v", e.Key)

	if e.Parent == nil || (e.Parent.Parent == nil && e.Parent.Key == nil) {
		return keyStr
	} else {
		return fmt.Sprintf("%s.%s", e.Parent.GetFullPath(), keyStr)
	}
}

func (e Element) ToMap() map[string]Element {
	return map[string]Element{e.GetFullPath(): e}
}

func NewElement(value interface{}, parent *Element, key interface{}) Element {
	return Element{
		Value:  value,
		Parent: parent,
		Key:    key,
	}
}

func NewElementRoot(value interface{}) Element {
	return NewElement(value, nil, nil)
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

func (e Element) Query(parts []parser.QueryPart) map[string]Element {
	if len(parts) == 0 {
		return e.ToMap()
	}

	part := parts[0]
	subquery := parts[1:]

	selected := map[string]Element{}
	spec := part.Specifier

	if spec == "" && e.Parent != nil {
		selected = e.Parent.ToMap()
	} else if spec == "*" {
		selected = e.GetChildren()
	} else if spec == "**" {
		selected = e.GetChildrenRecursive()
	} else {
		switch t := e.Value.(type) {
		case map[string]interface{}:
			if child, exists := t[spec]; exists {
				selected = NewElement(child, &e, spec).ToMap()
			}

		case map[interface{}]interface{}:
			for k, v := range t {
				if compareKey(k, spec) {
					elem := NewElement(v, &e, k)
					selected[elem.GetFullPath()] = elem
				}
			}

		case []interface{}:
			if index, err := strconv.Atoi(spec); err == nil && (index >= 0 && index < len(t)) {
				selected = NewElement(t[index], &e, strconv.Itoa(index)).ToMap()
			}
		}
	}

	results := map[string]Element{}
	for _, e := range selected {
		if e.MatchesFilters(part.Filters) {
			for k, v := range e.Query(subquery) {
				results[k] = v
			}
		}
	}
	return results
}

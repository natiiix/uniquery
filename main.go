package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	jsonFile  = "test.json"
	testQuery = `**.href!~".* item$"..title`
)

func must(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

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

func (e Element) MatchesFilters(filters []Filter) bool {
	for _, f := range filters {
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

type EqualityFilter struct {
	Value string
}

func (f EqualityFilter) IsMatch(value interface{}) bool {
	if valueStr, ok := value.(string); ok {
		return valueStr == f.Value
	}

	return false
}

type RegexFilter struct {
	Regex *regexp.Regexp
}

func (f RegexFilter) IsMatch(value interface{}) bool {
	if valueStr, ok := value.(string); ok {
		return f.Regex.MatchString(valueStr)
	}

	return false
}

type InvertFilter struct {
	InnerFilter Filter
}

func (f InvertFilter) IsMatch(value interface{}) bool {
	return !f.InnerFilter.IsMatch(value)
}

type Filter interface {
	IsMatch(value interface{}) bool
}

func (e Element) Query(parts []QueryPart) []Element {
	if len(parts) == 0 {
		return []Element{e}
	}

	part := parts[0]
	subquery := parts[1:]

	selected := []Element{}
	spec := part.Specifier

	if spec == "" {
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

		case []interface{}:
			if index, err := strconv.Atoi(spec); err == nil && (index >= 0 || index < len(t)) {
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

type QueryPart struct {
	Specifier string
	Filters   []Filter
}

const (
	escapeRune    = '\\'
	quoteRune     = '"'
	specifierRune = '.'
	equalityRune  = '='
	regexRune     = '~'
	invertRune    = '!'
)

const (
	filterEquality = iota
	filterRegex
)

func ParseSinglePart(query []rune) (string, int) {
	sb := strings.Builder{}
	escaped := false
	quoted := false

	for i, r := range query {
		if escaped {
			sb.WriteRune(r)
			escaped = false
		} else if quoted {
			if r == quoteRune {
				quoted = false
			} else {
				sb.WriteRune(r)
			}
		} else {
			switch r {
			case specifierRune, equalityRune, regexRune, invertRune:
				return sb.String(), i

			case escapeRune:
				escaped = true

			case quoteRune:
				quoted = true

			default:
				sb.WriteRune(r)
			}
		}
	}

	if escaped || quoted {
		log.Fatalln("Unexpected end of query - trailing escape or quote")
	}

	return sb.String(), len(query)
}

func ParseSingleFilter(query []rune) (Filter, int) {
	if len(query) <= 0 {
		log.Fatalln("Unexpected end of filter")
		return nil, 0
	}

	switch query[0] {
	case specifierRune:
		return nil, 0

	case equalityRune:
		value, len := ParseSinglePart(query[1:])
		return EqualityFilter{Value: value}, 1 + len

	case regexRune:
		regex, len := ParseSinglePart(query[1:])
		return RegexFilter{Regex: regexp.MustCompile(regex)}, 1 + len

	case invertRune:
		inner, len := ParseSingleFilter(query[1:])
		return InvertFilter{InnerFilter: inner}, 1 + len

	default:
		log.Fatalln("Unexpected filter prefix:", string(query[0]))
		return nil, 0
	}
}

func ParseQuery(query string) []QueryPart {
	queryRunes := []rune(query)
	parts := []QueryPart{}

	for i := 0; i < len(queryRunes); i++ {
		specifier, specifierLength := ParseSinglePart(queryRunes[i:])
		i += specifierLength

		filters := []Filter{}

		for i < len(queryRunes) {
			if filter, filterLength := ParseSingleFilter(queryRunes[i:]); filter != nil {
				filters = append(filters, filter)
				i += filterLength
			} else {
				break
			}
		}

		parts = append(parts, QueryPart{Specifier: specifier, Filters: filters})
	}

	return parts
}

func main() {
	f, err := os.Open(jsonFile)
	must(err)
	defer f.Close()

	var root interface{}
	err = json.NewDecoder(f).Decode(&root)
	must(err)

	queryParts := ParseQuery(testQuery)
	fmt.Printf("%+q\n", queryParts)

	rootElem := NewElementRoot(root)
	for i, v := range rootElem.Query(queryParts) {
		fmt.Printf("%d: %#v\n", i, v.Value)
	}
}

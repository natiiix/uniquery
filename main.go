package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	jsonFile  = "test.json"
	testQuery = `p..div.*.title="ZXCV"..`
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
	} else {
		switch t := e.Value.(type) {
		case map[string]interface{}:
			if spec == "*" {
				selected = []Element{}
				for _, v := range t {
					selected = append(selected, NewElement(v, &e))
				}
			} else if child, exists := t[spec]; exists {
				selected = []Element{NewElement(child, &e)}
			}

		case []interface{}:
			if spec == "*" {
				selected = []Element{}
				for _, v := range t {
					selected = append(selected, NewElement(v, &e))
				}
			} else if index, err := strconv.Atoi(spec); err == nil && index >= 0 || index < len(t) {
				selected = []Element{NewElement(t[index], &e)}
			} else {
				log.Fatalln("Invalid index:", spec)
			}

		default:
			log.Fatalln("Unexpected JSON type:", t)
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
)

const (
	filterEquality = iota
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
			case specifierRune, equalityRune:
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

func ParseQuery(query string) []QueryPart {
	queryRunes := []rune(query)
	parts := []QueryPart{}

	for i := 0; i < len(queryRunes); i++ {
		specifier, specifierLength := ParseSinglePart(queryRunes[i:])
		i += specifierLength

		filters := []Filter{}

		for i < len(queryRunes) && queryRunes[i] != specifierRune {
			filterValue, filterLength := ParseSinglePart(queryRunes[i+1:])

			switch queryRunes[i] {
			case equalityRune:
				filters = append(filters, EqualityFilter{Value: filterValue})

			default:
				log.Fatalln("Unexpected meta rune:", queryRunes[i])
			}

			i += 1 + filterLength
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

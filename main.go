package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
)

const (
	jsonFile  = "test.json"
	testQuery = "p..div.*.title=ZXCV..*"
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

func FilterMatches(f Filter, values []interface{}) []interface{} {
	filtered := []interface{}{}

	for _, val := range values {
		if f.IsMatch(val) {
			filtered = append(filtered, val)
		}
	}

	return filtered
}

func (e Element) Query(parts []QueryPart) []Element {
	if len(parts) == 0 {
		return []Element{e}
	}

	part := parts[0]
	subquery := parts[1:]

	spec := part.Specifier

	if spec == "" {
		return e.Parent.Query(subquery)
	}

	switch t := e.Value.(type) {
	case map[string]interface{}:
		if spec == "*" {
			elems := []Element{}

			for _, v := range t {
				elems = append(elems, NewElement(v, &e).Query(subquery)...)
			}

			return elems
		} else if child, exists := t[spec]; exists {
			return NewElement(child, &e).Query(subquery)
		}

	case []interface{}:
		if spec == "*" {
			elems := []Element{}

			for _, v := range t {
				elems = append(elems, NewElement(v, &e).Query(subquery)...)
			}

			return elems
		} else if index, err := strconv.Atoi(spec); err == nil && index >= 0 || index < len(t) {
			return NewElement(t[index], &e).Query(subquery)
		} else {
			log.Fatalln("Invalid index:", spec)
		}

	default:
		log.Fatalln("Unexpected JSON type:", t)
	}

	return []Element{}
}

type QueryPart struct {
	Specifier string
	Filters   []Filter
}

const (
	specifierRune = '.'
	equalityRune  = '='
)

const (
	filterEquality = iota
)

func CreateFilter(valueRunes []rune, filterType int) Filter {
	valueString := string(valueRunes)

	switch filterType {
	case filterEquality:
		return EqualityFilter{Value: valueString}

	default:
		log.Fatalln("Unexpected filter type:", filterType)
		return nil
	}
}

func ParseSinglePart(query []rune) (string, int) {
	// TODO: Add espacing via backslash and quotes
	for i, r := range query {
		if r == specifierRune || r == equalityRune {
			return string(query[:i]), i
		}
	}

	return string(query), len(query)
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

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
	testQuery = "p..div.*.desc..*"
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

func NewElement(value interface{}, parent *Element) *Element {
	return &Element{
		Value:  value,
		Parent: parent,
	}
}

func (e *Element) Query(parts []string) []*Element {
	if len(parts) == 0 {
		return []*Element{e}
	}

	specifier := parts[0]
	subquery := parts[1:]

	if specifier == "" {
		return e.Parent.Query(subquery)
	}

	switch t := e.Value.(type) {
	case map[string]interface{}:
		if specifier == "*" {
			elems := []*Element{}

			for _, v := range t {
				elems = append(elems, NewElement(v, e).Query(subquery)...)
			}

			return elems
		} else if child, exists := t[specifier]; exists {
			return NewElement(child, e).Query(subquery)
		}

	case []interface{}:
		if specifier == "*" {
			elems := []*Element{}

			for _, v := range t {
				elems = append(elems, NewElement(v, e).Query(subquery)...)
			}

			return elems
		} else if index, err := strconv.Atoi(specifier); err == nil && index >= 0 || index < len(t) {
			return NewElement(t[index], e).Query(subquery)
		} else {
			log.Fatalln("Invalid index:", specifier)
		}

	default:
		log.Fatalln("Unexpected JSON type:", t)
	}

	return []*Element{}
}

func ParseQuery(query string) []string {
	queryRunes := []rune(query)
	parts := []string{}

	specifierBuilder := strings.Builder{}
	escape := false

	for i := 0; i < len(queryRunes); i++ {
		r := queryRunes[i]

		if escape {
			specifierBuilder.WriteRune(r)
			escape = false
		} else if r == '\\' {
			escape = true
		} else if r == '.' {
			parts = append(parts, specifierBuilder.String())
			specifierBuilder.Reset()
		} else {
			specifierBuilder.WriteRune(r)
		}
	}

	parts = append(parts, specifierBuilder.String())
	return parts
}

func main() {
	f, err := os.Open(jsonFile)
	must(err)
	defer f.Close()

	var root interface{}
	err = json.NewDecoder(f).Decode(&root)
	must(err)

	rootElem := NewElement(root, nil)
	queryParts := ParseQuery(testQuery)

	fmt.Printf("%#v\n", queryParts)

	for i, v := range rootElem.Query(queryParts) {
		fmt.Printf("%d: %#v\n", i, v.Value)
	}
}

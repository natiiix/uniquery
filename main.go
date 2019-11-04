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
	testQuery = "p..div.1.href"
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
	// TODO: Implement `*` wildcard specifier support.
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
		if child, exists := t[specifier]; exists {
			return NewElement(child, e).Query(subquery)
		}

	case []interface{}:
		index, err := strconv.Atoi(specifier)
		if err != nil || index < 0 || index >= len(t) {
			log.Fatalln("Invalid index:", specifier)
		}

		return NewElement(t[index], e).Query(subquery)

	default:
		log.Fatalln("Unexpected JSON type:", t)
	}

	return []*Element{}
}

func main() {
	f, err := os.Open(jsonFile)
	must(err)
	defer f.Close()

	var root interface{}
	err = json.NewDecoder(f).Decode(&root)
	must(err)

	rootElem := NewElement(root, nil)
	queryParts := strings.Split(testQuery, ".")

	for i, v := range rootElem.Query(queryParts) {
		fmt.Printf("%d: %#v\n", i, v)
	}
}

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

func NewElement(value interface{}, parent *Element) *Element {
	return &Element{
		Value:  value,
		Parent: parent,
	}
}

func (e *Element) Query(parts []*QueryPart) []*Element {
	if len(parts) == 0 {
		return []*Element{e}
	}

	part := parts[0]
	subquery := parts[1:]

	if part.Specifier != nil {
		spec := *part.Specifier

		if spec == "" {
			return e.Parent.Query(subquery)
		}

		switch t := e.Value.(type) {
		case map[string]interface{}:
			if spec == "*" {
				elems := []*Element{}

				for _, v := range t {
					elems = append(elems, NewElement(v, e).Query(subquery)...)
				}

				return elems
			} else if child, exists := t[spec]; exists {
				return NewElement(child, e).Query(subquery)
			}

		case []interface{}:
			if spec == "*" {
				elems := []*Element{}

				for _, v := range t {
					elems = append(elems, NewElement(v, e).Query(subquery)...)
				}

				return elems
			} else if index, err := strconv.Atoi(spec); err == nil && index >= 0 || index < len(t) {
				return NewElement(t[index], e).Query(subquery)
			} else {
				log.Fatalln("Invalid index:", spec)
			}

		default:
			log.Fatalln("Unexpected JSON type:", t)
		}
	}

	return []*Element{}
}

type QueryPart struct {
	Specifier *string
	Equality  *string
	Regex     *string
}

func StringOrNil(strPtr *string) string {
	if strPtr == nil {
		return fmt.Sprint(strPtr)
	}

	return *strPtr
}

func (qp *QueryPart) String() string {
	return fmt.Sprintf("%s=%s~%s",
		StringOrNil(qp.Specifier),
		StringOrNil(qp.Equality),
		StringOrNil(qp.Regex))
}

func NewQueryPart(specifier, equality *string) *QueryPart {
	return &QueryPart{
		Specifier: specifier,
		Equality:  equality,
		Regex:     nil,
	}
}

func ParseQuery(query string) []*QueryPart {
	queryRunes := []rune(query)
	parts := []*QueryPart{}

	sb := strings.Builder{}
	escape := false

	var specifier *string = nil
	var equality *string = nil

	equalityPresent := false

	for i := 0; i < len(queryRunes); i++ {
		r := queryRunes[i]

		if escape {
			sb.WriteRune(r)
			escape = false
		} else if r == '\\' {
			escape = true
		} else if r == '.' {
			str := sb.String()
			sb.Reset()

			if equalityPresent {
				equality = &str
				equalityPresent = false
			} else {
				specifier = &str
			}

			parts = append(parts, NewQueryPart(specifier, equality))
		} else if r == '=' {
			equalityPresent = true
			str := sb.String()
			sb.Reset()
			specifier = &str

			// TODO: Make sure it wasn't present before
		} else {
			sb.WriteRune(r)
		}
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

	rootElem := NewElement(root, nil)
	queryParts := ParseQuery(testQuery)

	fmt.Printf("%+q\n", queryParts)

	for i, v := range rootElem.Query(queryParts) {
		fmt.Printf("%d: %#v\n", i, v.Value)
	}
}

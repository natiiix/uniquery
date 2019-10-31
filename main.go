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
	testQuery = "div..href"
)

func must(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func find(elem interface{}, queryParts []string) []interface{} {
	if len(queryParts) == 0 {
		return []interface{}{elem}
	}

	specifier := queryParts[0]
	subquery := queryParts[1:]

	switch t := elem.(type) {
	case map[string]interface{}:
		if child, exists := t[specifier]; exists {
			return find(child, subquery)
		}

	case []interface{}:
		if specifier == "" {
			results := []interface{}{}

			for _, item := range t {
				results = append(results, find(item, subquery)...)
			}

			return results
		}

		index, err := strconv.Atoi(specifier)
		if err != nil || index < 0 || index >= len(t) {
			log.Fatalln("Invalid index:", specifier)
		}

		return find(t[index], subquery)

	default:
		log.Fatalln("Unexpected JSON type:", t)
	}

	return []interface{}{}
}

func main() {
	f, err := os.Open(jsonFile)
	must(err)
	defer f.Close()

	var root interface{}
	err = json.NewDecoder(f).Decode(&root)
	must(err)

	queryParts := strings.Split(testQuery, ".")

	fmt.Printf("%#v\n", find(root, queryParts))
}

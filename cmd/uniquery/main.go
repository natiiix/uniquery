package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/natiiix/uniquery/pkg/parser"
	"github.com/natiiix/uniquery/pkg/runner"
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

func main() {
	f, err := os.Open(jsonFile)
	must(err)
	defer f.Close()

	var root interface{}
	err = json.NewDecoder(f).Decode(&root)
	must(err)

	queryParts := parser.ParseQuery(testQuery)
	fmt.Printf("%+q\n", queryParts)

	rootElem := runner.NewElementRoot(root)
	for i, v := range rootElem.Query(queryParts) {
		fmt.Printf("%d: %#v\n", i, v.Value)
	}
}

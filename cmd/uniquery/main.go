package main

import (
	"log"

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
	if results, err := runner.RunJsonFile(testQuery, jsonFile); err != nil {
		log.Fatalln(err)
	} else {
		for i, v := range results {
			log.Printf("%d: %#v\n", i, v.Value)
		}
	}
}

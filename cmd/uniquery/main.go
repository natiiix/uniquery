package main

import (
	"log"
	"flag"

	"github.com/natiiix/uniquery/pkg/runner"
)

var (
	query string = ""
	jsonPath string = ""
	verbose bool = false
)

func init() {
	flag.StringVar(&query, "query", query, "Query to run on the data")
	flag.StringVar(&jsonPath, "json", jsonPath, "Path of a JSON file to run the query on")
	flag.BoolVar(&verbose, "v", verbose, "Enable verbose mode - additional information will be printed, mostly for debugging purposes")
	flag.Parse()

	if jsonPath == "" {
		log.Fatalln("Please specify a JSON file path")
	}

	runner.Verbose = verbose
}

func main() {
	if results, err := runner.RunJsonFile(query, jsonPath); err != nil {
		log.Fatalln(err)
	} else {
		for i, v := range results {
			log.Printf("%d: %#v\n", i, v.Value)
		}
	}
}

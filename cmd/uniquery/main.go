package main

import (
	"flag"
	"log"

	"github.com/natiiix/uniquery/pkg/runner"
)

var (
	query    string = ""
	jsonPath string = ""
	yamlPath string = ""
	verbose  bool   = false
)

func init() {
	flag.StringVar(&query, "query", query, "Query to run on the data")
	flag.StringVar(&jsonPath, "json", jsonPath, "Path of a JSON file to run the query on")
	flag.StringVar(&yamlPath, "yaml", jsonPath, "Path of a YAML file to run the query on")
	flag.BoolVar(&verbose, "v", verbose, "Enable verbose mode - additional information will be printed, mostly for debugging purposes")
	flag.Parse()

	if jsonPath == "" && yamlPath == "" {
		log.Fatalln("Please specify an input file path")
	}

	runner.Verbose = verbose
}

func main() {
	results := map[string]runner.Element{}

	if jsonPath != "" {
		if jsonResults, err := runner.RunJsonFile(query, jsonPath); err != nil {
			log.Fatalln(err)
		} else {
			for k, v := range jsonResults {
				results[k] = v
			}
		}
	}

	if yamlPath != "" {
		if yamlResults, err := runner.RunYamlFile(query, yamlPath); err != nil {
			log.Fatalln(err)
		} else {
			for k, v := range yamlResults {
				results[k] = v
			}
		}
	}

	for k, v := range results {
		log.Printf("[%v] -- %#v\n", k, v.Value)
	}
}

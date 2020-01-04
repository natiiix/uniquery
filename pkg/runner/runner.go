package runner

import (
	"encoding/json"
	"log"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/natiiix/uniquery/pkg/parser"
)

var Verbose bool = false

func Run(query string, root interface{}) map[string]Element {
	queryParts := parser.ParseQuery(query)
	if Verbose {
		log.Printf("Parsed query: %+q\n", queryParts)
	}
	rootElem := NewElementRoot(root)
	results := rootElem.Query(queryParts)
	return results
}

func RunJson(query string, jsonData []byte) (map[string]Element, error) {
	var root interface{}
	err := json.Unmarshal(jsonData, &root)
	if err != nil {
		return nil, err
	}

	return Run(query, root), nil
}

func RunJsonString(query string, jsonStr string) (map[string]Element, error) {
	return RunJson(query, []byte(jsonStr))
}

func RunJsonFile(query string, jsonPath string) (map[string]Element, error) {
	f, err := os.Open(jsonPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var root interface{}
	if err = json.NewDecoder(f).Decode(&root); err != nil {
		return nil, err
	}

	return Run(query, root), nil
}

func RunYaml(query string, yamlData []byte) (map[string]Element, error) {
	var root interface{}
	err := yaml.Unmarshal(yamlData, &root)
	if err != nil {
		return nil, err
	}

	return Run(query, root), nil
}

func RunYamlString(query string, yamlStr string) (map[string]Element, error) {
	return RunYaml(query, []byte(yamlStr))
}

func RunYamlFile(query string, jsonPath string) (map[string]Element, error) {
	f, err := os.Open(jsonPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var root interface{}
	if err = yaml.NewDecoder(f).Decode(&root); err != nil {
		return nil, err
	}

	return Run(query, root), nil
}

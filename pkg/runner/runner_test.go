package runner

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type testTab []struct {
	query   string
	json    string
	results []interface{}
}

var testTabJSONChildlessRoot = testTab{
	// Empty query returns the root element.
	{``, `"root"`, []interface{}{"root"}},
	{``, `"1234"`, []interface{}{"1234"}},
	// All numbers parsed from JSON are float64 in Golang due to JavaScript's ambiguous Number type.
	{``, `1234`, []interface{}{1234.0}},
	{``, `1234.56`, []interface{}{1234.56}},
	{``, `true`, []interface{}{true}},
	{``, `null`, []interface{}{nil}},

	// All of these root elements are childless.
	{`*`, `"root"`, []interface{}{}},
	{`*`, `1234.56`, []interface{}{}},
	{`*`, `true`, []interface{}{}},
	{`*`, `null`, []interface{}{}},

	// The wildcard query should return the root element.
	{`**`, `"root"`, []interface{}{"root"}},
	{`**`, `1234.56`, []interface{}{1234.56}},
	{`**`, `true`, []interface{}{true}},
	{`**`, `null`, []interface{}{nil}},

	// There are no children, so there is nothing to return.
	{`child`, `"root"`, []interface{}{}},
	{`child`, `1234.56`, []interface{}{}},
	{`child`, `true`, []interface{}{}},
	{`child`, `null`, []interface{}{}},

	{`child.another.abc`, `"root"`, []interface{}{}},
	{`child.another.abc`, `1234.56`, []interface{}{}},
	{`child.another.abc`, `true`, []interface{}{}},

	{`first.second.0.last`, `"root"`, []interface{}{}},
	{`first.second.0.last`, `1234.56`, []interface{}{}},
	{`first.second.0.last`, `true`, []interface{}{}},

	{`child.`, `"root"`, []interface{}{}},
	{`child.`, `1234.56`, []interface{}{}},
	{`child.`, `true`, []interface{}{}},

	{`child..`, `"root"`, []interface{}{}},
	{`child..`, `1234.56`, []interface{}{}},
	{`child..`, `true`, []interface{}{}},

	{`child........`, `"root"`, []interface{}{}},
	{`child........`, `1234.56`, []interface{}{}},
	{`child........`, `true`, []interface{}{}},

	{`*.`, `"root"`, []interface{}{}},
	{`*.`, `1234.56`, []interface{}{}},
	{`*.`, `true`, []interface{}{}},

	{`**.`, `"root"`, []interface{}{}},
	{`**.`, `1234.56`, []interface{}{}},
	{`**.`, `true`, []interface{}{}},
}

var testTabJSONSingleChildRoot = testTab{
	// Empty array and map.
	{``, `[]`, []interface{}{[]interface{}{}}},
	{``, `{}`, []interface{}{map[string]interface{}{}}},

	{`*`, `[]`, []interface{}{}},
	{`*`, `{}`, []interface{}{}},

	{`**`, `[]`, []interface{}{[]interface{}{}}},
	{`**`, `{}`, []interface{}{map[string]interface{}{}}},

	{`child`, `[]`, []interface{}{}},
	{`child`, `{}`, []interface{}{}},

	{`0`, `[]`, []interface{}{}},
	{`0`, `{}`, []interface{}{}},

	// Single-item array.
	{``, `[123]`, []interface{}{[]interface{}{123.0}}},
	{``, `[123.45]`, []interface{}{[]interface{}{123.45}}},
	{``, `["abc"]`, []interface{}{[]interface{}{"abc"}}},
	{``, `[false]`, []interface{}{[]interface{}{false}}},
	{``, `[null]`, []interface{}{[]interface{}{nil}}},
	{``, `[[]]`, []interface{}{[]interface{}{[]interface{}{}}}},
	{``, `[{}]`, []interface{}{[]interface{}{map[string]interface{}{}}}},

	{`*`, `[123]`, []interface{}{123.0}},
	{`*`, `["abc"]`, []interface{}{"abc"}},

	{`**`, `[123]`, []interface{}{[]interface{}{123.0}, 123.0}},
	{`**`, `["abc"]`, []interface{}{[]interface{}{"abc"}, "abc"}},

	{`child`, `[123]`, []interface{}{}},
	{`child`, `["abc"]`, []interface{}{}},

	{`0`, `[123]`, []interface{}{123.0}},
	{`0`, `["abc"]`, []interface{}{"abc"}},

	{`1`, `[123]`, []interface{}{}},
	{`1`, `["abc"]`, []interface{}{}},

	{`0.`, `[123]`, []interface{}{[]interface{}{123.0}}},
	{`0.`, `["abc"]`, []interface{}{[]interface{}{"abc"}}},

	{`0..`, `[123]`, []interface{}{}},
	{`0..`, `["abc"]`, []interface{}{}},

	// Single-item map.
	{``, `{"child": 123}`, []interface{}{map[string]interface{}{"child": 123.0}}},
	{``, `{"child": 123.45}`, []interface{}{map[string]interface{}{"child": 123.45}}},
	{``, `{"child": "value"}`, []interface{}{map[string]interface{}{"child": "value"}}},
	{``, `{"child": false}`, []interface{}{map[string]interface{}{"child": false}}},
	{``, `{"child": null}`, []interface{}{map[string]interface{}{"child": nil}}},
	{``, `{"child": []}`, []interface{}{map[string]interface{}{"child": []interface{}{}}}},
	{``, `{"child": {}}`, []interface{}{map[string]interface{}{"child": map[string]interface{}{}}}},

	{`*`, `{"child": 123}`, []interface{}{123.0}},
	{`*`, `{"child": "value"}`, []interface{}{"value"}},

	{`*.`, `{"child": 123}`, []interface{}{map[string]interface{}{"child": 123.0}}},
	{`*.`, `{"child": "value"}`, []interface{}{map[string]interface{}{"child": "value"}}},

	{`*..`, `{"child": 123}`, []interface{}{}},
	{`*..`, `{"child": "value"}`, []interface{}{}},

	{`**`, `{"child": 123}`, []interface{}{map[string]interface{}{"child": 123.0}, 123.0}},
	{`**`, `{"child": "value"}`, []interface{}{map[string]interface{}{"child": "value"}, "value"}},

	{`**.`, `{"child": 123}`, []interface{}{map[string]interface{}{"child": 123.0}}},
	{`**.`, `{"child": "value"}`, []interface{}{map[string]interface{}{"child": "value"}}},

	{`**..`, `{"child": 123}`, []interface{}{}},
	{`**..`, `{"child": "value"}`, []interface{}{}},

	{`child`, `{"child": 123}`, []interface{}{123.0}},
	{`child`, `{"child": "value"}`, []interface{}{"value"}},

	{`child.`, `{"child": 123}`, []interface{}{map[string]interface{}{"child": 123.0}}},
	{`child.`, `{"child": "value"}`, []interface{}{map[string]interface{}{"child": "value"}}},

	{`child..`, `{"child": 123}`, []interface{}{}},
	{`child..`, `{"child": "value"}`, []interface{}{}},

	{`another`, `{"child": 123}`, []interface{}{}},
	{`another`, `{"child": "value"}`, []interface{}{}},

	{`0`, `{"child": 123}`, []interface{}{}},
	{`0`, `{"child": "value"}`, []interface{}{}},

	{`1`, `{"child": 123}`, []interface{}{}},
	{`1`, `{"child": "value"}`, []interface{}{}},
}

const complexJSON string = `[
	{
		"name": "John Doe",
		"debt": 1000
	},
	{
		"name": "Jane Doe",
		"debt": 2000
	},
	{
		"name": "John Daniel",
		"debt": 0
	},
	{
		"name": "Robert Denver",
		"debt": 0
	},
	{
		"name": "Clark Denver",
		"debt": 10000
	}
]`

var testTabJSONEquality = testTab{
	{`*.debt=0..name`, complexJSON, []interface{}{"John Daniel", "Robert Denver"}},
	{`*.debt=1..name`, complexJSON, []interface{}{}},
	{`*.debt=10..name`, complexJSON, []interface{}{}},
	{`*.debt=100..name`, complexJSON, []interface{}{}},
	{`*.debt=1000..name`, complexJSON, []interface{}{"John Doe"}},
}

var testTabJSONEqualityInverted = testTab{
	{`*.debt!=0..name`, complexJSON, []interface{}{"John Doe", "Jane Doe", "Clark Denver"}},
	{`*.debt!=1..name`, complexJSON, []interface{}{"John Doe", "Jane Doe", "John Daniel", "Robert Denver", "Clark Denver"}},
}

var testTabJSONRegex = testTab{
	{`*.name~" Doe$"`, complexJSON, []interface{}{"John Doe", "Jane Doe"}},
	{`*.name~"^John "`, complexJSON, []interface{}{"John Doe", "John Daniel"}},
}

var testTabJSONRegexInverted = testTab{
	{`*.name!~" Doe$"`, complexJSON, []interface{}{"John Daniel", "Robert Denver", "Clark Denver"}},
	{`*.name!~"^John "`, complexJSON, []interface{}{"Jane Doe", "Robert Denver", "Clark Denver"}},
}

const complexYAML = `name: Go
on: [push, pull_request]
jobs:

  build:
    name: Build
    runs-on: ubuntu-latest
    steps:

    - name: Set up Go 1.13
      uses: actions/setup-go@v1
      with:
        go-version: 1.13
      id: go

    - name: Check out code into the Go module directory
      uses: actions/checkout@v1

    - name: Get dependencies
      run: go get -v -t -d ./...

    - name: Build
      run: go build -v ./...

    - name: Run tests
      run: go test -v ./...`

var testTabYAMLGeneral = testTab{
	{`name`, complexYAML, []interface{}{"Go"}},
	{`name=Go`, complexYAML, []interface{}{"Go"}},
	{`name~Go`, complexYAML, []interface{}{"Go"}},
	{`name~^Go$`, complexYAML, []interface{}{"Go"}},

	{`on`, complexYAML, []interface{}{[]interface{}{"push", "pull_request"}}},
	{`on.*~pu`, complexYAML, []interface{}{"push", "pull_request"}},
	{`on.*~^pu`, complexYAML, []interface{}{"push", "pull_request"}},
	{`on.*=push`, complexYAML, []interface{}{"push"}},
	{`on.*~push`, complexYAML, []interface{}{"push"}},
	{`on.*~^push$`, complexYAML, []interface{}{"push"}},
	{`on.*=push.`, complexYAML, []interface{}{[]interface{}{"push", "pull_request"}}},
	{`on.*~push.`, complexYAML, []interface{}{[]interface{}{"push", "pull_request"}}},
	{`on.*~^push$.`, complexYAML, []interface{}{[]interface{}{"push", "pull_request"}}},

	{`jobs.*.steps.*.run..name`, complexYAML, []interface{}{"Get dependencies", "Build", "Run tests"}},
	{`**.run..name`, complexYAML, []interface{}{"Get dependencies", "Build", "Run tests"}},

	// TODO: While this is currently the expected result, it is not the desired final result.
	// Some kind of filtering should happen to remove duplicates.
	{`on.*~^pu.`, complexYAML, []interface{}{[]interface{}{"push", "pull_request"}, []interface{}{"push", "pull_request"}}},
}

func runTests(t *testing.T, tab testTab, verboseName bool, runFunc func(string, string) ([]Element, error)) {
	for index, entry := range tab {
		var testName string
		if verboseName {
			testName = fmt.Sprintf("Query:`%s`, JSON:`%s`", entry.query, entry.json)
		} else {
			testName = strconv.Itoa(index)
		}

		t.Run(testName, func(t *testing.T) {
			results, err := runFunc(entry.query, entry.json)

			if err != nil {
				t.Error(err)
				return
			}

			if len(results) != len(entry.results) {
				t.Errorf("Unexpected number of results: %d instead of %d", len(results), len(entry.results))
				return
			}

			for i, expected := range entry.results {
				if reality := results[i].Value; !cmp.Equal(reality, expected) {
					t.Errorf("Unexpected result at index %d: `%#v` (%T) instead of `%#v` (%T)", i, reality, reality, expected, expected)
				}
			}
		})
	}
}

func runTestsJSON(t *testing.T, tab testTab, verboseName bool) {
	runTests(t, tab, verboseName, RunJsonString)
}

func runTestsYAML(t *testing.T, tab testTab, verboseName bool) {
	runTests(t, tab, verboseName, RunYamlString)
}

func TestRunJSONChildlessRoot(t *testing.T) {
	runTestsJSON(t, testTabJSONChildlessRoot, true)
}

func TestRunJSONSingleChildRoot(t *testing.T) {
	runTestsJSON(t, testTabJSONSingleChildRoot, true)
}

func TestRunJSONEquality(t *testing.T) {
	runTestsJSON(t, testTabJSONEquality, false)
}

func TestRunJSONEqualityInverted(t *testing.T) {
	runTestsJSON(t, testTabJSONEqualityInverted, false)
}

func TestRunJSONRegex(t *testing.T) {
	runTestsJSON(t, testTabJSONRegex, false)
}

func TestRunJSONRegexInverted(t *testing.T) {
	runTestsJSON(t, testTabJSONRegexInverted, false)
}

func TestRunYAMLGeneral(t *testing.T) {
	runTestsYAML(t, testTabYAMLGeneral, false)
}

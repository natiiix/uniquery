package runner

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type testTab []struct {
	query   string
	source  string
	results map[string]interface{}
}

var testTabJSONChildlessRoot = testTab{
	// Empty query returns the root element.
	{``, `"root"`, map[string]interface{}{``: "root"}},
	{``, `"1234"`, map[string]interface{}{``: "1234"}},
	// All numbers parsed from JSON are float64 in Golang due to JavaScript's ambiguous Number type.
	{``, `1234`, map[string]interface{}{``: 1234.0}},
	{``, `1234.56`, map[string]interface{}{``: 1234.56}},
	{``, `true`, map[string]interface{}{``: true}},
	{``, `null`, map[string]interface{}{``: nil}},

	// All of these root elements are childless.
	{`*`, `"root"`, map[string]interface{}{}},
	{`*`, `1234.56`, map[string]interface{}{}},
	{`*`, `true`, map[string]interface{}{}},
	{`*`, `null`, map[string]interface{}{}},

	// The wildcard query should return the root element.
	{`**`, `"root"`, map[string]interface{}{``: "root"}},
	{`**`, `1234.56`, map[string]interface{}{``: 1234.56}},
	{`**`, `true`, map[string]interface{}{``: true}},
	{`**`, `null`, map[string]interface{}{``: nil}},

	// There are no children, so there is nothing to return.
	{`child`, `"root"`, map[string]interface{}{}},
	{`child`, `1234.56`, map[string]interface{}{}},
	{`child`, `true`, map[string]interface{}{}},
	{`child`, `null`, map[string]interface{}{}},

	{`child.another.abc`, `"root"`, map[string]interface{}{}},
	{`child.another.abc`, `1234.56`, map[string]interface{}{}},
	{`child.another.abc`, `true`, map[string]interface{}{}},

	{`first.second.0.last`, `"root"`, map[string]interface{}{}},
	{`first.second.0.last`, `1234.56`, map[string]interface{}{}},
	{`first.second.0.last`, `true`, map[string]interface{}{}},

	{`child.`, `"root"`, map[string]interface{}{}},
	{`child.`, `1234.56`, map[string]interface{}{}},
	{`child.`, `true`, map[string]interface{}{}},

	{`child..`, `"root"`, map[string]interface{}{}},
	{`child..`, `1234.56`, map[string]interface{}{}},
	{`child..`, `true`, map[string]interface{}{}},

	{`child........`, `"root"`, map[string]interface{}{}},
	{`child........`, `1234.56`, map[string]interface{}{}},
	{`child........`, `true`, map[string]interface{}{}},

	{`*.`, `"root"`, map[string]interface{}{}},
	{`*.`, `1234.56`, map[string]interface{}{}},
	{`*.`, `true`, map[string]interface{}{}},

	{`**.`, `"root"`, map[string]interface{}{}},
	{`**.`, `1234.56`, map[string]interface{}{}},
	{`**.`, `true`, map[string]interface{}{}},
}

var testTabJSONSingleChildRoot = testTab{
	// Empty array and map.
	{``, `[]`, map[string]interface{}{``: []interface{}{}}},
	{``, `{}`, map[string]interface{}{``: map[string]interface{}{}}},

	{`*`, `[]`, map[string]interface{}{}},
	{`*`, `{}`, map[string]interface{}{}},

	{`**`, `[]`, map[string]interface{}{``: []interface{}{}}},
	{`**`, `{}`, map[string]interface{}{``: map[string]interface{}{}}},

	{`child`, `[]`, map[string]interface{}{}},
	{`child`, `{}`, map[string]interface{}{}},

	{`0`, `[]`, map[string]interface{}{}},
	{`0`, `{}`, map[string]interface{}{}},

	// Single-item array.
	{``, `[123]`, map[string]interface{}{``: []interface{}{123.0}}},
	{``, `[123.45]`, map[string]interface{}{``: []interface{}{123.45}}},
	{``, `["abc"]`, map[string]interface{}{``: []interface{}{"abc"}}},
	{``, `[false]`, map[string]interface{}{``: []interface{}{false}}},
	{``, `[null]`, map[string]interface{}{``: []interface{}{nil}}},
	{``, `[[]]`, map[string]interface{}{``: []interface{}{[]interface{}{}}}},
	{``, `[{}]`, map[string]interface{}{``: []interface{}{map[string]interface{}{}}}},

	{`*`, `[123]`, map[string]interface{}{`0`: 123.0}},
	{`*`, `["abc"]`, map[string]interface{}{`0`: "abc"}},

	{`**`, `[123]`, map[string]interface{}{``: []interface{}{123.0}, `0`: 123.0}},
	{`**`, `["abc"]`, map[string]interface{}{``: []interface{}{"abc"}, `0`: "abc"}},

	{`child`, `[123]`, map[string]interface{}{}},
	{`child`, `["abc"]`, map[string]interface{}{}},

	{`0`, `[123]`, map[string]interface{}{`0`: 123.0}},
	{`0`, `["abc"]`, map[string]interface{}{`0`: "abc"}},

	{`1`, `[123]`, map[string]interface{}{}},
	{`1`, `["abc"]`, map[string]interface{}{}},

	{`0.`, `[123]`, map[string]interface{}{``: []interface{}{123.0}}},
	{`0.`, `["abc"]`, map[string]interface{}{``: []interface{}{"abc"}}},

	{`0..`, `[123]`, map[string]interface{}{}},
	{`0..`, `["abc"]`, map[string]interface{}{}},

	// Single-item map.
	{``, `{"child": 123}`, map[string]interface{}{``: map[string]interface{}{"child": 123.0}}},
	{``, `{"child": 123.45}`, map[string]interface{}{``: map[string]interface{}{"child": 123.45}}},
	{``, `{"child": "value"}`, map[string]interface{}{``: map[string]interface{}{"child": "value"}}},
	{``, `{"child": false}`, map[string]interface{}{``: map[string]interface{}{"child": false}}},
	{``, `{"child": null}`, map[string]interface{}{``: map[string]interface{}{"child": nil}}},
	{``, `{"child": []}`, map[string]interface{}{``: map[string]interface{}{"child": []interface{}{}}}},
	{``, `{"child": {}}`, map[string]interface{}{``: map[string]interface{}{"child": map[string]interface{}{}}}},

	{`*`, `{"child": 123}`, map[string]interface{}{`"child"`: 123.0}},
	{`*`, `{"child": "value"}`, map[string]interface{}{`"child"`: "value"}},

	{`*.`, `{"child": 123}`, map[string]interface{}{``: map[string]interface{}{"child": 123.0}}},
	{`*.`, `{"child": "value"}`, map[string]interface{}{``: map[string]interface{}{"child": "value"}}},

	{`*..`, `{"child": 123}`, map[string]interface{}{}},
	{`*..`, `{"child": "value"}`, map[string]interface{}{}},

	{`**`, `{"child": 123}`, map[string]interface{}{``: map[string]interface{}{"child": 123.0}, `"child"`: 123.0}},
	{`**`, `{"child": "value"}`, map[string]interface{}{``: map[string]interface{}{"child": "value"}, `"child"`: "value"}},

	{`**.`, `{"child": 123}`, map[string]interface{}{``: map[string]interface{}{"child": 123.0}}},
	{`**.`, `{"child": "value"}`, map[string]interface{}{``: map[string]interface{}{"child": "value"}}},

	{`**..`, `{"child": 123}`, map[string]interface{}{}},
	{`**..`, `{"child": "value"}`, map[string]interface{}{}},

	{`child`, `{"child": 123}`, map[string]interface{}{`"child"`: 123.0}},
	{`child`, `{"child": "value"}`, map[string]interface{}{`"child"`: "value"}},

	{`child.`, `{"child": 123}`, map[string]interface{}{``: map[string]interface{}{"child": 123.0}}},
	{`child.`, `{"child": "value"}`, map[string]interface{}{``: map[string]interface{}{"child": "value"}}},

	{`child..`, `{"child": 123}`, map[string]interface{}{}},
	{`child..`, `{"child": "value"}`, map[string]interface{}{}},

	{`another`, `{"child": 123}`, map[string]interface{}{}},
	{`another`, `{"child": "value"}`, map[string]interface{}{}},

	{`0`, `{"child": 123}`, map[string]interface{}{}},
	{`0`, `{"child": "value"}`, map[string]interface{}{}},

	{`1`, `{"child": 123}`, map[string]interface{}{}},
	{`1`, `{"child": "value"}`, map[string]interface{}{}},

	{`0`, `{"0": 123}`, map[string]interface{}{`"0"`: 123.0}},
	{`0`, `{"0": "abc"}`, map[string]interface{}{`"0"`: "abc"}},

	{`1`, `{"0": 123}`, map[string]interface{}{}},
	{`1`, `{"0": "abc"}`, map[string]interface{}{}},
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
	{`*.debt=0..name`, complexJSON, map[string]interface{}{`2."name"`: "John Daniel", `3."name"`: "Robert Denver"}},
	{`*.debt=1..name`, complexJSON, map[string]interface{}{}},
	{`*.debt=10..name`, complexJSON, map[string]interface{}{}},
	{`*.debt=100..name`, complexJSON, map[string]interface{}{}},
	{`*.debt=1000..name`, complexJSON, map[string]interface{}{`0."name"`: "John Doe"}},
}

var testTabJSONEqualityInverted = testTab{
	{`*.debt!=0..name`, complexJSON, map[string]interface{}{`0."name"`: "John Doe", `1."name"`: "Jane Doe", `4."name"`: "Clark Denver"}},
	{`*.debt!=1..name`, complexJSON, map[string]interface{}{`0."name"`: "John Doe", `1."name"`: "Jane Doe", `2."name"`: "John Daniel", `3."name"`: "Robert Denver", `4."name"`: "Clark Denver"}},
}

var testTabJSONRegex = testTab{
	{`*.name~" Doe$"`, complexJSON, map[string]interface{}{`0."name"`: "John Doe", `1."name"`: "Jane Doe"}},
	{`*.name~"^John "`, complexJSON, map[string]interface{}{`0."name"`: "John Doe", `2."name"`: "John Daniel"}},
}

var testTabJSONRegexInverted = testTab{
	{`*.name!~" Doe$"`, complexJSON, map[string]interface{}{`2."name"`: "John Daniel", `3."name"`: "Robert Denver", `4."name"`: "Clark Denver"}},
	{`*.name!~"^John "`, complexJSON, map[string]interface{}{`1."name"`: "Jane Doe", `3."name"`: "Robert Denver", `4."name"`: "Clark Denver"}},
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
	{`name`, complexYAML, map[string]interface{}{`"name"`: "Go"}},
	{`name=Go`, complexYAML, map[string]interface{}{`"name"`: "Go"}},
	{`name~Go`, complexYAML, map[string]interface{}{`"name"`: "Go"}},
	{`name~^Go$`, complexYAML, map[string]interface{}{`"name"`: "Go"}},

	{`on`, complexYAML, map[string]interface{}{`true`: []interface{}{"push", "pull_request"}}},
	{`on.*~pu`, complexYAML, map[string]interface{}{`true.0`: "push", `true.1`: "pull_request"}},
	{`on.*~^pu`, complexYAML, map[string]interface{}{`true.0`: "push", `true.1`: "pull_request"}},
	{`on.*=push`, complexYAML, map[string]interface{}{`true.0`: "push"}},
	{`on.*~push`, complexYAML, map[string]interface{}{`true.0`: "push"}},
	{`on.*~^push$`, complexYAML, map[string]interface{}{`true.0`: "push"}},
	{`on.*=push.`, complexYAML, map[string]interface{}{`true`: []interface{}{"push", "pull_request"}}},
	{`on.*~push.`, complexYAML, map[string]interface{}{`true`: []interface{}{"push", "pull_request"}}},
	{`on.*~^push$.`, complexYAML, map[string]interface{}{`true`: []interface{}{"push", "pull_request"}}},

	{`jobs.*.steps.*.run..name`, complexYAML, map[string]interface{}{`"jobs"."build"."steps".2."name"`: "Get dependencies", `"jobs"."build"."steps".3."name"`: "Build", `"jobs"."build"."steps".4."name"`: "Run tests"}},
	{`**.run..name`, complexYAML, map[string]interface{}{`"jobs"."build"."steps".2."name"`: "Get dependencies", `"jobs"."build"."steps".3."name"`: "Build", `"jobs"."build"."steps".4."name"`: "Run tests"}},

	// NOTE: This checks that duplicate results are filtered out.
	{`on.*~^pu.`, complexYAML, map[string]interface{}{`true`: []interface{}{"push", "pull_request"}}},
}

func runTests(t *testing.T, tab testTab, verboseName bool, runFunc func(string, string) (map[string]Element, error)) {
	for index, entry := range tab {
		var testName string
		if verboseName {
			testName = fmt.Sprintf("Query:`%s`, JSON:`%s`", entry.query, entry.source)
		} else {
			testName = strconv.Itoa(index)
		}

		t.Run(testName, func(t *testing.T) {
			results, err := runFunc(entry.query, entry.source)

			if err != nil {
				t.Error(err)
				return
			}

			if len(results) != len(entry.results) {
				t.Errorf("Unexpected number of results: %d instead of %d", len(results), len(entry.results))
				return
			}

			// for k, expected := range entry.results {
			// 	if reality, exists := results[k]; !exists {
			// 		t.Errorf("Expected element path `%s` missing from results -- %v", k, results)
			// 	} else if !cmp.Equal(reality, expected) {
			// 		t.Errorf("Unexpected value of result with path `%s`: `%#v` (%T) instead of `%#v` (%T)", k, reality, reality, expected, expected)
			// 	}
			// }

			for k, expected := range results {
				if reality, exists := entry.results[k]; !exists {
					t.Errorf("Expected element path `%s` missing from results -- %v", k, results)
				} else if !cmp.Equal(reality, expected.Value) {
					t.Errorf("Unexpected value of result with path `%s`: `%#v` (%T) instead of `%#v` (%T)", k, reality, reality, expected, expected)
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

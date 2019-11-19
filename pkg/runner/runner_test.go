package runner

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type testTab []struct {
	query   string
	json    string
	results []interface{}
}

var testTabChildlessRoot = testTab{
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

var testTabParentRoot = testTab{
	// Empty array and map.
	{``, `[]`, []interface{}{[]interface{}{}}},
	{``, `{}`, []interface{}{make(map[string]interface{})}},

	{`*`, `[]`, []interface{}{}},
	{`*`, `{}`, []interface{}{}},

	{`**`, `[]`, []interface{}{[]interface{}{}}},
	{`**`, `{}`, []interface{}{make(map[string]interface{})}},

	{`child`, `[]`, []interface{}{}},
	{`child`, `{}`, []interface{}{}},

	{`0`, `[]`, []interface{}{}},
	{`0`, `{}`, []interface{}{}},
}

func runTests(t *testing.T, tab testTab) {
	for _, entry := range tab {
		t.Run(fmt.Sprintf("Query:`%s`, JSON:`%s`", entry.query, entry.json), func(t *testing.T) {
			results, err := RunJsonString(entry.query, entry.json)

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

func TestRunChildlessRoot(t *testing.T) {
	runTests(t, testTabChildlessRoot)
}

func TestRunParentRoot(t *testing.T) {
	runTests(t, testTabParentRoot)
}

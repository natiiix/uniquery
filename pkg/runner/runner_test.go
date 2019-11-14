package runner

import (
	"fmt"
	"testing"
)

var testTab = []struct {
	query   string
	json    string
	results []interface{}
}{
	// Empty query returns the root element.
	{``, `"root"`, []interface{}{"root"}},
	{``, `"1234"`, []interface{}{"1234"}},
	// All numbers parsed from JSON are float64 in Golang due to JavaScript's ambiguous Number type.
	{``, `1234`, []interface{}{1234.0}},
	{``, `1234.56`, []interface{}{1234.56}},
	{``, `true`, []interface{}{true}},

	// All of these root elements are childless.
	{`*`, `"root"`, []interface{}{}},
	{`*`, `1234.56`, []interface{}{}},
	{`*`, `true`, []interface{}{}},

	// The wildcard query should return the root element.
	{`**`, `"root"`, []interface{}{"root"}},
	{`**`, `1234.56`, []interface{}{1234.56}},
	{`**`, `true`, []interface{}{true}},

	// There are no children, so there is nothing to return.
	{`child`, `"root"`, []interface{}{}},
	{`child`, `1234.56`, []interface{}{}},
	{`child`, `true`, []interface{}{}},

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

func TestRun(t *testing.T) {
	for _, tt := range testTab {
		t.Run(fmt.Sprintf(`Query: "%s", JSON: "%s"`, tt.query, tt.json), func(t *testing.T) {
			results, err := RunJsonString(tt.query, tt.json)

			if err != nil {
				t.Error(err)
				return
			}

			if len(results) != len(tt.results) {
				t.Errorf("Unexpected number of results: %d instead of %d", len(results), len(tt.results))
				return
			}

			for i, expected := range tt.results {
				if reality := results[i].Value; reality != expected {
					t.Errorf(`Unexpected result at index %d: "%#v" (%T) instead of "%#v" (%T)`, i, reality, reality, expected, expected)
				}
			}
		})
	}
}

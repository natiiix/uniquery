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
	// All numbers in JSON are float64 in Golang due to JavaScript's Number type.
	{``, `1234`, []interface{}{1234.0}},
	{``, `1234.56`, []interface{}{1234.56}},
	{``, `true`, []interface{}{true}},
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
				t.Errorf("Unexpected number of results: %d instead of %d", len(tt.results), len(results))
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

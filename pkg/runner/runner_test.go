package runner

import (
	"testing"
)

func TestRun(t *testing.T) {
	results, err := RunJsonString(``, `"Test"`)

	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 1 {
		t.Fatal("Unexpected number of results:", len(results))
	}

	if results[0].Value != "Test" {
		t.Error("Unexpected value of the result:", results[0].Value)
	}
}

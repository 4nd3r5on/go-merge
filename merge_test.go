package merge_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/4nd3r5on/go-merge"
)

// TestCase represents a single merge operation test case
type TestCase struct {
	Name      string
	Mode      merge.Mode
	Original  any
	Merge     any
	Expected  any
	ShouldErr bool
	ErrMsg    string
}

// RunTestCase executes a single test case
func RunTestCase(t *testing.T, tc TestCase) {
	t.Helper()

	result, err := merge.Data(tc.Mode, tc.Original, tc.Merge)

	if tc.ShouldErr {
		if err == nil {
			t.Errorf("Expected error but got none")
			return
		}
		if tc.ErrMsg != "" && !strings.Contains(err.Error(), tc.ErrMsg) {
			t.Errorf("Expected error containing %q, got: %v", tc.ErrMsg, err)
		}
		return
	}

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if !reflect.DeepEqual(result, tc.Expected) {
		t.Errorf("Result mismatch:\nGot:      %s\nExpected: %s",
			toJSON(result), toJSON(tc.Expected))
	}
}

// TableTest executes multiple test cases
func TableTest(t *testing.T, cases []TestCase) {
	t.Helper()
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			RunTestCase(t, tc)
		})
	}
}

// M creates a map[string]any for test data
func M(pairs ...any) map[string]any {
	if len(pairs)%2 != 0 {
		panic("M() requires even number of arguments")
	}

	m := make(map[string]any)
	for i := 0; i < len(pairs); i += 2 {
		m[pairs[i].(string)] = pairs[i+1]
	}
	return m
}

func toJSON(v any) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintf("%#v", v)
	}
	return string(b)
}

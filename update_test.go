package merge_test

import (
	"testing"

	"github.com/4nd3r5on/go-merge"
)

func TestUpdateMode_Maps(t *testing.T) {
	cases := []TestCase{
		// Basic map operations
		{
			Name:     "update existing keys in flat map",
			Mode:     merge.ModeUpdate,
			Original: map[string]any{"a": 1, "b": 2},
			Merge:    map[string]any{"a": 10, "b": 20},
			Expected: map[string]any{"a": 10, "b": 20},
		},
		{
			Name:     "ignore new keys in flat map",
			Mode:     merge.ModeUpdate,
			Original: map[string]any{"a": 1, "b": 2},
			Merge:    map[string]any{"a": 10, "c": 30},
			Expected: map[string]any{"a": 10, "b": 2},
		},
		{
			Name:     "update some keys and ignore new keys",
			Mode:     merge.ModeUpdate,
			Original: map[string]any{"a": 1, "b": 2, "c": 3},
			Merge:    map[string]any{"b": 20, "d": 40, "e": 50},
			Expected: map[string]any{"a": 1, "b": 20, "c": 3},
		},
		{
			Name:     "merge with empty map does nothing",
			Mode:     merge.ModeUpdate,
			Original: map[string]any{"a": 1, "b": 2},
			Merge:    map[string]any{},
			Expected: map[string]any{"a": 1, "b": 2},
		},
		{
			Name:     "empty original with non-empty merge stays empty",
			Mode:     merge.ModeUpdate,
			Original: map[string]any{},
			Merge:    map[string]any{"a": 1, "b": 2},
			Expected: map[string]any{},
		},
		{
			Name:     "both maps empty",
			Mode:     merge.ModeUpdate,
			Original: map[string]any{},
			Merge:    map[string]any{},
			Expected: map[string]any{},
		},

		// Nested map operations
		{
			Name:     "recursively merge nested maps - update nested values",
			Mode:     merge.ModeUpdate,
			Original: map[string]any{"a": map[string]any{"x": 1, "y": 2}},
			Merge:    map[string]any{"a": map[string]any{"x": 10}},
			Expected: map[string]any{"a": map[string]any{"x": 10, "y": 2}},
		},
		{
			Name:     "recursively merge nested maps - ignore new nested keys",
			Mode:     merge.ModeUpdate,
			Original: map[string]any{"a": map[string]any{"x": 1, "y": 2}},
			Merge:    map[string]any{"a": map[string]any{"z": 30}},
			Expected: map[string]any{"a": map[string]any{"x": 1, "y": 2}},
		},
		{
			Name:     "deeply nested map update",
			Mode:     merge.ModeUpdate,
			Original: map[string]any{"a": map[string]any{"b": map[string]any{"c": 1, "d": 2}}},
			Merge:    map[string]any{"a": map[string]any{"b": map[string]any{"c": 10}}},
			Expected: map[string]any{"a": map[string]any{"b": map[string]any{"c": 10, "d": 2}}},
		},
		{
			Name:     "multiple nested maps with mixed updates",
			Mode:     merge.ModeUpdate,
			Original: map[string]any{"a": map[string]any{"x": 1}, "b": map[string]any{"y": 2}},
			Merge:    map[string]any{"a": map[string]any{"x": 10}, "b": map[string]any{"y": 20}},
			Expected: map[string]any{"a": map[string]any{"x": 10}, "b": map[string]any{"y": 20}},
		},

		// Type changes and replacements
		{
			Name:     "replace primitive with map is ignored (new structure)",
			Mode:     merge.ModeUpdate,
			Original: map[string]any{"a": 1},
			Merge:    map[string]any{"a": map[string]any{"x": 10}},
			Expected: map[string]any{"a": map[string]any{"x": 10}},
		},
		{
			Name:     "replace array with map",
			Mode:     merge.ModeUpdate,
			Original: map[string]any{"a": []int{1, 2, 3}},
			Merge:    map[string]any{"a": map[string]any{"x": 10}},
			Expected: map[string]any{"a": map[string]any{"x": 10}},
		},

		// Nil handling
		{
			Name:     "update with nil value",
			Mode:     merge.ModeUpdate,
			Original: map[string]any{"a": 1, "b": 2},
			Merge:    map[string]any{"a": nil},
			Expected: map[string]any{"a": nil, "b": 2},
		},
		{
			Name:     "nested map with nil values",
			Mode:     merge.ModeUpdate,
			Original: map[string]any{"a": map[string]any{"x": nil, "y": 2}},
			Merge:    map[string]any{"a": map[string]any{"x": 10}},
			Expected: map[string]any{"a": map[string]any{"x": nil, "y": 2}},
		},
		{
			Name:     "original map is nil",
			Mode:     merge.ModeUpdate,
			Original: nil,
			Merge:    map[string]any{"a": 1},
			Expected: nil,
		},
		{
			Name:     "merge map is nil",
			Mode:     merge.ModeUpdate,
			Original: map[string]any{"a": 1},
			Merge:    nil,
			Expected: map[string]any{"a": 1},
		},

		// Different map key types
		{
			Name:     "map with integer keys",
			Mode:     merge.ModeUpdate,
			Original: map[int]any{1: "a", 2: "b"},
			Merge:    map[int]any{1: "x", 3: "c"},
			Expected: map[int]any{1: "x", 2: "b"},
		},
		{
			Name:     "nested maps with mixed key types",
			Mode:     merge.ModeUpdate,
			Original: map[string]any{"a": map[int]any{1: "x", 2: "y"}},
			Merge:    map[string]any{"a": map[int]any{1: "updated"}},
			Expected: map[string]any{"a": map[int]any{1: "updated", 2: "y"}},
		},

		// Complex nested structures
		{
			Name: "complex nested structure with multiple levels",
			Mode: merge.ModeUpdate,
			Original: map[string]any{
				"level1": map[string]any{
					"level2": map[string]any{
						"level3": map[string]any{
							"value": 1,
							"other": 2,
						},
					},
				},
			},
			Merge: map[string]any{
				"level1": map[string]any{
					"level2": map[string]any{
						"level3": map[string]any{
							"value": 100,
						},
					},
				},
			},
			Expected: map[string]any{
				"level1": map[string]any{
					"level2": map[string]any{
						"level3": map[string]any{
							"value": 100,
							"other": 2,
						},
					},
				},
			},
		},
		{
			Name: "multiple nested maps at same level",
			Mode: merge.ModeUpdate,
			Original: map[string]any{
				"config": map[string]any{"timeout": 30, "retries": 3},
				"data":   map[string]any{"count": 100, "size": 200},
			},
			Merge: map[string]any{
				"config": map[string]any{"timeout": 60},
				"data":   map[string]any{"count": 150, "new": 300},
			},
			Expected: map[string]any{
				"config": map[string]any{"timeout": 60, "retries": 3},
				"data":   map[string]any{"count": 150, "size": 200},
			},
		},
		{
			Name: "partial nested structure updates",
			Mode: merge.ModeUpdate,
			Original: map[string]any{
				"a": map[string]any{"x": 1},
				"b": map[string]any{"y": 2},
				"c": 3,
			},
			Merge: map[string]any{
				"a": map[string]any{"x": 10},
				"d": map[string]any{"z": 40},
			},
			Expected: map[string]any{
				"a": map[string]any{"x": 10},
				"b": map[string]any{"y": 2},
				"c": 3,
			},
		},

		// Edge cases with empty nested structures
		{
			Name:     "nested empty map in original",
			Mode:     merge.ModeUpdate,
			Original: map[string]any{"a": map[string]any{}},
			Merge:    map[string]any{"a": map[string]any{"x": 1}},
			Expected: map[string]any{"a": map[string]any{}},
		},
		{
			Name:     "nested empty map in merge",
			Mode:     merge.ModeUpdate,
			Original: map[string]any{"a": map[string]any{"x": 1}},
			Merge:    map[string]any{"a": map[string]any{}},
			Expected: map[string]any{"a": map[string]any{"x": 1}},
		},

		// Maps containing arrays (mixed structures)
		{
			Name:     "map with array value - existing key",
			Mode:     merge.ModeUpdate,
			Original: map[string]any{"arr": []any{1, 2, 3}},
			Merge:    map[string]any{"arr": []any{4, 5}},
			Expected: map[string]any{"arr": []any{1, 2, 3, 4, 5}},
		},
		{
			Name:     "nested map with array updates",
			Mode:     merge.ModeUpdate,
			Original: map[string]any{"a": map[string]any{"arr": []any{1, 2}}},
			Merge:    map[string]any{"a": map[string]any{"arr": []any{3}}},
			Expected: map[string]any{"a": map[string]any{"arr": []any{1, 2, 3}}},
		},

		// Zero values
		{
			Name:     "update with zero values",
			Mode:     merge.ModeUpdate,
			Original: map[string]any{"a": 10, "b": "text", "c": true},
			Merge:    map[string]any{"a": 0, "b": "", "c": false},
			Expected: map[string]any{"a": 0, "b": "", "c": false},
		},
		{
			Name:     "nested zero values",
			Mode:     merge.ModeUpdate,
			Original: map[string]any{"config": map[string]any{"value": 100}},
			Merge:    map[string]any{"config": map[string]any{"value": 0}},
			Expected: map[string]any{"config": map[string]any{"value": 0}},
		},
	}
	TableTest(t, cases)
}

func TestUpdateMode_Arrays(t *testing.T) {
	cases := []TestCase{
		// Basic array operations - appending unique values
		{
			Name:     "append unique values to array",
			Mode:     merge.ModeUpdate,
			Original: []any{1, 2, 3},
			Merge:    []any{4, 5},
			Expected: []any{1, 2, 3, 4, 5},
		},
		{
			Name:     "ignore duplicate values",
			Mode:     merge.ModeUpdate,
			Original: []any{1, 2, 3},
			Merge:    []any{2, 3, 4},
			Expected: []any{1, 2, 3, 4},
		},
		{
			Name:     "all values are duplicates",
			Mode:     merge.ModeUpdate,
			Original: []any{1, 2, 3},
			Merge:    []any{1, 2, 3},
			Expected: []any{1, 2, 3},
		},
		{
			Name:     "merge empty array does nothing",
			Mode:     merge.ModeUpdate,
			Original: []any{1, 2, 3},
			Merge:    []any{},
			Expected: []any{1, 2, 3},
		},
		{
			Name:     "append to empty array",
			Mode:     merge.ModeUpdate,
			Original: []any{},
			Merge:    []any{1, 2, 3},
			Expected: []any{1, 2, 3},
		},
		{
			Name:     "both arrays empty",
			Mode:     merge.ModeUpdate,
			Original: []any{},
			Merge:    []any{},
			Expected: []any{},
		},

		// Different data types
		{
			Name:     "string arrays",
			Mode:     merge.ModeUpdate,
			Original: []any{"a", "b", "c"},
			Merge:    []any{"c", "d", "e"},
			Expected: []any{"a", "b", "c", "d", "e"},
		},
		{
			Name:     "boolean arrays",
			Mode:     merge.ModeUpdate,
			Original: []any{true, false},
			Merge:    []any{false, true, false},
			Expected: []any{true, false},
		},
		{
			Name:     "mixed type arrays",
			Mode:     merge.ModeUpdate,
			Original: []any{1, "a", true},
			Merge:    []any{"b", 2, false, 1},
			Expected: []any{1, "a", true, "b", 2, false},
		},
		{
			Name:     "float arrays",
			Mode:     merge.ModeUpdate,
			Original: []any{1.5, 2.5, 3.5},
			Merge:    []any{2.5, 4.5},
			Expected: []any{1.5, 2.5, 3.5, 4.5},
		},

		// Nil handling
		{
			Name:     "array with nil values",
			Mode:     merge.ModeUpdate,
			Original: []any{1, nil, 3},
			Merge:    []any{nil, 4},
			Expected: []any{1, nil, 3, 4},
		},
		{
			Name:     "append nil to array without nil",
			Mode:     merge.ModeUpdate,
			Original: []any{1, 2, 3},
			Merge:    []any{nil, 4},
			Expected: []any{1, 2, 3, nil, 4},
		},
		{
			Name:     "multiple nils in merge",
			Mode:     merge.ModeUpdate,
			Original: []any{1},
			Merge:    []any{nil, nil, 2},
			Expected: []any{1, nil, 2},
		},
		{
			Name:     "original is nil array",
			Mode:     merge.ModeUpdate,
			Original: nil,
			Merge:    []any{1, 2, 3},
			Expected: nil,
		},
		{
			Name:     "merge is nil array",
			Mode:     merge.ModeUpdate,
			Original: []any{1, 2, 3},
			Merge:    nil,
			Expected: []any{1, 2, 3},
		},

		// Complex objects in arrays
		{
			Name:     "arrays of maps - deep equality check",
			Mode:     merge.ModeUpdate,
			Original: []any{map[string]any{"id": 1, "name": "a"}},
			Merge:    []any{map[string]any{"id": 1, "name": "a"}, map[string]any{"id": 2, "name": "b"}},
			Expected: []any{map[string]any{"id": 1, "name": "a"}, map[string]any{"id": 2, "name": "b"}},
		},
		{
			Name:     "arrays of maps - different values are unique",
			Mode:     merge.ModeUpdate,
			Original: []any{map[string]any{"id": 1}},
			Merge:    []any{map[string]any{"id": 2}},
			Expected: []any{map[string]any{"id": 1}, map[string]any{"id": 2}},
		},
		{
			Name:     "nested arrays",
			Mode:     merge.ModeUpdate,
			Original: []any{[]any{1, 2}, []any{3, 4}},
			Merge:    []any{[]any{1, 2}, []any{5, 6}},
			Expected: []any{[]any{1, 2}, []any{3, 4}, []any{5, 6}},
		},
		{
			Name:     "arrays with nested maps",
			Mode:     merge.ModeUpdate,
			Original: []any{map[string]any{"x": 1, "y": 2}},
			Merge:    []any{map[string]any{"x": 1, "y": 3}},
			Expected: []any{map[string]any{"x": 1, "y": 2}, map[string]any{"x": 1, "y": 3}},
		},

		// Zero values
		{
			Name:     "zero value integers",
			Mode:     merge.ModeUpdate,
			Original: []any{0, 1, 2},
			Merge:    []any{0, 3},
			Expected: []any{0, 1, 2, 3},
		},
		{
			Name:     "empty strings",
			Mode:     merge.ModeUpdate,
			Original: []any{"", "a", "b"},
			Merge:    []any{"", "c"},
			Expected: []any{"", "a", "b", "c"},
		},
		{
			Name:     "false boolean values",
			Mode:     merge.ModeUpdate,
			Original: []any{false, true},
			Merge:    []any{false},
			Expected: []any{false, true},
		},
		{
			Name:     "zero floats",
			Mode:     merge.ModeUpdate,
			Original: []any{0.0, 1.5},
			Merge:    []any{0.0, 2.5},
			Expected: []any{0.0, 1.5, 2.5},
		},

		// Deep equality edge cases
		{
			Name:     "deeply nested identical structures",
			Mode:     merge.ModeUpdate,
			Original: []any{map[string]any{"a": map[string]any{"b": []any{1, 2}}}},
			Merge:    []any{map[string]any{"a": map[string]any{"b": []any{1, 2}}}},
			Expected: []any{map[string]any{"a": map[string]any{"b": []any{1, 2}}}},
		},
		{
			Name:     "deeply nested different structures",
			Mode:     merge.ModeUpdate,
			Original: []any{map[string]any{"a": map[string]any{"b": []any{1, 2}}}},
			Merge:    []any{map[string]any{"a": map[string]any{"b": []any{1, 3}}}},
			Expected: []any{
				map[string]any{"a": map[string]any{"b": []any{1, 2}}},
				map[string]any{"a": map[string]any{"b": []any{1, 3}}},
			},
		},
		{
			Name:     "maps with different key order but same content",
			Mode:     merge.ModeUpdate,
			Original: []any{map[string]any{"x": 1, "y": 2}},
			Merge:    []any{map[string]any{"y": 2, "x": 1}},
			Expected: []any{map[string]any{"x": 1, "y": 2}},
		},

		// Large arrays
		{
			Name:     "large array with many duplicates",
			Mode:     merge.ModeUpdate,
			Original: []any{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			Merge:    []any{5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
			Expected: []any{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
		},
		{
			Name:     "large array with no duplicates",
			Mode:     merge.ModeUpdate,
			Original: []any{1, 2, 3, 4, 5},
			Merge:    []any{6, 7, 8, 9, 10},
			Expected: []any{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		},

		// Order preservation
		{
			Name:     "maintains original order and appends in merge order",
			Mode:     merge.ModeUpdate,
			Original: []any{3, 1, 2},
			Merge:    []any{5, 4, 6},
			Expected: []any{3, 1, 2, 5, 4, 6},
		},
		{
			Name:     "duplicate detection maintains order",
			Mode:     merge.ModeUpdate,
			Original: []any{"z", "a", "m"},
			Merge:    []any{"a", "b", "z", "c"},
			Expected: []any{"z", "a", "m", "b", "c"},
		},

		// Edge cases with complex equality
		{
			Name:     "similar but not equal maps",
			Mode:     merge.ModeUpdate,
			Original: []any{map[string]any{"id": 1, "value": "a"}},
			Merge:    []any{map[string]any{"id": 1, "value": "b"}},
			Expected: []any{
				map[string]any{"id": 1, "value": "a"},
				map[string]any{"id": 1, "value": "b"},
			},
		},
		{
			Name:     "arrays with different lengths",
			Mode:     merge.ModeUpdate,
			Original: []any{[]any{1, 2, 3}},
			Merge:    []any{[]any{1, 2}},
			Expected: []any{[]any{1, 2, 3}, []any{1, 2}},
		},
		{
			Name:     "single element arrays",
			Mode:     merge.ModeUpdate,
			Original: []any{1},
			Merge:    []any{2},
			Expected: []any{1, 2},
		},
		{
			Name:     "single duplicate element",
			Mode:     merge.ModeUpdate,
			Original: []any{1},
			Merge:    []any{1},
			Expected: []any{1},
		},
	}
	TableTest(t, cases)
}

package merge_test

import (
	"testing"

	"github.com/4nd3r5on/go-merge"
)

func TestInsertMode_Maps(t *testing.T) {
	cases := []TestCase{
		{
			Name:     "Insert new keys only",
			Mode:     merge.MergeModeInsert,
			Original: M("a", 1, "b", 2),
			Merge:    M("b", 3, "c", 4),
			Expected: M("a", 1, "b", 2, "c", 4),
		},
		{
			Name:     "Insert into empty map",
			Mode:     merge.MergeModeInsert,
			Original: M(),
			Merge:    M("a", 1, "b", 2),
			Expected: M("a", 1, "b", 2),
		},
		{
			Name:     "Insert with all existing keys",
			Mode:     merge.MergeModeInsert,
			Original: M("a", 1, "b", 2, "c", 3),
			Merge:    M("a", 10, "b", 20),
			Expected: M("a", 1, "b", 2, "c", 3),
		},
		{
			Name:     "Insert nested maps recursively",
			Mode:     merge.MergeModeInsert,
			Original: M("user", M("name", "John", "email", "john@example.com")),
			Merge:    M("user", M("email", "new@example.com", "phone", "555-1234")),
			Expected: M("user", M("name", "John", "email", "john@example.com", "phone", "555-1234")),
		},
		{
			Name:     "Insert with deep nesting",
			Mode:     merge.MergeModeInsert,
			Original: M("level1", M("level2", M("value", "original"))),
			Merge:    M("level1", M("level2", M("value", "new", "extra", "added"))),
			Expected: M("level1", M("level2", M("value", "original", "extra", "added"))),
		},
		{
			Name:     "Insert preserves original non-zero primitives",
			Mode:     merge.MergeModeInsert,
			Original: M("str", "original", "num", 42, "bool", true),
			Merge:    M("str", "new", "num", 100, "bool", false),
			Expected: M("str", "original", "num", 42, "bool", true),
		},
	}

	TableTest(t, cases)
}

func TestInsertMode_Arrays(t *testing.T) {
	cases := []TestCase{
		{
			Name:     "Insert appends to array",
			Mode:     merge.MergeModeInsert,
			Original: []any{1, 2},
			Merge:    []any{3, 4},
			Expected: []any{1, 2, 3, 4},
		},
		{
			Name:     "Insert into empty array",
			Mode:     merge.MergeModeInsert,
			Original: []any{},
			Merge:    []any{1, 2, 3},
			Expected: []any{1, 2, 3},
		},
		{
			Name:     "Insert appends mixed types",
			Mode:     merge.MergeModeInsert,
			Original: []any{1, "two"},
			Merge:    []any{3.0, true, "five"},
			Expected: []any{1, "two", 3.0, true, "five"},
		},
		{
			Name:     "Insert appends complex objects",
			Mode:     merge.MergeModeInsert,
			Original: []any{M("id", 1)},
			Merge:    []any{M("id", 2), M("id", 3)},
			Expected: []any{M("id", 1), M("id", 3)},
		},
		{
			Name:     "Insert with nested arrays",
			Mode:     merge.MergeModeInsert,
			Original: []any{[]any{1, 2}},
			Merge:    []any{[]any{3, 4}},
			Expected: []any{[]any{1, 2, 3, 4}},
		},
	}

	TableTest(t, cases)
}

func TestInsertMode_SparseArrays(t *testing.T) {
	cases := []TestCase{
		{
			Name:     "Insert into sparse array merges maps recursively",
			Mode:     merge.MergeModeInsert,
			Original: []any{M("id", 1, "name", "Alice"), M("id", 2, "name", "Bob")},
			Merge:    map[int]any{0: M("email", "alice@example.com")},
			Expected: []any{M("id", 1, "name", "Alice", "email", "alice@example.com"), M("id", 2, "name", "Bob")},
		},
		{
			Name:     "Insert sparse array appends at non-existent index",
			Mode:     merge.MergeModeInsert,
			Original: []any{M("id", 1, "name", "Alice"), M("id", 2, "name", "Bob")},
			Merge:    map[int]any{2: M("id", 3, "name", "Charlie")},
			Expected: []any{M("id", 1, "name", "Alice"), M("id", 2, "name", "Bob"), M("id", 3, "name", "Charlie")},
		},
		{
			Name:     "Insert sparse array replaces non-maps",
			Mode:     merge.MergeModeInsert,
			Original: []any{"a", "b", "c"},
			Merge:    map[int]any{1: "B"},
			Expected: []any{"a", "b", "c"},
		},
		{
			Name:     "Insert sparse array with multiple indices",
			Mode:     merge.MergeModeInsert,
			Original: []any{M("id", 1), M("id", 2), M("id", 3)},
			Merge:    map[int]any{0: M("extra", "A"), 2: M("extra", "C"), 3: M("id", 4)},
			Expected: []any{M("id", 1, "extra", "A"), M("id", 2), M("id", 3, "extra", "C"), M("id", 4)},
		},
		{
			Name:     "Insert sparse preserves existing map values",
			Mode:     merge.MergeModeInsert,
			Original: []any{M("id", 1, "name", "Alice", "email", "alice@old.com")},
			Merge:    map[int]any{0: M("name", "Bob", "email", "new@example.com", "phone", "555-1234")},
			Expected: []any{M("id", 1, "name", "Alice", "email", "alice@old.com", "phone", "555-1234")},
		},
		{
			Name:     "Insert sparse into empty array",
			Mode:     merge.MergeModeInsert,
			Original: []any{},
			Merge:    map[int]any{0: "first", 1: "second"},
			Expected: []any{"first", "second"},
		},
	}

	TableTest(t, cases)
}

func TestInsertMode_Primitives(t *testing.T) {
	cases := []TestCase{
		{
			Name:     "Insert keeps non-zero string",
			Mode:     merge.MergeModeInsert,
			Original: "original",
			Merge:    "new",
			Expected: "original",
		},
		{
			Name:     "Insert replaces empty string",
			Mode:     merge.MergeModeInsert,
			Original: "",
			Merge:    "new",
			Expected: "new",
		},
		{
			Name:     "Insert keeps non-zero number",
			Mode:     merge.MergeModeInsert,
			Original: 42,
			Merge:    100,
			Expected: 42,
		},
		{
			Name:     "Insert replaces zero number",
			Mode:     merge.MergeModeInsert,
			Original: 0,
			Merge:    100,
			Expected: 100,
		},
		{
			Name:     "Insert keeps true boolean",
			Mode:     merge.MergeModeInsert,
			Original: true,
			Merge:    false,
			Expected: true,
		},
		{
			Name:     "Insert replaces false boolean",
			Mode:     merge.MergeModeInsert,
			Original: false,
			Merge:    true,
			Expected: true,
		},
		{
			Name:     "Insert replaces nil",
			Mode:     merge.MergeModeInsert,
			Original: nil,
			Merge:    "value",
			Expected: "value",
		},
	}

	TableTest(t, cases)
}

func TestInsertMode_ComplexScenarios(t *testing.T) {
	cases := []TestCase{
		{
			Name: "Insert complex nested structure",
			Mode: merge.MergeModeInsert,
			Original: M(
				"config", M(
					"database", M("host", "localhost", "port", 5432),
					"cache", M("enabled", true),
				),
				"features", []any{"auth", "logging"},
			),
			Merge: M(
				"config", M(
					"database", M("port", 3306, "user", "admin"),
					"cache", M("enabled", false, "ttl", 3600),
					"api", M("timeout", 30),
				),
				"features", []any{"metrics"},
				"version", "1.0.0",
			),
			Expected: M(
				"config", M(
					"database", M("host", "localhost", "port", 5432, "user", "admin"),
					"cache", M("enabled", true, "ttl", 3600),
					"api", M("timeout", 30),
				),
				"features", []any{"auth", "logging", "metrics"},
				"version", "1.0.0",
			),
		},
		{
			Name: "Insert with array of maps",
			Mode: merge.MergeModeInsert,
			Original: M(
				"users", []any{
					M("id", 1, "name", "Alice"),
					M("id", 2, "name", "Bob"),
				},
			),
			Merge: M(
				"users", []any{
					M("id", 1, "name", "Alice", "kurwa", "Bober"),
					M("id", 4),
					M("id", 3, "name", "Charlie"),
				},
			),
			Expected: M(
				"users", []any{
					M("id", 1, "name", "Alice", "kurwa", "Bober"),
					M("id", 2, "name", "Bob"),
					M("id", 3, "name", "Charlie"),
				},
			),
		},
	}

	TableTest(t, cases)
}

func TestInsertMode_EdgeCases(t *testing.T) {
	cases := []TestCase{
		{
			Name:     "Insert empty merge data into map",
			Mode:     merge.MergeModeInsert,
			Original: M("a", 1),
			Merge:    M(),
			Expected: M("a", 1),
		},
		{
			Name:     "Insert empty merge data into array",
			Mode:     merge.MergeModeInsert,
			Original: []any{1, 2},
			Merge:    []any{},
			Expected: []any{1, 2},
		},
		{
			Name:     "Insert with nil original map",
			Mode:     merge.MergeModeInsert,
			Original: map[string]any(nil),
			Merge:    M("a", 1),
			Expected: M("a", 1),
		},
		{
			Name:     "Insert with nil original array",
			Mode:     merge.MergeModeInsert,
			Original: []any(nil),
			Merge:    []any{1, 2},
			Expected: []any{1, 2},
		},
	}

	TableTest(t, cases)
}

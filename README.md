# Merge Package Documentation

## Overview

The `merge` package provides a flexible data merging system that recursively merges complex data structures with different merge strategies. It supports maps, arrays, and primitive values with five distinct merge modes.

May be useful for merging configurations or doing templates.

In plans is to extend library to support more different data types and be more flexible and cover functionality with tests.

## Merge Modes

### 1. ModeFullReplace (`"replace"`) and ModePartialReplace (`"replace_p"`)

Completely replaces the original data with the merge data.

**Behavior:**
- **Arrays:** Replaces entire array with merge data
- **Primitives:** Replaces with new value
- **Sparse Arrays:** Iterates through indices sequentially (0 to len-1); replaces values at existing indices; appends if index exceeds length; skips missing indices in sparse map

If merge data is nil -- will replace it anyway (with nil value)


**Behavior changes for replace_p (partial):**
- **Arrays:** Replaces elements by index; preserves original length if merge is shorter; extends if merge is longer

**Example:**
```go
orig := map[string]any{"a": 1, "b": 2}
merge := map[string]any{"b": 3, "c": 4}
// Result: {"a": 1, "b": 3, "c": 4}

origArray := []any{1, 2, 3}
mergeArray := []any{4, 5}
// Result: [4, 5]
```

**Example for replace_p:**
```go
origArray := []any{1, 2, 3, 4}
mergeArray := []any{10, 20}
// Result: [10, 20, 3, 4]

origArray := []any{1, 2}
mergeArray := []any{10, 20, 30}
// Result: [10, 20, 30]
```

### 3. ModeInsert (`"insert"`)

Inserts new values only; does not overwrite existing values.

**Behavior:**
- **Maps:** Adds new keys only; recursively merges nested maps; skips existing keys
- **Arrays:** Appends merge data to original array
- **Primitives:** Keeps original if exists
- **Sparse Arrays:** Merges into existing maps at specified indices (recursively); replaces non-maps; appends if index doesn't exist

**Example:**
```go
orig := map[string]any{"a": 1, "b": 2}
merge := map[string]any{"b": 3, "c": 4}
// Result: {"a": 1, "b": 2, "c": 4}

origArray := []any{1, 2}
mergeArray := []any{3, 4}
// Result: [1, 2, 3, 4]
```

### 4. ModeAppend (`"append"`)

Appends new data to existing data.

**Behavior:**
- **Maps:** Same as insert mode
- **Arrays:** Appends merge data to original array
- **Primitives:** Keeps original if exists
- **Sparse Arrays:** Appends all values from sparse map (ignores indices, treats as list of values)

**Example:**
```go
origArray := []any{1, 2, 3}
mergeArray := []any{4, 5}
// Result: [1, 2, 3, 4, 5]
```

### 5. ModeUpdate (`"update"`)

Updates existing values only; does not add new entries.

**Behavior:**
- **Maps:** Updates existing keys only; recursively merges nested maps; ignores new keys
- **Arrays:** Appends unique values only (uses deep equality check)
- **Primitives:** Replaces if original is non-nil
- **Sparse Arrays:** Updates values at specified indices only if index exists in original

**Example:**
```go
orig := map[string]any{"a": 1, "b": 2}
merge := map[string]any{"b": 3, "c": 4}
// Result: {"a": 1, "b": 3}

origArray := []any{1, 2, 3}
mergeArray := []any{2, 4}
// Result: [1, 2, 3, 4]  // 2 not added (duplicate), 4 added (unique)
```

## API Reference

### MergeData

```go
func MergeData(mode Mode, orig, mergeData any) (any, error)
```

Recursively merges `mergeData` into `orig` using the specified merge mode.

**Parameters:**
- `mode`: The merge strategy to use (see Merge Modes)
- `orig`: The original data structure
- `mergeData`: The data to merge into original

**Returns:**
- Merged data structure
- Error if mode is invalid or type mismatch occurs

**Supported Types:**
- `map[string]any`: Key-value maps
- `[]any`: Arrays/slices
- `map[int]any`: Sparse arrays (for array merging)
- Primitives: string, int, float, bool, etc.

### MergeMap

```go
var MergeMap = map[string]Mode{
    "replace":   ModeFullReplace,
    "replace_p": ModePartialReplace,
    "insert":    ModeInsert,
    "append":    ModeAppend,
    "update":    ModeUpdate,
}
```

Convenience map for converting string mode names to Mode constants.

## Usage Examples

### Basic Map Merging

```go
import "merge"

original := map[string]any{
    "name": "John",
    "age": 30,
    "city": "NYC",
}

updates := map[string]any{
    "age": 31,
    "country": "USA",
}

// Replace mode
result, err := merge.MergeData(merge.ModeFullReplace, original, updates)
// Result: {"name": "John", "age": 31, "city": "NYC", "country": "USA"}

// Update mode (only existing keys)
result, err = merge.MergeData(merge.ModeUpdate, original, updates)
// Result: {"name": "John", "age": 31, "city": "NYC"}

// Insert mode (only new keys)
result, err = merge.MergeData(merge.ModeInsert, original, updates)
// Result: {"name": "John", "age": 30, "city": "NYC", "country": "USA"}
```

### Nested Map Merging

```go
original := map[string]any{
    "user": map[string]any{
        "name": "John",
        "email": "john@example.com",
    },
    "settings": map[string]any{
        "theme": "dark",
    },
}

updates := map[string]any{
    "user": map[string]any{
        "email": "newemail@example.com",
        "phone": "555-1234",
    },
}

result, _ := merge.MergeData(merge.ModeFullReplace, original, updates)
// Result: Nested user map is merged recursively
```

### Array Merging

```go
original := []any{1, 2, 3}
additions := []any{4, 5}

// Append mode
result, _ := merge.MergeData(merge.ModeAppend, original, additions)
// Result: [1, 2, 3, 4, 5]

// Full replace mode
result, _ = merge.MergeData(merge.ModeFullReplace, original, additions)
// Result: [4, 5]

// Partial replace mode
result, _ = merge.MergeData(merge.ModePartialReplace, original, additions)
// Result: [4, 5, 3]
```

### Sparse Array Merging

```go
original := []any{"a", "b", "c", "d"}
sparseUpdates := map[int]any{
    1: "B",  // Update index 1
    3: "D",  // Update index 3
}

// Full replace mode
result, _ := merge.MergeData(merge.ModeFullReplace, original, sparseUpdates)
// Result: ["a", "B", "c", "D"]

// Insert mode - merges into existing maps
original := []any{
    map[string]any{"id": 1, "name": "Alice"},
    map[string]any{"id": 2, "name": "Bob"},
}
sparseInserts := map[int]any{
    0: map[string]any{"email": "alice@example.com"},
    2: map[string]any{"id": 3, "name": "Charlie"},
}
result, _ = merge.MergeData(merge.ModeInsert, original, sparseInserts)
// Result: [
//   {"id": 1, "name": "Alice", "email": "alice@example.com"},  // merged into existing map
//   {"id": 2, "name": "Bob"},
//   {"id": 3, "name": "Charlie"}  // appended (index 2 didn't exist)
// ]

// Append mode - appends values (ignores indices)
sparseAppends := map[int]any{
    5: "E",
    10: "F",
}
result, _ = merge.MergeData(merge.ModeAppend, original, sparseAppends)
// Result: ["a", "b", "c", "d", "E", "F"]
```

### String Mode Lookup

```go
modeStr := "replace_p"
mode := merge.MergeMap[modeStr]

result, err := merge.MergeData(mode, original, updates)
```

## Type Compatibility

### Valid Combinations

- `map[string]any` ↔ `map[string]any`
- `map[int]any` ↔ `map[int]any`
- `[]any` ↔ `[]any`
- `[]any` ↔ `map[int]any` (sparse array)
- Primitives ↔ Primitives (same or compatible types)

### Type Mismatches

The following combinations will return an error:
- `map[string]any` ↔ `[]any`
- `[]any` ↔ `map[string]any`
- Incompatible primitive types

## Zero Value Handling

The package includes special handling for zero values:

**Zero Values:**
- `nil`
- Empty strings `""`
- Empty arrays/slices `[]`
- Empty maps `map[]`
- Boolean `false`
- Numeric `0`, `0.0`

**Usage**

```go
myMergeMode := len(merge.Mergers)
merge.Mergers[myMergeMode] = MyMerger // anything that implements merger interface

// Warning! Doesn't guarantee safety for the orig data
result, err := merge.Data(mode, orig, mergeData)
if err != nil {
    log.Fatal(err)
}
// Use bulk merge for sequentually merging changes in different modes
merge.Bulk(
    orig,
    merge.ModeDataPair{
		Mode: merge.ModeAppend,
        Data: MDAppend
    },
    merge.ModeDataPair{
		Mode: merge.ModeUpdate,
        Data: MDUpd
    },
    merge.ModeDataPair{
		Mode: merge.ModeAppend,
        Data: MDAppend2
    },
    /* ... */
)
```

## Best Practices

1. **Choose the Right Mode:** Select the merge mode that matches your intent
   - Use `replace` for complete overwrites
   - Use `insert` when adding new data without changing existing
   - Use `update` when modifying existing data only
   - Use `append` for accumulating arrays
   - Use `replace_p` for partial array updates

2. **Handle Errors:** Always check the error return value

3. **Type Safety:** Ensure orig and mergeData have compatible types before merging

4. **Deep Nesting:** The merge is recursive, so it handles deeply nested structures automatically

5. **Immutability:** Note that maps are modified in place, but arrays are copied

## Performance Considerations

- Map merging modifies the original map in place
- Array merging creates new slices
- Deep equality checks in Update mode can be expensive for large arrays
- Recursive merging may have performance implications for deeply nested structures

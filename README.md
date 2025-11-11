# Merge Package Documentation

## Overview

The `merge` package provides a flexible data merging system that recursively merges complex data structures with different merge strategies. It supports maps, arrays, and primitive values with five distinct merge modes.

May be useful for merging configurations or doing templates.

In plans is to extend library to support more different data types and be more flexible and cover functionality with tests.

## Merge Modes

### 1. MergeModeFullReplace (`"replace"`)

Completely replaces the original data with the merge data.

**Behavior:**
- **Maps:** Recursively merges nested maps; replaces non-map values
- **Arrays:** Replaces entire array with merge data
- **Primitives:** Replaces with new value
- **Sparse Arrays:** Iterates through indices sequentially (0 to len-1); replaces values at existing indices; appends if index exceeds length; skips missing indices in sparse map

**Example:**
```go
orig := map[string]any{"a": 1, "b": 2}
merge := map[string]any{"b": 3, "c": 4}
// Result: {"a": 1, "b": 3, "c": 4}

origArray := []any{1, 2, 3}
mergeArray := []any{4, 5}
// Result: [4, 5]
```

### 2. MergeModePartialReplace (`"replace_p"`)

Similar to full replace, but preserves array length when possible.

**Behavior:**
- **Maps:** Same as full replace
- **Arrays:** Replaces elements by index; preserves original length if merge is shorter; extends if merge is longer
- **Primitives:** Replaces with new value
- **Sparse Arrays:** Replaces values at specified indices

**Example:**
```go
origArray := []any{1, 2, 3, 4}
mergeArray := []any{10, 20}
// Result: [10, 20, 3, 4]

origArray := []any{1, 2}
mergeArray := []any{10, 20, 30}
// Result: [10, 20, 30]
```

### 3. MergeModeInsert (`"insert"`)

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

### 4. MergeModeAppend (`"append"`)

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

### 5. MergeModeUpdate (`"update"`)

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
func MergeData(mode MergeMode, orig, mergeData any) (any, error)
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
var MergeMap = map[string]MergeMode{
    "replace":   MergeModeFullReplace,
    "replace_p": MergeModePartialReplace,
    "insert":    MergeModeInsert,
    "append":    MergeModeAppend,
    "update":    MergeModeUpdate,
}
```

Convenience map for converting string mode names to MergeMode constants.

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
result, err := merge.MergeData(merge.MergeModeFullReplace, original, updates)
// Result: {"name": "John", "age": 31, "city": "NYC", "country": "USA"}

// Update mode (only existing keys)
result, err = merge.MergeData(merge.MergeModeUpdate, original, updates)
// Result: {"name": "John", "age": 31, "city": "NYC"}

// Insert mode (only new keys)
result, err = merge.MergeData(merge.MergeModeInsert, original, updates)
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

result, _ := merge.MergeData(merge.MergeModeFullReplace, original, updates)
// Result: Nested user map is merged recursively
```

### Array Merging

```go
original := []any{1, 2, 3}
additions := []any{4, 5}

// Append mode
result, _ := merge.MergeData(merge.MergeModeAppend, original, additions)
// Result: [1, 2, 3, 4, 5]

// Full replace mode
result, _ = merge.MergeData(merge.MergeModeFullReplace, original, additions)
// Result: [4, 5]

// Partial replace mode
result, _ = merge.MergeData(merge.MergeModePartialReplace, original, additions)
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
result, _ := merge.MergeData(merge.MergeModeFullReplace, original, sparseUpdates)
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
result, _ = merge.MergeData(merge.MergeModeInsert, original, sparseInserts)
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
result, _ = merge.MergeData(merge.MergeModeAppend, original, sparseAppends)
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

**Behavior:**
In Insert and Append modes, if the original value is a zero value, it will be replaced with the merge value.

## Error Handling

The function returns errors in the following cases:

1. **Type Mismatch:** When orig and mergeData have incompatible types
   ```go
   // Error: type mismatch: map[string]any vs []interface {}
   ```

2. **Invalid Mode:** When using an invalid or unsupported merge mode
   ```go
   // Error: invalid mode for map: 99
   ```

Always check for errors when calling MergeData:

```go
result, err := merge.MergeData(mode, orig, mergeData)
if err != nil {
    log.Fatal(err)
}
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

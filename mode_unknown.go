package merge

import "fmt"

// UnknownMerger is a default merger that returns errors for all operations.
// Used for ModeUnknown to indicate that merging is not supported.
type UnknownMerger struct{}

// MergeMap returns an error indicating merging is not supported.
func (u *UnknownMerger) MergeMap(next Merger, path []string, orig, mergeData map[string]any) (map[string]any, error) {
	return nil, fmt.Errorf("merge not supported for unknown mode")
}

// MergeArray returns an error indicating merging is not supported.
func (u *UnknownMerger) MergeArray(next Merger, path []string, orig, mergeData []any) ([]any, error) {
	return nil, fmt.Errorf("merge not supported for unknown mode")
}

// MergeSparseArray returns an error indicating merging is not supported.
func (u *UnknownMerger) MergeSparseArray(next Merger, path []string, orig []any, mergeData map[int]any) ([]any, error) {
	return nil, fmt.Errorf("merge not supported for unknown mode")
}

// MergeIntMap returns an error indicating merging is not supported.
func (u *UnknownMerger) MergeIntMap(next Merger, path []string, orig, mergeData map[int]any) (map[int]any, error) {
	return nil, fmt.Errorf("merge not supported for unknown mode")
}

// MergePrimitive returns an error indicating merging is not supported.
func (u *UnknownMerger) MergePrimitive(next Merger, path []string, orig, mergeData any) (any, error) {
	return nil, fmt.Errorf("merge not supported for unknown mode")
}

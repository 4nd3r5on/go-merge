package merge

import (
	"fmt"
)

type MergeMode int

const (
	MergeModeFullReplace MergeMode = iota
	MergeModePartialReplace
	MergeModeInsert
	MergeModeAppend
	MergeModeUpdate

	DefaultMergeMode = MergeModeInsert
)

type Merger interface {
	MergeMap(next Merger, orig, mergeData map[string]any) (map[string]any, error)
	MergeArray(next Merger, orig, mergeData []any) ([]any, error)
	MergeSparseArray(next Merger, orig []any, mergeData map[int]any) ([]any, error)
	MergeIntMap(next Merger, orig, mergeData map[int]any) (map[int]any, error)
	MergePrimitive(next Merger, orig, mergeData any) (any, error)
}

var MergeModeMap = map[string]MergeMode{
	"replace":   MergeModeFullReplace,
	"replace_p": MergeModePartialReplace,
	"insert":    MergeModeInsert,
	"append":    MergeModeAppend,
	"update":    MergeModeUpdate,
}

var Mergers = map[MergeMode]Merger{
	MergeModeFullReplace:    &ReplaceMerger{Mode: MergeModeFullReplace, Conf: ReplaceMode{Partial: false}},
	MergeModePartialReplace: &ReplaceMerger{Mode: MergeModePartialReplace, Conf: ReplaceMode{Partial: true}},
	MergeModeInsert:         &InsertMerger{Mode: MergeModeInsert},
	MergeModeAppend:         &InsertMerger{Mode: MergeModeAppend, Conf: InsertMode{Append: true}},
	MergeModeUpdate:         &UpdateMerger{Mode: MergeModeUpdate},
}

// MergeData recursively merges `mergeData` into `orig`.
// Returns the resulting merged structure or an error on invalid mode or type mismatch.
func MergeData(mode MergeMode, orig, mergeData any) (any, error) {
	if merger, found := Mergers[mode]; found {
		return useMerger(merger, orig, mergeData)
	}
	return nil, fmt.Errorf("merger mode %q doesn't exist", mode)
}

func useMerger(m Merger, orig, mergeData any) (any, error) {
	switch o := orig.(type) {
	case map[string]any:
		md, ok := mergeData.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("type mismatch: map[string]any vs %T", mergeData)
		}
		return m.MergeMap(m, o, md)

	case []any:
		switch md := mergeData.(type) {
		case []any:
			return m.MergeArray(m, o, md)
		case map[int]any:
			return m.MergeSparseArray(m, o, md)
		default:
			return nil, fmt.Errorf("type mismatch: []any vs %T", mergeData)
		}

	case map[int]any:
		md, ok := mergeData.(map[int]any)
		if !ok {
			return nil, fmt.Errorf("type mismatch: map[int]any vs %T", mergeData)
		}
		return m.MergeIntMap(m, o, md)

	default:
		return m.MergePrimitive(m, orig, mergeData)
	}
}

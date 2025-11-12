package merge

import (
	"fmt"
)

type Mode int

const (
	ModeFullReplace Mode = iota
	ModePartialReplace
	ModeInsert
	ModeAppend
	ModeUpdate

	DefaultMergeMode = ModeInsert
)

type Merger interface {
	MergeMap(next Merger, orig, mergeData map[string]any) (map[string]any, error)
	MergeArray(next Merger, orig, mergeData []any) ([]any, error)
	MergeSparseArray(next Merger, orig []any, mergeData map[int]any) ([]any, error)
	MergeIntMap(next Merger, orig, mergeData map[int]any) (map[int]any, error)
	MergePrimitive(next Merger, orig, mergeData any) (any, error)
}

type ModeDataPair struct {
	Mode Mode
	Data any
}

var ModeMap = map[string]Mode{
	"replace":   ModeFullReplace,
	"replace_p": ModePartialReplace,
	"insert":    ModeInsert,
	"append":    ModeAppend,
	"update":    ModeUpdate,
}

var Mergers = map[Mode]Merger{
	ModeFullReplace:    &ReplaceMerger{Mode: ModeFullReplace, Conf: ReplaceMode{Partial: false}},
	ModePartialReplace: &ReplaceMerger{Mode: ModePartialReplace, Conf: ReplaceMode{Partial: true}},
	ModeInsert:         &InsertMerger{Mode: ModeInsert},
	ModeAppend:         &InsertMerger{Mode: ModeAppend, Conf: InsertMode{Append: true}},
	ModeUpdate:         &UpdateMerger{Mode: ModeUpdate},
}

func Bulk(orig any, mergeData ...ModeDataPair) (any, error) {
	var err error
	for _, mergePart := range mergeData {
		orig, err = Data(mergePart.Mode, orig, mergePart.Data)
		if err != nil {
			return nil, err
		}
	}
	return orig, nil
}

// Data recursively merges `mergeData` into `orig`.
// Returns the resulting merged structure or an error on invalid mode or type mismatch.
func Data(mode Mode, orig, mergeData any) (any, error) {
	if merger, found := Mergers[mode]; found {
		return UseMerger(merger, orig, mergeData)
	}
	return nil, fmt.Errorf("merger mode %q doesn't exist", mode)
}

func UseMerger(m Merger, orig, mergeData any) (any, error) {
	switch o := orig.(type) {
	case map[string]any:
		if mergeData == nil {
			return orig, nil
		}
		md, ok := mergeData.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("type mismatch: map[string]any vs %T", mergeData)
		}
		return m.MergeMap(m, o, md)

	case []any:
		if mergeData == nil {
			return orig, nil
		}
		switch md := mergeData.(type) {
		case []any:
			return m.MergeArray(m, o, md)
		case map[int]any:
			return m.MergeSparseArray(m, o, md)
		default:
			return nil, fmt.Errorf("type mismatch: []any vs %T", mergeData)
		}

	case map[int]any:
		if mergeData == nil {
			return orig, nil
		}
		md, ok := mergeData.(map[int]any)
		if !ok {
			return nil, fmt.Errorf("type mismatch: map[int]any vs %T", mergeData)
		}
		return m.MergeIntMap(m, o, md)

	default:
		return m.MergePrimitive(m, orig, mergeData)
	}
}

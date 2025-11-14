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
	DefaultMergersCount

	DefaultMergeMode = ModeInsert
)

type Merger interface {
	MergeMap(next Merger, path []string, orig, mergeData map[string]any) (map[string]any, error)
	MergeArray(next Merger, path []string, orig, mergeData []any) ([]any, error)
	MergeSparseArray(next Merger, path []string, orig []any, mergeData map[int]any) ([]any, error)
	MergeIntMap(next Merger, path []string, orig, mergeData map[int]any) (map[int]any, error)
	MergePrimitive(next Merger, path []string, orig, mergeData any) (any, error)
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
		return UseMerger(merger, nil, orig, mergeData)
	}
	return nil, fmt.Errorf("merger mode %q doesn't exist", mode)
}

func formatPath(p []string) any {
	if len(p) == 0 {
		return "root"
	}
	return p
}

func typeMismatch(p []string, expected string, got any) error {
	return fmt.Errorf("merge error at %v: expected %s, got %T",
		formatPath(p), expected, got)
}

func UseMerger(m Merger, path []string, orig, mergeData any) (any, error) {
	if path == nil {
		path = make([]string, 0)
	}

	switch o := orig.(type) {
	case map[string]any:
		if mergeData == nil {
			return orig, nil
		}
		md, ok := mergeData.(map[string]any)
		if !ok {
			return nil, typeMismatch(path, "map[string]any", mergeData)
		}
		res, err := m.MergeMap(m, path, o, md)
		if err != nil {
			return nil, fmt.Errorf("map merge failed at %v: %w", formatPath(path), err)
		}
		return res, nil

	case []any:
		if mergeData == nil {
			return orig, nil
		}

		switch md := mergeData.(type) {
		case []any:
			res, err := m.MergeArray(m, path, o, md)
			if err != nil {
				return nil, fmt.Errorf("array merge failed at %v: %w", formatPath(path), err)
			}
			return res, nil

		case map[int]any:
			res, err := m.MergeSparseArray(m, path, o, md)
			if err != nil {
				return nil, fmt.Errorf("sparse array merge failed at %v: %w", formatPath(path), err)
			}
			return res, nil

		default:
			return nil, typeMismatch(path, "[]any or map[int]any", mergeData)
		}

	case map[int]any:
		if mergeData == nil {
			return orig, nil
		}
		md, ok := mergeData.(map[int]any)
		if !ok {
			return nil, typeMismatch(path, "map[int]any", mergeData)
		}
		res, err := m.MergeIntMap(m, path, o, md)
		if err != nil {
			return nil, fmt.Errorf("int map merge failed at %v: %w", formatPath(path), err)
		}
		return res, nil

	default:
		res, err := m.MergePrimitive(m, path, orig, mergeData)
		if err != nil {
			return nil, fmt.Errorf("primitive merge failed at %v: %w", formatPath(path), err)
		}
		return res, nil
	}
}

func init() {
	orig := map[string]any{}

	Bulk(orig,
		ModeDataPair{
			Mode: ModeAppend,
		},
		ModeDataPair{},
		ModeDataPair{},
	)
}

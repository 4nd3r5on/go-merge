package merge

import (
	"fmt"
	"reflect"
)

type MergeMode int8

const (
	MergeModeFullReplace MergeMode = iota
	MergeModePartialReplace
	MergeModeInsert
	MergeModeAppend
	MergeModeUpdate

	DefaultMergeMode = MergeModeInsert
)

var MergeModeMap = map[string]MergeMode{
	"replace":   MergeModeFullReplace,
	"replace_p": MergeModePartialReplace,
	"insert":    MergeModeInsert,
	"append":    MergeModeAppend,
	"update":    MergeModeUpdate,
}

// MergeData recursively merges `mergeData` into `orig`.
// Returns the resulting merged structure or an error on invalid mode or type mismatch.
func MergeData(mode MergeMode, orig, mergeData any) (any, error) {
	switch o := orig.(type) {
	case map[string]any:
		md, ok := mergeData.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("type mismatch: map[string]any vs %T", mergeData)
		}
		return mergeMap(mode, o, md)

	case []any:
		switch md := mergeData.(type) {
		case []any:
			return mergeArray(mode, o, md)
		case map[int]any:
			return mergeSparseArray(mode, o, md)
		default:
			return nil, fmt.Errorf("type mismatch: []any vs %T", mergeData)
		}

	default:
		return mergePrimitive(mode, o, mergeData)
	}
}

func mergeMap(mode MergeMode, orig, mergeData map[string]any) (map[string]any, error) {
	for k, v := range mergeData {
		old, exists := orig[k]

		switch mode {
		case MergeModeFullReplace, MergeModePartialReplace:
			if exists && isBothMap(old, v) {
				m, err := mergeMap(mode, old.(map[string]any), v.(map[string]any))
				if err != nil {
					return nil, err
				}
				orig[k] = m
			} else {
				orig[k] = v
			}

		case MergeModeInsert, MergeModeAppend:
			if !exists {
				orig[k] = v
			} else if isBothMap(old, v) {
				m, err := mergeMap(mode, old.(map[string]any), v.(map[string]any))
				if err != nil {
					return nil, err
				}
				orig[k] = m
			}

		case MergeModeUpdate:
			if exists {
				if isBothMap(old, v) {
					m, err := mergeMap(mode, old.(map[string]any), v.(map[string]any))
					if err != nil {
						return nil, err
					}
					orig[k] = m
				} else {
					orig[k] = v
				}
			}

		default:
			return nil, fmt.Errorf("invalid mode for map: %v", mode)
		}
	}
	return orig, nil
}

func mergeArray(mode MergeMode, orig, mergeData []any) ([]any, error) {
	switch mode {
	case MergeModeFullReplace:
		return append([]any(nil), mergeData...), nil

	case MergeModePartialReplace:
		out := make([]any, len(orig))
		copy(out, orig)
		for i := range mergeData {
			if i < len(out) {
				if isBothMap(out[i], mergeData[i]) {
					m, err := mergeMap(mode, out[i].(map[string]any), mergeData[i].(map[string]any))
					if err != nil {
						return nil, err
					}
					out[i] = m
				} else {
					out[i] = mergeData[i]
				}
			} else {
				out = append(out, mergeData[i])
			}
		}
		return out, nil

	case MergeModeInsert, MergeModeAppend:
		return append(orig, mergeData...), nil

	case MergeModeUpdate:
		out := make([]any, len(orig))
		copy(out, orig)
		for _, v := range mergeData {
			if !contains(out, v) {
				out = append(out, v)
			}
		}
		return out, nil

	default:
		return nil, fmt.Errorf("invalid mode for array: %v", mode)
	}
}

func mergeSparseArray(mode MergeMode, orig []any, mergeData map[int]any) ([]any, error) {
	out := make([]any, len(orig))
	copy(out, orig)

	switch mode {
	case MergeModeFullReplace, MergeModePartialReplace:
		for i := 0; i < len(mergeData); i++ {
			v, ok := mergeData[i]
			if !ok {
				continue
			}
			if i < len(out) {
				if isBothMap(out[i], v) {
					m, err := mergeMap(mode, out[i].(map[string]any), v.(map[string]any))
					if err != nil {
						return nil, err
					}
					out[i] = m
				} else {
					out[i] = v
				}
			} else {
				out = append(out, v)
			}
		}

	case MergeModeInsert:
		// INSERT: merge into existing map at index if possible,
		// or append new map if index doesn't exist
		for i, v := range mergeData {
			if i < len(out) {
				if isBothMap(out[i], v) {
					m, err := mergeMap(MergeModeInsert, out[i].(map[string]any), v.(map[string]any))
					if err != nil {
						return nil, err
					}
					out[i] = m
				} else {
					out[i] = v // replace non-map with new map
				}
			} else {
				out = append(out, v)
			}
		}

	case MergeModeAppend:
		// APPEND: just add new items (whole maps) at the end
		for _, v := range mergeData {
			out = append(out, v)
		}

	case MergeModeUpdate:
		for i, v := range mergeData {
			if i < len(out) {
				if isBothMap(out[i], v) {
					m, err := mergeMap(MergeModeUpdate, out[i].(map[string]any), v.(map[string]any))
					if err != nil {
						return nil, err
					}
					out[i] = m
				} else {
					out[i] = v
				}
			}
		}

	default:
		return nil, fmt.Errorf("invalid mode for sparse array: %v", mode)
	}

	return out, nil
}

func mergePrimitive(mode MergeMode, orig, mergeData any) (any, error) {
	switch mode {
	case MergeModeFullReplace, MergeModePartialReplace:
		return mergeData, nil
	case MergeModeInsert, MergeModeAppend:
		if isZeroValue(orig) {
			return mergeData, nil
		}
		return orig, nil
	case MergeModeUpdate:
		if orig != nil {
			return mergeData, nil
		}
		return orig, nil
	default:
		return nil, fmt.Errorf("invalid mode for primitive: %v", mode)
	}
}

// --- helpers ---

func isBothMap(a, b any) bool {
	_, okA := a.(map[string]any)
	_, okB := b.(map[string]any)
	return okA && okB
}

func contains(arr []any, val any) bool {
	for _, v := range arr {
		if reflect.DeepEqual(v, val) {
			return true
		}
	}
	return false
}

func isZeroValue(x any) bool {
	if x == nil {
		return true
	}
	v := reflect.ValueOf(x)
	switch v.Kind() {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	}
	return false
}

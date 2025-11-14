package merge

import "fmt"

type InsertMode struct {
	Append bool
}

type InsertMerger struct {
	Mode Mode
	Conf InsertMode
}

func (m *InsertMerger) MergeMap(next Merger, path []string, orig, mergeData map[string]any) (map[string]any, error) {
	if orig == nil {
		return mergeData, nil
	}

	for k, v := range mergeData {
		old, exists := orig[k]
		if !exists {
			orig[k] = v
			continue
		}

		// mutate path
		path = append(path, k)

		merged, err := UseMerger(next, path, old, v)

		// restore
		path = path[:len(path)-1]

		if err != nil {
			return nil, err
		}

		orig[k] = merged
	}
	return orig, nil
}

func (m *InsertMerger) MergeArray(next Merger, path []string, orig, mergeData []any) ([]any, error) {
	out := make([]any, len(orig))
	copy(out, orig)

	if m.Conf.Append {
		return append(orig, mergeData...), nil
	}

	for i := range mergeData {
		if i >= len(orig) {
			out = append(out, mergeData[i:]...)
			break
		}

		switch orig[i].(type) {
		case map[string]any, map[int]any, []any:
			path = append(path, fmt.Sprintf("%v", i))

			merged, err := UseMerger(next, path, orig[i], mergeData[i])

			path = path[:len(path)-1]

			if err != nil {
				return nil, err
			}
			out[i] = merged

		default:
			out = append(out, mergeData[i])
		}
	}
	return out, nil
}

func (m *InsertMerger) MergeSparseArray(next Merger, path []string, orig []any, mergeData map[int]any) ([]any, error) {
	out := make([]any, len(orig))
	copy(out, orig)

	if m.Conf.Append {
		return append(orig, sparseArrayToArray(mergeData)...), nil
	}

	leftToMerge := make(map[int]any, len(mergeData))

	for i, v := range mergeData {
		if i < len(out) {
			path = append(path, fmt.Sprintf("%v", i))

			merged, err := UseMerger(next, path, out[i], v)

			path = path[:len(path)-1]

			if err != nil {
				return nil, err
			}

			out[i] = merged
		} else {
			leftToMerge[i] = v
		}
	}

	if len(leftToMerge) == 0 {
		return out, nil
	}

	return append(out, sparseArrayToArray(leftToMerge)...), nil
}

func (m *InsertMerger) MergeIntMap(next Merger, path []string, orig, mergeData map[int]any) (map[int]any, error) {
	for k, v := range mergeData {
		old, exists := orig[k]
		if !exists {
			orig[k] = v
			continue
		}

		path = append(path, fmt.Sprintf("%v", k))

		merged, err := UseMerger(next, path, old, v)

		path = path[:len(path)-1]

		if err != nil {
			return nil, err
		}

		orig[k] = merged
	}
	return orig, nil
}

func (m *InsertMerger) MergePrimitive(_ Merger, _ []string, orig, mergeData any) (any, error) {
	return nonZero(orig, mergeData), nil
}

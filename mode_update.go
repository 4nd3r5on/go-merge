package merge

import "fmt"

type UpdateMerger struct{ Mode Mode }

func updateMergeMap[K comparable](
	next Merger,
	path []string,
	orig, mergeData map[K]any,
) (map[K]any, error) {

	for k, v := range mergeData {
		old, exists := orig[k]
		if !exists {
			continue
		}

		// mutate path in-place
		seg := fmt.Sprintf("%v", k)
		path = append(path, seg)

		merged, err := UseMerger(next, path, old, v)

		// restore path
		path = path[:len(path)-1]

		if err != nil {
			return nil, err
		}

		orig[k] = merged
	}
	return orig, nil
}

func (m *UpdateMerger) MergeArray(_ Merger, _ []string, orig, mergeData []any) ([]any, error) {
	out := make([]any, len(orig))
	copy(out, orig)
	for _, v := range mergeData {
		if !contains(out, v) {
			out = append(out, v)
		}
	}
	return out, nil
}

func (m *UpdateMerger) MergeSparseArray(next Merger, path []string, orig []any, mergeData map[int]any) ([]any, error) {
	out := make([]any, len(orig))
	copy(out, orig)

	for i, v := range mergeData {
		if i >= len(out) {
			continue
		}

		old := out[i]

		// path mutation
		path = append(path, fmt.Sprintf("%v", i))

		merged, err := UseMerger(next, path, old, v)

		// restore
		path = path[:len(path)-1]

		if err != nil {
			return nil, err
		}

		out[i] = merged
	}

	return out, nil
}

func (m *UpdateMerger) MergePrimitive(_ Merger, _ []string, orig, mergeData any) (any, error) {
	if isZeroValue(orig) {
		return orig, nil
	}
	return mergeData, nil
}

func (m *UpdateMerger) MergeMap(next Merger, path []string, orig, mergeData map[string]any) (map[string]any, error) {
	return updateMergeMap(next, path, orig, mergeData)
}

func (m *UpdateMerger) MergeIntMap(next Merger, path []string, orig, mergeData map[int]any) (map[int]any, error) {
	return updateMergeMap(next, path, orig, mergeData)
}

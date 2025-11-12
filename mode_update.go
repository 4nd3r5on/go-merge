package merge

type UpdateMerger struct{ Mode Mode }

func updateMergeMap[K comparable](
	next Merger,
	orig, mergeData map[K]any,
) (map[K]any, error) {
	for k, v := range mergeData {
		old, exists := orig[k]
		if !exists {
			continue
		}
		merged, err := UseMerger(next, old, v)
		if err != nil {
			return nil, err
		}
		orig[k] = merged
	}
	return orig, nil
}

func (m *UpdateMerger) MergeArray(_ Merger, orig, mergeData []any) ([]any, error) {
	out := make([]any, len(orig))
	copy(out, orig)
	for _, v := range mergeData {
		if !contains(out, v) {
			out = append(out, v)
		}
	}
	return out, nil
}

func (m *UpdateMerger) MergeSparseArray(next Merger, orig []any, mergeData map[int]any) ([]any, error) {
	out := make([]any, len(orig))
	copy(out, orig)

	for i, v := range mergeData {
		if i >= len(out) {
			continue
		}

		old := out[i]
		merged, err := UseMerger(next, old, v)
		if err != nil {
			return nil, err
		}
		out[i] = merged
	}

	return out, nil
}

func (m *UpdateMerger) MergePrimitive(_ Merger, orig, mergeData any) (any, error) {
	if isZeroValue(orig) {
		return orig, nil
	}
	return mergeData, nil
}

// refactored methods using the generic function
func (m *UpdateMerger) MergeMap(next Merger, orig, mergeData map[string]any) (map[string]any, error) {
	return updateMergeMap(next, orig, mergeData)
}

func (m *UpdateMerger) MergeIntMap(next Merger, orig, mergeData map[int]any) (map[int]any, error) {
	return updateMergeMap(next, orig, mergeData)
}

package merge

type InsertMode struct {
	Append bool
}

type InsertMerger struct {
	Mode MergeMode
	Conf InsertMode
}

func (m *InsertMerger) MergeMap(next Merger, orig, mergeData map[string]any) (map[string]any, error) {
	if orig == nil {
		return mergeData, nil
	}

	for k, v := range mergeData {
		old, exists := orig[k]

		if !exists {
			orig[k] = v
			continue
		}

		merged, err := useMerger(next, old, v)
		if err != nil {
			return nil, err
		}
		orig[k] = merged
	}
	return orig, nil
}

func (m *InsertMerger) MergeArray(next Merger, orig, mergeData []any) ([]any, error) {
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
			merged, err := useMerger(next, orig[i], mergeData[i])
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

func (m *InsertMerger) MergeSparseArray(next Merger, orig []any, mergeData map[int]any) ([]any, error) {
	out := make([]any, len(orig))
	copy(out, orig)

	if m.Conf.Append {
		return append(orig, sparseArrayToArray(mergeData)...), nil
	}

	leftToMerge := make(map[int]any, len(mergeData))

	for i, v := range mergeData {
		if i < len(out) {
			merged, err := useMerger(next, out[i], v)
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

func (m *InsertMerger) MergeIntMap(next Merger, orig, mergeData map[int]any) (map[int]any, error) {
	for k, v := range mergeData {
		old, exists := orig[k]
		if !exists {
			orig[k] = v
			continue
		}

		merged, err := useMerger(next, old, v)
		if err != nil {
			return nil, err
		}
		orig[k] = merged
	}
	return orig, nil
}

func (m *InsertMerger) MergePrimitive(_ Merger, orig, mergeData any) (any, error) {
	return nonZero(orig, mergeData), nil
}

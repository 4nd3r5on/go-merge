package merge

type ReplaceMode struct {
	Partial bool
}

type ReplaceMerger struct {
	Mode MergeMode
	Conf ReplaceMode
}

// mergeMapCore merges maps based on ReplaceMode.
func (m *ReplaceMerger) mergeMapCore(next Merger, orig, mergeData map[string]any, conf ReplaceMode) (map[string]any, error) {
	for k, v := range mergeData {
		old, exists := orig[k]
		if conf.Partial && !exists {
			continue
		}

		if exists {
			merged, err := useMerger(next, old, v)
			if err != nil {
				return nil, err
			}
			orig[k] = merged
		} else {
			orig[k] = v
		}
	}
	return orig, nil
}

// mergeArrayCore merges slices by index with optional partial mode.
func (m *ReplaceMerger) mergeArrayCore(next Merger, orig, mergeData []any, conf ReplaceMode) ([]any, error) {
	out := make([]any, len(orig))
	copy(out, orig)

	limit := len(mergeData)
	if conf.Partial && limit > len(out) {
		limit = len(out)
	}

	for i := range limit {
		if i < len(out) {
			merged, err := useMerger(next, out[i], mergeData[i])
			if err != nil {
				return nil, err
			}
			out[i] = merged
		} else {
			out = append(out, mergeData[i])
		}
	}
	if !conf.Partial && len(mergeData) > len(out) {
		out = append(out, mergeData[len(out):]...)
	}
	return out, nil
}

// mergeSparseArrayCore merges sparse arrays up to given indices.
func (m *ReplaceMerger) mergeSparseArrayCore(next Merger, orig []any, mergeData map[int]any, conf ReplaceMode) ([]any, error) {
	if mergeData == nil {
		return orig, nil
	}
	out := make([]any, len(orig))
	copy(out, orig)

	maxIdx := -1
	for i := range mergeData {
		if i > maxIdx {
			maxIdx = i
		}
	}

	for i := 0; i <= maxIdx; i++ {
		v, ok := mergeData[i]
		if !ok {
			continue
		}
		if i < len(out) {
			if conf.Partial && i >= len(orig) {
				break
			}
			merged, err := useMerger(next, out[i], v)
			if err != nil {
				return nil, err
			}
			out[i] = merged
		} else if !conf.Partial {
			out = append(out, v)
		}
	}
	return out, nil
}

func (m *ReplaceMerger) mergeIntMapCore(next Merger, orig, mergeData map[int]any, conf ReplaceMode) (map[int]any, error) {
	for k, v := range mergeData {
		old, exists := orig[k]
		if conf.Partial && !exists {
			continue
		}
		if exists {
			merged, err := useMerger(next, old, v)
			if err != nil {
				return nil, err
			}
			orig[k] = merged
		} else {
			orig[k] = v
		}
	}
	return orig, nil
}

// ===== MERGER INTERFACE IMPLEMENTATION =====

func (m *ReplaceMerger) MergeMap(next Merger, orig, mergeData map[string]any) (map[string]any, error) {
	return m.mergeMapCore(next, orig, mergeData, m.Conf)
}

func (m *ReplaceMerger) MergeArray(next Merger, orig, mergeData []any) ([]any, error) {
	return m.mergeArrayCore(next, orig, mergeData, m.Conf)
}

func (m *ReplaceMerger) MergeSparseArray(next Merger, orig []any, mergeData map[int]any) ([]any, error) {
	return m.mergeSparseArrayCore(next, orig, mergeData, m.Conf)
}

func (m *ReplaceMerger) MergeIntMap(next Merger, orig, mergeData map[int]any) (map[int]any, error) {
	return m.mergeIntMapCore(next, orig, mergeData, m.Conf)
}

func (m *ReplaceMerger) MergePrimitive(_ Merger, _, mergeData any) (any, error) {
	return mergeData, nil
}

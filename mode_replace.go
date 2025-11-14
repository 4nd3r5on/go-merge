package merge

import "fmt"

type ReplaceMode struct {
	Partial bool
}

type ReplaceMerger struct {
	Mode Mode
	Conf ReplaceMode
}

func replaceMergeMap[K comparable](
	next Merger,
	path []string,
	orig, mergeData map[K]any,
	conf ReplaceMode,
) (map[K]any, error) {

	for k, v := range mergeData {
		old, exists := orig[k]
		if conf.Partial && !exists {
			continue
		}

		if exists {
			path = append(path, fmt.Sprintf("%v", k))

			merged, err := UseMerger(next, path, old, v)

			path = path[:len(path)-1]

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

func (m *ReplaceMerger) MergeArray(next Merger, path []string, orig, mergeData []any) ([]any, error) {
	out := make([]any, len(orig))
	copy(out, orig)

	limit := len(mergeData)
	if m.Conf.Partial && limit > len(out) {
		limit = len(out)
	}

	for i := 0; i < limit; i++ {
		if i < len(out) {
			path = append(path, fmt.Sprintf("%v", i))

			merged, err := UseMerger(next, path, out[i], mergeData[i])

			path = path[:len(path)-1]

			if err != nil {
				return nil, err
			}

			out[i] = merged
		} else {
			out = append(out, mergeData[i])
		}
	}

	if !m.Conf.Partial && len(mergeData) > len(out) {
		out = append(out, mergeData[len(out):]...)
	}

	return out, nil
}

func (m *ReplaceMerger) MergeSparseArray(next Merger, path []string, orig []any, mergeData map[int]any) ([]any, error) {
	if mergeData == nil {
		return nil, nil
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
			if m.Conf.Partial && i >= len(orig) {
				break
			}

			path = append(path, fmt.Sprintf("%v", i))

			merged, err := UseMerger(next, path, out[i], v)

			path = path[:len(path)-1]

			if err != nil {
				return nil, err
			}

			out[i] = merged

		} else if !m.Conf.Partial {
			out = append(out, v)
		}
	}

	return out, nil
}

func (m *ReplaceMerger) MergeMap(
	next Merger,
	path []string,
	orig, mergeData map[string]any,
) (map[string]any, error) {
	return replaceMergeMap(next, path, orig, mergeData, m.Conf)
}

func (m *ReplaceMerger) MergeIntMap(
	next Merger,
	path []string,
	orig, mergeData map[int]any,
) (map[int]any, error) {
	return replaceMergeMap(next, path, orig, mergeData, m.Conf)
}

func (m *ReplaceMerger) MergePrimitive(_ Merger, _ []string, _, mergeData any) (any, error) {
	return mergeData, nil
}

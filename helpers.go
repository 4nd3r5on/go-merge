package merge

import (
	"math"
	"reflect"
)

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

func nonZero(a, b any) any {
	if isZeroValue(a) {
		return b
	}
	return a
}

func sparseArrayToArray[T any](sparceArr map[int]T) []T {
	var (
		minK int = math.MaxInt32
		maxK int = 0
	)
	for k := range sparceArr {
		maxK = max(k, maxK)
		minK = min(k, minK)
	}
	outArr := make([]T, 0, len(sparceArr))
	for i := minK; i <= maxK; i++ {
		if v, exists := sparceArr[i]; exists {
			outArr = append(outArr, v)
		}
	}
	return outArr
}

func contains(arr []any, val any) bool {
	for _, v := range arr {
		if reflect.DeepEqual(v, val) {
			return true
		}
	}
	return false
}

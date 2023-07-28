package parser

import "log"

func TransformISliceToStrSlice(in []interface{}) []string {
	out := make([]string, len(in))
	for idx, val := range in {
		if str, ok := val.(string); ok {
			out[idx] = str
		} else {
			log.Fatalf("Element %v at index %d is not a string", val, idx)
		}
	}
	return out
}

func IsFloatInLimits(target, low, high float64) bool {
    return low < target  && target < high
}

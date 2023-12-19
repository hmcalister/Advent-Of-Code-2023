package lib

// General interface for a comparison function, either lessThanFunc or greaterThanFunc
type comparisonFunctionType func(a, b int) bool

func lessThanFunc(a, b int) bool {
	return a < b
}

func greaterThanFunc(a, b int) bool {
	return a > b
}

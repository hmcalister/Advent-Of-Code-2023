package lib

//go:generate stringer --type comparisonTypeEnum
type comparisonTypeEnum rune

const (
	LESS_THAN    comparisonTypeEnum = '<'
	GREATER_THAN comparisonTypeEnum = '>'
)

type comparisonFunctionType func(int, int) int

// Given the type of comparison, return the function used to map the rangeStart and rangeEnd appropriately.
//
// If a function requires, for example x LESS_THAN 2000, then we simply update the rangeStart and rangeEnd like:
//
// ```
// rangeStart = min(rangeStart, 2000)
// rangeEnd = min(rangeEnd, 2000)
// ```
//
// such that the range values can never be greater than 2000
func getUpdateFuncs(comparisonType comparisonTypeEnum) (comparisonFunctionType, comparisonFunctionType) {
	if comparisonType == LESS_THAN {
		return func(i1, i2 int) int { return min(i1, i2) }, func(i1, i2 int) int { return max(i1, i2) }
	} else {
		return func(i1, i2 int) int { return max(i1, i2+1) }, func(i1, i2 int) int { return min(i1, i2+1) }
	}
}

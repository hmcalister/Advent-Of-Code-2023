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
	min2Arg := func(i1, i2 int) int { return min(i1, i2) }
	max2Arg := func(i1, i2 int) int { return max(i1, i2) }

	if comparisonType == LESS_THAN {
		return min2Arg, max2Arg
	} else {
		return max2Arg, min2Arg
	}
}

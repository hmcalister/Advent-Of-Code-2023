// Code generated by "stringer --type comparisonTypeEnum"; DO NOT EDIT.

package lib

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[LESS_THAN-60]
	_ = x[GREATER_THAN-62]
}

const (
	_comparisonTypeEnum_name_0 = "LESS_THAN"
	_comparisonTypeEnum_name_1 = "GREATER_THAN"
)

func (i comparisonTypeEnum) String() string {
	switch {
	case i == 60:
		return _comparisonTypeEnum_name_0
	case i == 62:
		return _comparisonTypeEnum_name_1
	default:
		return "comparisonTypeEnum(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}

// Code generated by "stringer --type PartPropertyEnum"; DO NOT EDIT.

package lib

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ExtremelyCoolProperty-120]
	_ = x[MusicalProperty-109]
	_ = x[AerodynamicProperty-97]
	_ = x[ShinyProperty-115]
}

const (
	_PartPropertyEnum_name_0 = "AerodynamicProperty"
	_PartPropertyEnum_name_1 = "MusicalProperty"
	_PartPropertyEnum_name_2 = "ShinyProperty"
	_PartPropertyEnum_name_3 = "ExtremelyCoolProperty"
)

func (i partPropertyEnum) String() string {
	switch {
	case i == 97:
		return _PartPropertyEnum_name_0
	case i == 109:
		return _PartPropertyEnum_name_1
	case i == 115:
		return _PartPropertyEnum_name_2
	case i == 120:
		return _PartPropertyEnum_name_3
	default:
		return "PartPropertyEnum(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}

// Code generated by "stringer -type=DirectionEnum"; DO NOT EDIT.

package lib

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[DIRECTION_NORTH-0]
	_ = x[DIRECTION_EAST-1]
	_ = x[DIRECTION_SOUTH-2]
	_ = x[DIRECTION_WEST-3]
}

const _DirectionEnum_name = "DIRECTION_NORTHDIRECTION_EASTDIRECTION_SOUTHDIRECTION_WEST"

var _DirectionEnum_index = [...]uint8{0, 15, 29, 44, 58}

func (i DirectionEnum) String() string {
	if i < 0 || i >= DirectionEnum(len(_DirectionEnum_index)-1) {
		return "DirectionEnum(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _DirectionEnum_name[_DirectionEnum_index[i]:_DirectionEnum_index[i+1]]
}

// Code generated by "stringer -type PulseTypeEnum"; DO NOT EDIT.

package lib

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[NO_PULSE-0]
	_ = x[LOW_PULSE-1]
	_ = x[HIGH_PULSE-2]
}

const _PulseTypeEnum_name = "NO_PULSELOW_PULSEHIGH_PULSE"

var _PulseTypeEnum_index = [...]uint8{0, 8, 17, 27}

func (i PulseTypeEnum) String() string {
	if i < 0 || i >= PulseTypeEnum(len(_PulseTypeEnum_index)-1) {
		return "PulseTypeEnum(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _PulseTypeEnum_name[_PulseTypeEnum_index[i]:_PulseTypeEnum_index[i+1]]
}
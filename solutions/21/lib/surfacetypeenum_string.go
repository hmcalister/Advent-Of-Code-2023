// Code generated by "stringer -type=SurfaceTypeEnum"; DO NOT EDIT.

package lib

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[SURFACE_PLOT-0]
	_ = x[SURFACE_ROCK-1]
}

const _SurfaceTypeEnum_name = "SURFACE_PLOTSURFACE_ROCK"

var _SurfaceTypeEnum_index = [...]uint8{0, 12, 24}

func (i SurfaceTypeEnum) String() string {
	if i < 0 || i >= SurfaceTypeEnum(len(_SurfaceTypeEnum_index)-1) {
		return "SurfaceTypeEnum(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _SurfaceTypeEnum_name[_SurfaceTypeEnum_index[i]:_SurfaceTypeEnum_index[i+1]]
}
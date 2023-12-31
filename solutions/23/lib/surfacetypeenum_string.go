// Code generated by "stringer -type SurfaceTypeEnum"; DO NOT EDIT.

package lib

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[SURFACE_FOREST-0]
	_ = x[SURFACE_PATH-1]
	_ = x[SURFACE_SLOPE_UP-2]
	_ = x[SURFACE_SLOPE_RIGHT-3]
	_ = x[SURFACE_SLOPE_DOWN-4]
	_ = x[SURFACE_SLOPE_LEFT-5]
}

const _SurfaceTypeEnum_name = "SURFACE_FORESTSURFACE_PATHSURFACE_SLOPE_UPSURFACE_SLOPE_RIGHTSURFACE_SLOPE_DOWNSURFACE_SLOPE_LEFT"

var _SurfaceTypeEnum_index = [...]uint8{0, 14, 26, 42, 61, 79, 97}

func (i SurfaceTypeEnum) String() string {
	if i < 0 || i >= SurfaceTypeEnum(len(_SurfaceTypeEnum_index)-1) {
		return "SurfaceTypeEnum(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _SurfaceTypeEnum_name[_SurfaceTypeEnum_index[i]:_SurfaceTypeEnum_index[i+1]]
}

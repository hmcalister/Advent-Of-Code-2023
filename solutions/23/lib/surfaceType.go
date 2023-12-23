package lib

//go:generate stringer -type SurfaceTypeEnum
type SurfaceTypeEnum int

const (
	SURFACE_FOREST      SurfaceTypeEnum = iota
	SURFACE_PATH        SurfaceTypeEnum = iota
	SURFACE_SLOPE_UP    SurfaceTypeEnum = iota
	SURFACE_SLOPE_RIGHT SurfaceTypeEnum = iota
	SURFACE_SLOPE_DOWN  SurfaceTypeEnum = iota
	SURFACE_SLOPE_LEFT  SurfaceTypeEnum = iota
)

var (
	runeToSurfaceTypeMap map[rune]SurfaceTypeEnum = map[rune]SurfaceTypeEnum{
		'#': SURFACE_FOREST,
		'.': SURFACE_PATH,
		'^': SURFACE_SLOPE_UP,
		'>': SURFACE_SLOPE_RIGHT,
		'v': SURFACE_SLOPE_DOWN,
		'<': SURFACE_SLOPE_LEFT,
	}

	surfaceTypeToRuneMap map[SurfaceTypeEnum]rune = map[SurfaceTypeEnum]rune{
		SURFACE_FOREST:      '#',
		SURFACE_PATH:        '.',
		SURFACE_SLOPE_UP:    '^',
		SURFACE_SLOPE_RIGHT: '>',
		SURFACE_SLOPE_DOWN:  'v',
		SURFACE_SLOPE_LEFT:  '<',
	}
)

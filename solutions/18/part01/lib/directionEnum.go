package lib

//go:generate stringer -type=DirectionEnum
type DirectionEnum int

const (
	DIRECTION_UP    DirectionEnum = iota
	DIRECTION_RIGHT DirectionEnum = iota
	DIRECTION_DOWN  DirectionEnum = iota
	DIRECTION_LEFT  DirectionEnum = iota
	DIRECTION_NONE  DirectionEnum = iota
)

var (
	directionDecoderMap = map[string]DirectionEnum{
		"U": DIRECTION_UP,
		"R": DIRECTION_RIGHT,
		"D": DIRECTION_DOWN,
		"L": DIRECTION_LEFT,
	}
)

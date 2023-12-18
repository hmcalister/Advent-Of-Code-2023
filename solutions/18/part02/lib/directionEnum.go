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
		"3": DIRECTION_UP,
		"0": DIRECTION_RIGHT,
		"1": DIRECTION_DOWN,
		"2": DIRECTION_LEFT,
	}
)

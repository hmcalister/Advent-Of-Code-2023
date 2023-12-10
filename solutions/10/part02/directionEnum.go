package part02

//go:generate stringer -type=directionEnum
type directionEnum int

const (
	DIRECTION_NORTH directionEnum = iota
	DIRECTION_EAST  directionEnum = iota
	DIRECTION_SOUTH directionEnum = iota
	DIRECTION_WEST  directionEnum = iota
)

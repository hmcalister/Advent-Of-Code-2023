package lib

//go:generate stringer -type=DirectionEnum
type DirectionEnum int

const (
	DIRECTION_NORTH DirectionEnum = iota
	DIRECTION_EAST  DirectionEnum = iota
	DIRECTION_SOUTH DirectionEnum = iota
	DIRECTION_WEST  DirectionEnum = iota
)

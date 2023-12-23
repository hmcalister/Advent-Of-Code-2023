package lib

//go:generate stringer -type DirectionEnum
type DirectionEnum int

const (
	DIRECTION_UP    DirectionEnum = iota
	DIRECTION_RIGHT DirectionEnum = iota
	DIRECTION_DOWN  DirectionEnum = iota
	DIRECTION_LEFT  DirectionEnum = iota
)

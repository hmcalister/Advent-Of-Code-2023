package part02

//go:generate stringer -type=directionEnum
type directionEnum int

const (
	DIRECTION_LEFT  directionEnum = iota
	DIRECTION_RIGHT directionEnum = iota
)

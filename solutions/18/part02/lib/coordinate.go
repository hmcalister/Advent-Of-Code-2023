package lib

import "fmt"

type coordinate struct {
	X int
	Y int
}

func (c coordinate) String() string {
	return fmt.Sprintf("(%v, %v)", c.X, c.Y)
}

func (c coordinate) Move(direction DirectionEnum, distance int) coordinate {
	switch direction {
	case DIRECTION_UP:
		c.Y -= distance
	case DIRECTION_RIGHT:
		c.X += distance
	case DIRECTION_DOWN:
		c.Y += distance
	case DIRECTION_LEFT:
		c.X -= distance
	}

	return c
}

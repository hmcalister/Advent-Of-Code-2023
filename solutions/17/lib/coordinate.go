package lib

import "fmt"

type coordinate struct {
	X int
	Y int
}

func (c coordinate) String() string {
	return fmt.Sprintf("(%v, %v)", c.X, c.Y)
}

func (c coordinate) Move(direction DirectionEnum) coordinate {
	switch direction {
	case DIRECTION_UP:
		c.Y -= 1
	case DIRECTION_RIGHT:
		c.X += 1
	case DIRECTION_DOWN:
		c.Y += 1
	case DIRECTION_LEFT:
		c.X -= 1
	}

	return c
}

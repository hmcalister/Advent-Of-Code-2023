package lib

import "fmt"

type Coordinate struct {
	X int
	Y int
}

func (coord Coordinate) String() string {
	return fmt.Sprintf("(%v, %v)", coord.X, coord.Y)
}

func (coord Coordinate) Move(direction DirectionEnum) Coordinate {
	newCoord := coord
	switch direction {
	case DIRECTION_UP:
		newCoord.Y -= 1
	case DIRECTION_RIGHT:
		newCoord.X += 1
	case DIRECTION_DOWN:
		newCoord.Y += 1
	case DIRECTION_LEFT:
		newCoord.X -= 1
	}

	return newCoord
}

package lib

import "fmt"

type LightRay struct {
	Direction DirectionEnum
	XCoord    int
	YCoord    int
}

func (ray *LightRay) String() string {
	return fmt.Sprintf("%v (%v, %v)", ray.Direction.String(), ray.XCoord, ray.YCoord)
}

func (ray *LightRay) MarchRay() {
	switch ray.Direction {
	case DIRECTION_NORTH:
		ray.YCoord -= 1
	case DIRECTION_EAST:
		ray.XCoord += 1
	case DIRECTION_SOUTH:
		ray.YCoord += 1
	case DIRECTION_WEST:
		ray.XCoord -= 1
	}
}

package lib

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

type coordinate struct {
	X int
	Y int
}

func (c coordinate) String() string {
	return fmt.Sprintf("(%v, %v)", c.X, c.Y)
}

func (c coordinate) Move(direction DirectionEnum) coordinate {
	switch direction {
	case DIRECTION_NORTH:
		c.Y -= 1
	case DIRECTION_EAST:
		c.X += 1
	case DIRECTION_SOUTH:
		c.Y += 1
	case DIRECTION_WEST:
		c.X -= 1
	case DIRECTION_UNDEFINED:
		log.Fatal().Msg("cannot add undefined direction to coordinate")
	}

	return c
}

type pathFindNodeData struct {
	Coordinate                coordinate
	Direction                 DirectionEnum
	DirectionConsecutiveCount int
	Distance                  int
	HashString                string
}

func newPathFindNode(c coordinate, direction DirectionEnum, directionConsecutiveCount int, distance int) *pathFindNodeData {
	return &pathFindNodeData{
		Coordinate:                c,
		Direction:                 direction,
		DirectionConsecutiveCount: directionConsecutiveCount,
		Distance:                  distance,
		// HashString:                c.String(),
		HashString: fmt.Sprintf("%v %v %v", c.String(), direction.String(), directionConsecutiveCount),
	}
}

func (node *pathFindNodeData) String() string {
	return fmt.Sprintf("%v %v (%v-%v)", node.Coordinate.String(), node.Distance, node.Direction.String(), node.DirectionConsecutiveCount)
}

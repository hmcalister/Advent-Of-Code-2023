package lib

import "fmt"

type PathNodeData struct {
	currentCoordinate  Coordinate
	visitedCoordinates map[Coordinate]interface{}
}

func (node PathNodeData) String() string {
	return fmt.Sprintf("%v Len %v", node.currentCoordinate.String(), len(node.visitedCoordinates))
}

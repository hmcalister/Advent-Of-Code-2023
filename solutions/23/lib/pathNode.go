package lib

import "fmt"

const (
	visitedCoordinatePresenceIndicator = true
)

type PathNodeData struct {
	currentCoordinate  Coordinate
	visitedCoordinates map[Coordinate]interface{}
}

func (node PathNodeData) NextPathNode(direction DirectionEnum) PathNodeData {
func (node PathNodeData) String() string {
	return fmt.Sprintf("%v Len %v", node.CurrentCoordinate.String(), len(node.VisitedCoordinates))
}

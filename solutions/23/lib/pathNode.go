package lib

import "fmt"

type PathNodeData struct {
	CurrentCoordinate  Coordinate
	VisitedCoordinates map[Coordinate]interface{}
}

func (node PathNodeData) String() string {
	return fmt.Sprintf("%v Len %v", node.CurrentCoordinate.String(), len(node.VisitedCoordinates))
}

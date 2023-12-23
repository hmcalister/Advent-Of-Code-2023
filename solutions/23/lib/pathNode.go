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
	nextCoord := node.currentCoordinate.Move(direction)
	nextNode := PathNodeData{
		currentCoordinate:  nextCoord,
		visitedCoordinates: make(map[Coordinate]interface{}),
	}
	for k, v := range node.visitedCoordinates {
		nextNode.visitedCoordinates[k] = v
	}
	nextNode.visitedCoordinates[nextCoord] = visitedCoordinatePresenceIndicator

	return nextNode
}

func (node PathNodeData) String() string {
	return fmt.Sprintf("%v Len %v", node.currentCoordinate.String(), len(node.visitedCoordinates))
}

func (node PathNodeData) PathLength() int {
	return len(node.visitedCoordinates)
}


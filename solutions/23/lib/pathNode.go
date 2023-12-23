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

type PathNodePriorityQueue []PathNodeData

func (pq PathNodePriorityQueue) Len() int { return len(pq) }
func (pq PathNodePriorityQueue) Less(i, j int) bool {
	return len(pq[i].visitedCoordinates) < len(pq[j].visitedCoordinates)
}
func (pq PathNodePriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}
func (pq *PathNodePriorityQueue) Push(x any) {
	item := x.(PathNodeData)
	*pq = append(*pq, item)
}
func (pq *PathNodePriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

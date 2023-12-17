package lib

import (
	"fmt"
)

type pathFindNodeData struct {
	Coordinate coordinate
	Direction  DirectionEnum
	Streak     int
	Cost       int
}

func newPathFindNode(c coordinate, direction DirectionEnum, streak int, cost int) pathFindNodeData {
	return pathFindNodeData{
		Coordinate: c,
		Direction:  direction,
		Streak:     streak,
		Cost:       cost,
	}
}

func (node pathFindNodeData) String() string {
	return fmt.Sprintf("%v COST: %v (%v STREAK %v)", node.Coordinate.String(), node.Cost, node.Direction.String(), node.Streak)
}

func (node pathFindNodeData) Hash() string {
	return fmt.Sprintf("%v %v %v", node.Coordinate.String(), node.Direction.String(), node.Streak)
}

type PriorityQueue []pathFindNodeData

func (pq PriorityQueue) Len() int            { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool  { return pq[i].Cost < pq[j].Cost }
func (pq PriorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) { *pq = append(*pq, x.(pathFindNodeData)) }
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

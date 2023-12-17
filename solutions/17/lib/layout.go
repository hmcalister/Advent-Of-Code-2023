package lib

import (
	"bufio"
	"container/heap"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

type LayoutData struct {
	CostMap       [][]int
	MapWidth      int
	MapHeight     int
	minPathStreak int
	maxPathStreak int
}

func NewLayoutFromFileScanner(fileScanner *bufio.Scanner, minPathStreak int, maxPathStreak int) *LayoutData {
	costMap := make([][]int, 0)
	for fileScanner.Scan() {
		lineStr := fileScanner.Text()
		lineFields := strings.Split(lineStr, "")
		costLine := make([]int, len(lineFields))

		for index, field := range lineFields {
			parsedInt, ok := strconv.Atoi(field)
			if ok != nil {
				log.Fatal().Msgf("failed to parse field %v: %v in line %v", index, field, lineStr)
			}
			costLine[index] = parsedInt
		}
		log.Debug().Interface("CostMapLine", costLine).Send()
		costMap = append(costMap, costLine)
	}

	return &LayoutData{
		CostMap:       costMap,
		MapWidth:      len(costMap[0]),
		MapHeight:     len(costMap),
		minPathStreak: minPathStreak,
		maxPathStreak: maxPathStreak,
	}
}

func (layout *LayoutData) checkValidCoordinate(coord coordinate) bool {
	return (0 <= coord.X && coord.X < layout.MapWidth) && (0 <= coord.Y && coord.Y < layout.MapHeight)
}

func (layout *LayoutData) checkGoalNode(node pathFindNodeData) bool {
	return node.Coordinate.X == layout.MapWidth-1 && node.Coordinate.Y == layout.MapHeight-1
}

func (layout *LayoutData) PathFind() int {
	visited := make(map[string]pathFindNodeData)
	priorityQueue := make(PriorityQueue, 0)
	heap.Push(&priorityQueue, newPathFindNode(coordinate{1, 0}, DIRECTION_RIGHT, 0, layout.CostMap[0][1]))
	heap.Push(&priorityQueue, newPathFindNode(coordinate{0, 1}, DIRECTION_DOWN, 0, layout.CostMap[1][0]))

	for priorityQueue.Len() > 0 {
		currentNode := heap.Pop(&priorityQueue).(pathFindNodeData)
		log.Debug().Str("CurrentNode", currentNode.String()).Send()
		// If we are at the goal and have a valid streak length, we are done
		if layout.checkGoalNode(currentNode) && layout.minPathStreak <= currentNode.Streak {
			return currentNode.Cost
		}
		// If we have already seen this node, we can skip it
		if _, ok := visited[currentNode.Hash()]; ok {
			continue
		}
		visited[currentNode.Hash()] = currentNode

		// If we have not yet exhausted this streak, add the next node to the queue
		nextCoord := currentNode.Coordinate.Move(currentNode.Direction)
		if currentNode.Streak < layout.maxPathStreak-1 && layout.checkValidCoordinate(nextCoord) {
			nextCost := currentNode.Cost + layout.CostMap[nextCoord.Y][nextCoord.X]
			nextNode := newPathFindNode(nextCoord, currentNode.Direction, currentNode.Streak+1, nextCost)
			heap.Push(&priorityQueue, nextNode)
			log.Trace().Interface("ConsideringNode", nextNode).Send()
		}
		// If we have continued on this path long enough, consider turning as well
		if layout.minPathStreak <= currentNode.Streak {
			var turnDirections []DirectionEnum
			switch currentNode.Direction {
			case DIRECTION_UP:
				fallthrough
			case DIRECTION_DOWN:
				turnDirections = []DirectionEnum{DIRECTION_LEFT, DIRECTION_RIGHT}
			case DIRECTION_LEFT:
				fallthrough
			case DIRECTION_RIGHT:
				turnDirections = []DirectionEnum{DIRECTION_UP, DIRECTION_DOWN}
			}

			for _, nextDirection := range turnDirections {
				nextCoord := currentNode.Coordinate.Move(nextDirection)
				if !layout.checkValidCoordinate(nextCoord) {
					continue
				}
				nextCost := currentNode.Cost + layout.CostMap[nextCoord.Y][nextCoord.X]
				nextNode := newPathFindNode(nextCoord, nextDirection, 0, nextCost)
				heap.Push(&priorityQueue, nextNode)
				log.Trace().Interface("ConsideringNode", nextNode).Send()
			}
		}
	}

	return 0
}

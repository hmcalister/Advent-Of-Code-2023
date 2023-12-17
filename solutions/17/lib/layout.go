package lib

import (
	"bufio"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

type LayoutData struct {
	CostMap   [][]int
	MapWidth  int
	MapHeight int
}

func NewLayoutFromFileScanner(fileScanner *bufio.Scanner) *LayoutData {
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
		CostMap:   costMap,
		MapWidth:  len(costMap[0]),
		MapHeight: len(costMap),
	}
}

func (layout *LayoutData) getNeighbors(node *pathFindNodeData) []*pathFindNodeData {
	var nextCoordinate coordinate
	var nextDistance int
	var nextDirection DirectionEnum

	neighbors := make([]*pathFindNodeData, 0)

	currentCoordinate := node.Coordinate
	currentDirection := node.Direction
	currentDirectionCount := node.DirectionConsecutiveCount

	directionArr := []DirectionEnum{
		DIRECTION_NORTH,
		DIRECTION_EAST,
		DIRECTION_SOUTH,
		DIRECTION_WEST,
	}
	conditionalArr := []bool{
		currentCoordinate.Y > 0 && currentDirection != DIRECTION_SOUTH,
		currentCoordinate.X < layout.MapWidth-1 && currentDirection != DIRECTION_WEST,
		currentCoordinate.Y < layout.MapHeight-1 && currentDirection != DIRECTION_NORTH,
		currentCoordinate.X > 0 && currentDirection != DIRECTION_EAST,
	}

	for index := range directionArr {
		if !conditionalArr[index] {
			continue
		}
		nextDirection = directionArr[index]
		nextCoordinate = currentCoordinate.Move(nextDirection)
		nextDistance = node.Distance + layout.CostMap[nextCoordinate.Y][nextCoordinate.Y]

		if currentDirection == nextDirection {
			if currentDirectionCount < 3 {
				neighbors = append(neighbors, newPathFindNode(nextCoordinate, nextDirection, currentDirectionCount+1, nextDistance))
			}
		} else {
			neighbors = append(neighbors, newPathFindNode(nextCoordinate, nextDirection, 1, nextDistance))
		}
	}

	return neighbors
}

func (layout *LayoutData) reconstructPath(cameFromMap map[string]*pathFindNodeData, goalPosition *pathFindNodeData) {
	var ok bool
	var current *pathFindNodeData

	pathVisualization := make([][]rune, layout.MapHeight)
	for y := 0; y < layout.MapHeight; y += 1 {
		line := make([]rune, layout.MapWidth)
		for x := 0; x < layout.MapWidth; x += 1 {
			line[x] = '.'
		}
		pathVisualization[y] = line
	}

	ok = true
	current = goalPosition
	for ok {
		var directionRune rune
		switch current.Direction {
		case DIRECTION_NORTH:
			directionRune = '^'
		case DIRECTION_EAST:
			directionRune = '>'
		case DIRECTION_SOUTH:
			directionRune = 'v'
		case DIRECTION_WEST:
			directionRune = '<'
		case DIRECTION_UNDEFINED:
			directionRune = '#'
		}
		pathVisualization[current.Coordinate.Y][current.Coordinate.X] = directionRune
		log.Info().Str("Node", current.String()).Send()
		current, ok = cameFromMap[current.HashString]
	}

	for y := 0; y < layout.MapHeight; y += 1 {
		log.Info().Str("Path", string(pathVisualization[y])).Send()
	}
}

func (layout *LayoutData) checkGoalNode(node *pathFindNodeData) bool {
	return node.Coordinate.X == layout.MapWidth-1 && node.Coordinate.Y == layout.MapHeight-1
}

func (layout *LayoutData) DijkstraPathFind() *pathFindNodeData {
	var current *pathFindNodeData
	startPosition := newPathFindNode(coordinate{0, 0}, DIRECTION_UNDEFINED, 0, 0)

	// List of the open nodes
	openSet := make([]*pathFindNodeData, 0)
	openSet = append(openSet, startPosition)

	// Map of all visited nodes, keyed by the coordinate alone
	visitedNodes := make(map[string][]*pathFindNodeData)

	// Map of best distance to each position
	distanceMap := make(map[string]int)
	distanceMap[startPosition.HashString] = 0

	// Map of how we arrived at each node
	cameFromMap := make(map[string]*pathFindNodeData)

	for len(openSet) > 0 {
		sort.Slice(openSet, func(i, j int) bool {
			return openSet[i].Distance < openSet[j].Distance
		})
		current, openSet = openSet[0], openSet[1:]
		neighbors := layout.getNeighbors(current)
		log.Debug().Interface("CurrentNode", current).Send()

		visitedNodes[current.Coordinate.String()] = append(visitedNodes[current.Coordinate.String()], current)

		for _, neighbor := range neighbors {
			currentBestDistance, ok := distanceMap[neighbor.HashString]
			if !ok {
				currentBestDistance = math.MaxInt
			}

			if neighbor.Distance < currentBestDistance {
				log.Debug().Interface("ConsideringNode", neighbor).Send()
				distanceMap[neighbor.HashString] = neighbor.Distance
				openSet = append(openSet, neighbor)
				cameFromMap[neighbor.HashString] = current
			}
		}
	}

	var bestGoalNode *pathFindNodeData
	goalNodeCoordinateString := coordinate{layout.MapWidth - 1, layout.MapHeight - 1}.String()
	bestGoalNodeDistance := math.MaxInt
	log.Debug().Interface("GoalNodes", visitedNodes[goalNodeCoordinateString]).Send()
	for _, goalNode := range visitedNodes[goalNodeCoordinateString] {
		if goalNode.Distance < bestGoalNodeDistance {
			bestGoalNode = goalNode
			bestGoalNodeDistance = goalNode.Distance
		}
	}
	layout.reconstructPath(cameFromMap, bestGoalNode)

	return bestGoalNode
}

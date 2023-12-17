package lib

import (
	"bufio"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

type LayoutData struct {
	CostMap [][]int
	xLim    int
	yLim    int
}

func GetCostMapFromFileScanner(fileScanner *bufio.Scanner) [][]int {
	costMap := make([][]int, 0)
	for fileScanner.Scan() {
		line := fileScanner.Text()
		lineInts := strings.Split(line, "")
		newLine := make([]int, len(line))

		for index, i := range lineInts {
			parsedInt, err := strconv.Atoi(i)
			if err != nil {
				log.Fatal().Msgf("failed to parse integer %v: %v from line %v: %v", index, i, len(costMap), line)
			}
			newLine[index] = parsedInt
		}
		costMap = append(costMap, newLine)
	}

	return costMap
}

func NewLayout(costMap [][]int) *LayoutData {
	log.Debug().Str("TargetState", fmt.Sprintf("(%v, %v)", len(costMap[0])-1, len(costMap)-1)).Send()
	return &LayoutData{
		CostMap: costMap,
		xLim:    len(costMap[0]) - 1,
		yLim:    len(costMap) - 1,
	}
}

func (layout *LayoutData) heuristicFunction(XCoord, YCoord int) int {
	return (layout.yLim - YCoord) + (layout.xLim - XCoord)
}

func (layout *LayoutData) newPositionData(XCoord, YCoord, gScore int, direction DirectionEnum, directionCount int) *positionData {
	h := layout.heuristicFunction(XCoord, YCoord)
	return &positionData{
		Direction:        direction,
		DirectionCount:   directionCount,
		CoordinateString: fmt.Sprintf("(%v, %v)", XCoord, YCoord),
		XCoord:           XCoord,
		YCoord:           YCoord,
		GScore:           gScore,
		HScore:           h,
		FScore:           gScore + h,
	}
}

func (layout *LayoutData) checkGoalState(position *positionData) bool {
	return position.XCoord == layout.xLim && position.YCoord == layout.yLim
}

func (layout *LayoutData) getNeighbors(position *positionData) []*positionData {
	var nextX int
	var nextY int
	var nextGScore int
	var nextDirection DirectionEnum
	neighbors := make([]*positionData, 0)

	x, y := position.XCoord, position.YCoord
	currentDirection := position.Direction
	currentDirectionCount := position.DirectionCount

	if x > 0 && currentDirection != DIRECTION_EAST {
		nextX = x - 1
		nextY = y
		nextGScore = position.GScore + layout.CostMap[nextX][nextY]
		nextDirection = DIRECTION_WEST
		if currentDirection == nextDirection {
			if currentDirectionCount < 3 {
				neighbors = append(neighbors, layout.newPositionData(nextX, nextY, nextGScore, nextDirection, currentDirectionCount+1))
			}
		} else {
			neighbors = append(neighbors, layout.newPositionData(nextX, nextY, nextGScore, nextDirection, 1))
		}
	}
	if y > 0 && currentDirection != DIRECTION_SOUTH {
		nextX = x
		nextY = y - 1
		nextGScore := position.GScore + layout.CostMap[nextX][nextY]
		nextDirection = DIRECTION_NORTH
		if currentDirection == nextDirection {
			if currentDirectionCount < 3 {
				neighbors = append(neighbors, layout.newPositionData(nextX, nextY, nextGScore, nextDirection, currentDirectionCount+1))
			}
		} else {
			neighbors = append(neighbors, layout.newPositionData(nextX, nextY, nextGScore, nextDirection, 1))
		}
	}
	if x < layout.xLim && currentDirection != DIRECTION_WEST {
		nextX = x + 1
		nextY = y
		nextGScore := position.GScore + layout.CostMap[nextX][nextY]
		nextDirection = DIRECTION_EAST
		if currentDirection == nextDirection {
			if currentDirectionCount < 3 {
				neighbors = append(neighbors, layout.newPositionData(nextX, nextY, nextGScore, nextDirection, currentDirectionCount+1))
			}
		} else {
			neighbors = append(neighbors, layout.newPositionData(nextX, nextY, nextGScore, nextDirection, 1))
		}
	}
	if y < layout.yLim && currentDirection != DIRECTION_NORTH {
		nextX = x
		nextY = y + 1
		nextGScore := position.GScore + layout.CostMap[nextX][nextY]
		nextDirection = DIRECTION_SOUTH
		if currentDirection == nextDirection {
			if currentDirectionCount < 3 {
				neighbors = append(neighbors, layout.newPositionData(nextX, nextY, nextGScore, nextDirection, currentDirectionCount+1))
			}
		} else {
			neighbors = append(neighbors, layout.newPositionData(nextX, nextY, nextGScore, nextDirection, 1))
		}
	}

	return neighbors
}

func (layout *LayoutData) reconstructPath(cameFromMap map[string]*positionData, goalPosition *positionData) {
	var ok bool
	var current *positionData
	log.Info().Msg("PATH FOUND")

	ok = true
	current = goalPosition
	for ok {
		log.Info().Interface("Node", current).Send()
		current, ok = cameFromMap[current.CoordinateString]
	}
}

func (layout *LayoutData) AStarPathFind() *positionData {
	var current *positionData
	startPosition := layout.newPositionData(0, 0, 0, DIRECTION_WEST, 0)

	openList := make([]*positionData, 0)
	openList = append(openList, startPosition)

	openMap := make(map[string]*positionData)
	openMap[startPosition.CoordinateString] = startPosition

	// Map of best paths to each node
	gScoreMap := make(map[string]int)
	gScoreMap[startPosition.CoordinateString] = startPosition.GScore

	// Map of best estimates from positions to end
	fScoreMap := make(map[string]int)
	fScoreMap[startPosition.CoordinateString] = startPosition.FScore

	// Map of how we arrived at each node
	cameFromMap := make(map[string]*positionData)

	for len(openList) > 0 {
		sort.Slice(openList, func(i, j int) bool {
			return openList[i].FScore < openList[j].FScore
		})
		current, openList = openList[0], openList[1:]
		delete(openMap, current.CoordinateString)

		log.Debug().
			Int("OpenListLen", len(openList)).
			Interface("Current", current).
			Send()

		if layout.checkGoalState(current) {
			log.Info().Interface("GoalState", current).Send()
			layout.reconstructPath(cameFromMap, current)
			return current
		}

		neighbors := layout.getNeighbors(current)
		for _, neighbor := range neighbors {
			currentBestGScore, ok := gScoreMap[neighbor.CoordinateString]
			if !ok {
				// We have never seen this node before
				currentBestGScore = math.MaxInt
			}

			if neighbor.GScore < currentBestGScore {
				cameFromMap[neighbor.CoordinateString] = current
				gScoreMap[neighbor.CoordinateString] = neighbor.GScore
				fScoreMap[neighbor.CoordinateString] = neighbor.FScore
				if _, ok := openMap[neighbor.CoordinateString]; !ok {
					openList = append(openList, neighbor)
					openMap[neighbor.CoordinateString] = neighbor
				}
			}
		}
	}

	return nil
}

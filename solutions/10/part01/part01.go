package part01

import (
	"bufio"
	"strings"

	"github.com/rs/zerolog/log"
)

//go:generate stringer -type=directionEnum
type directionEnum int

const (
	START_RUNE rune = 'S'

	DIRECTION_NORTH directionEnum = iota
	DIRECTION_EAST  directionEnum = iota
	DIRECTION_SOUTH directionEnum = iota
	DIRECTION_WEST  directionEnum = iota
)

var (
	PipeMaze [][]rune

	NorthwardsRunes map[rune]directionEnum
	EastwardsRunes  map[rune]directionEnum
	SouthwardsRunes map[rune]directionEnum
	WestwardsRunes  map[rune]directionEnum
	DirectionMap    map[directionEnum]map[rune]directionEnum
)

type NodeData struct {
	XCoordinate int
	YCoordinate int
	NodeRune    rune
}

func (node NodeData) nextNode(direction directionEnum) NodeData {
	var nextXCoord int
	var nextYCoord int
	switch direction {
	case DIRECTION_NORTH:
		nextXCoord, nextYCoord = node.XCoordinate, node.YCoordinate-1
	case DIRECTION_EAST:
		nextXCoord, nextYCoord = node.XCoordinate+1, node.YCoordinate
	case DIRECTION_SOUTH:
		nextXCoord, nextYCoord = node.XCoordinate, node.YCoordinate+1
	case DIRECTION_WEST:
		nextXCoord, nextYCoord = node.XCoordinate-1, node.YCoordinate
	}

	nextNode := NodeData{
		XCoordinate: nextXCoord,
		YCoordinate: nextYCoord,
		NodeRune:    PipeMaze[nextYCoord][nextXCoord],
	}

	return nextNode
}

func determineStartDirection(startNode NodeData) directionEnum {
	var nextNode NodeData
	nextNode = startNode.nextNode(DIRECTION_NORTH)
	if strings.ContainsRune("|7F", nextNode.NodeRune) {
		return DIRECTION_NORTH
	}

	nextNode = startNode.nextNode(DIRECTION_EAST)
	if strings.ContainsRune("-J7", nextNode.NodeRune) {
		return DIRECTION_EAST
	}

	nextNode = startNode.nextNode(DIRECTION_SOUTH)
	if strings.ContainsRune("|LJ", nextNode.NodeRune) {
		return DIRECTION_SOUTH
	}

	return DIRECTION_WEST
}

type LoopData struct {
	StartXCoordinate int
	StartYCoordinate int
	LoopNodes        []NodeData
}

func ProcessInput(fileScanner *bufio.Scanner) (int, error) {
	NorthwardsRunes = map[rune]directionEnum{
		'|': DIRECTION_NORTH,
		'7': DIRECTION_WEST,
		'F': DIRECTION_EAST,
	}
	EastwardsRunes = map[rune]directionEnum{
		'-': DIRECTION_EAST,
		'7': DIRECTION_SOUTH,
		'J': DIRECTION_NORTH,
	}
	SouthwardsRunes = map[rune]directionEnum{
		'|': DIRECTION_SOUTH,
		'L': DIRECTION_EAST,
		'J': DIRECTION_WEST,
	}
	WestwardsRunes = map[rune]directionEnum{
		'-': DIRECTION_WEST,
		'L': DIRECTION_NORTH,
		'F': DIRECTION_SOUTH,
	}

	DirectionMap = map[directionEnum]map[rune]directionEnum{
		DIRECTION_NORTH: NorthwardsRunes,
		DIRECTION_EAST:  EastwardsRunes,
		DIRECTION_SOUTH: SouthwardsRunes,
		DIRECTION_WEST:  WestwardsRunes,
	}

	PipeMaze = make([][]rune, 0)
	for fileScanner.Scan() {
		line := fileScanner.Text()
		PipeMaze = append(PipeMaze, []rune(line))
		log.Debug().
			Int("YCoord", len(PipeMaze)).
			Str("Line", line).
			Send()
	}

	var startNode NodeData
	for yCoord, row := range PipeMaze {
		for xCoord, col := range row {
			if col == START_RUNE {
				startNode = NodeData{
					XCoordinate: xCoord,
					YCoordinate: yCoord,
					NodeRune:    START_RUNE,
				}
			}
		}
	}
	loop := LoopData{
		StartXCoordinate: startNode.XCoordinate,
		StartYCoordinate: startNode.YCoordinate,
		LoopNodes:        []NodeData{startNode},
	}

	direction := determineStartDirection(startNode)
	currentNode := startNode
	log.Debug().
		Interface("StartNode", currentNode).
		Str("StartNodeRune", string(currentNode.NodeRune)).
		Str("Direction", direction.String()).
		Send()

	for {
		currentNode = currentNode.nextNode(direction)
		direction = DirectionMap[direction][currentNode.NodeRune]
		log.Debug().
			Interface("CurrentNode", currentNode).
			Str("CurrentNodeRune", string(currentNode.NodeRune)).
			Str("Direction", direction.String()).
			Send()

		if currentNode.NodeRune == START_RUNE {
			break
		}

		loop.LoopNodes = append(loop.LoopNodes, currentNode)
	}

	return len(loop.LoopNodes) / 2, nil
}

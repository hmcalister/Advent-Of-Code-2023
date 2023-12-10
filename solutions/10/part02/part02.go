package part02

import (
	"bufio"

	"github.com/rs/zerolog/log"
)

const (
	LINE_LENGTH int  = 20
	START_RUNE  rune = 'S'
	GROUND_RUNE rune = '.'
)

var (
	PipeMaze                  [][]rune
	DirectionMap              map[directionEnum]map[rune]directionEnum
	LoopPipeLinearCoordinates map[int]NodeData
)

func ProcessInput(fileScanner *bufio.Scanner) (int, error) {
	DirectionMap = createDirectionMap()

	PipeMaze = make([][]rune, 0)
	LoopPipeLinearCoordinates = make(map[int]NodeData)

	partition := newPartitionData()
	for fileScanner.Scan() {
		line := fileScanner.Text()
		PipeMaze = append(PipeMaze, []rune(line))
		partition.PartitionIndicesArray = append(partition.PartitionIndicesArray, make([]int, len(line)))
		log.Trace().
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

	direction := determineStartDirection(startNode)
	currentNode := startNode
	partition.PartitionIndicesArray[currentNode.YCoordinate][currentNode.XCoordinate] = -1
	log.Debug().
		Interface("StartNode", currentNode).
		Str("StartNodeRune", string(currentNode.NodeRune)).
		Str("Direction", direction.String()).
		Send()

	loop := LoopData{
		StartXCoordinate: startNode.XCoordinate,
		StartYCoordinate: startNode.YCoordinate,
		LoopNodes:        []NodeData{startNode},
		LoopDirection:    []directionEnum{direction},
	}

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
		loop.LoopDirection = append(loop.LoopDirection, direction)
		partition.PartitionIndicesArray[currentNode.YCoordinate][currentNode.XCoordinate] = -1
	}

	log.Debug().Msg("finished parsing loop")
	log.Trace().
		Interface("LoopNodes", loop.LoopNodes).
		Interface("LoopDirections", loop.LoopDirection).
		Send()

	for xCoord := 0; xCoord < LINE_LENGTH; xCoord += 1 {
		partition.determinePartition(1, xCoord, 0)
		partition.determinePartition(1, xCoord, len(PipeMaze)-1)
	}
	for yCoord := 0; yCoord < len(PipeMaze); yCoord += 1 {
		partition.determinePartition(1, 0, yCoord)
		partition.determinePartition(1, LINE_LENGTH-1, yCoord)
	}

	enclosedNodeCount := 0
	for yCoord := 0; yCoord < len(PipeMaze); yCoord += 1 {
		for xCoord := 0; xCoord < LINE_LENGTH; xCoord += 1 {
			if partition.PartitionIndicesArray[yCoord][xCoord] == 0 {
				enclosedNodeCount += 1
			}
		}
	}

	return enclosedNodeCount, nil
}

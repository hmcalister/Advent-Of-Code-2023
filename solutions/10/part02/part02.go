package part02

import (
	"bufio"

	"github.com/rs/zerolog/log"
)

const (
	START_RUNE  rune = 'S'
	GROUND_RUNE rune = '.'
)

var (
	PipeMaze     [][]rune
	DirectionMap map[directionEnum]map[rune]directionEnum
)

func ProcessInput(fileScanner *bufio.Scanner) (int, error) {
	DirectionMap = createDirectionMap()

	PipeMaze = make([][]rune, 0)
	partition := newPartitionData()
	for fileScanner.Scan() {
		line := fileScanner.Text()
		PipeMaze = append(PipeMaze, []rune(line))
		partition.PartitionArray = append(partition.PartitionArray, make([]int, len(line)))
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
	partition.PartitionArray[currentNode.YCoordinate][currentNode.XCoordinate] = -1
	log.Debug().
		Interface("StartNode", currentNode).
		Str("StartNodeRune", string(currentNode.NodeRune)).
		Str("Direction", direction.String()).
		Send()

	for {
		currentNode = currentNode.nextNode(direction)
		direction = DirectionMap[direction][currentNode.NodeRune]
		partition.PartitionArray[currentNode.YCoordinate][currentNode.XCoordinate] = -1
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

	for yCoord, row := range PipeMaze {
		for xCoord := range row {
			if partition.PartitionArray[yCoord][xCoord] == 0 {
				partition.determinePartition(xCoord, yCoord)
			}
		}
	}

	log.Info().
		Interface("PartitionGroundCounts", partition.PartitionsGroundCount).
		Interface("PartitionSizes", partition.PartitionSizes).
		Send()

	return 0, nil
}

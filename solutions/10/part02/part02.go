package part02

import (
	"bufio"
	"slices"

	"github.com/rs/zerolog/log"
)

const (
	START_RUNE  rune = 'S'
	GROUND_RUNE rune = '.'
	LINE_LENGTH int  = 141
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

	for fileScanner.Scan() {
		line := fileScanner.Text()
		PipeMaze = append(PipeMaze, []rune(line))
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
	LoopPipeLinearCoordinates[currentNode.YCoordinate*LINE_LENGTH+currentNode.XCoordinate] = currentNode
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
		LoopPipeLinearCoordinates[currentNode.YCoordinate*LINE_LENGTH+currentNode.XCoordinate] = currentNode
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
	}

	log.Debug().Msg("finished parsing loop")
	log.Trace().
		Interface("LoopNodes", loop.LoopNodes).
		Interface("LoopDirections", loop.LoopDirection).
		Send()

	enclosedNodeCoordinates := make(map[int]NodeData)
	directionArr := []directionEnum{DIRECTION_NORTH, DIRECTION_EAST, DIRECTION_SOUTH, DIRECTION_WEST}
	for loopNodeIndex := range loop.LoopNodes {
		currentNode := loop.LoopNodes[loopNodeIndex]
		loopDirection := loop.LoopDirection[loopNodeIndex]
		probeDirection := directionArr[(slices.Index(directionArr, loopDirection)+3)%len(directionArr)]
		log.Debug().
			Int("LoopNodeIndex", loopNodeIndex).
			Interface("LoopNode", currentNode).
			Str("LoopDirection", loopDirection.String()).
			Str("ProbeDirection", probeDirection.String()).Send()

		probeNode := currentNode.nextNode(probeDirection)
		for {
			linearProbeNodeCoordinate := probeNode.YCoordinate*LINE_LENGTH + probeNode.XCoordinate
			_, probeNodeIsInLoop := LoopPipeLinearCoordinates[linearProbeNodeCoordinate]
			log.Debug().
				Interface("ProbeNode", probeNode).
				Str("ProbeDirection", probeDirection.String()).
				Bool("ProbeNodeIsInLoop", probeNodeIsInLoop).
				Send()
			if probeNodeIsInLoop {
				break
			}
			if (probeDirection == DIRECTION_NORTH && probeNode.YCoordinate == 0) ||
				(probeDirection == DIRECTION_EAST && probeNode.XCoordinate == LINE_LENGTH-1) ||
				(probeDirection == DIRECTION_SOUTH && probeNode.YCoordinate == len(PipeMaze)-1) ||
				(probeDirection == DIRECTION_WEST && probeNode.XCoordinate == 0) {
				log.Panic().Msgf("loop appears to have opposite handedness")
			}

			enclosedNodeCoordinates[linearProbeNodeCoordinate] = probeNode
			probeNode = probeNode.nextNode(probeDirection)
		}
	}

	return len(enclosedNodeCoordinates), nil
}

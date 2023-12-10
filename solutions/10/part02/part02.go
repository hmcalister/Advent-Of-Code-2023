package part02

import (
	"bufio"
	"regexp"

	"github.com/rs/zerolog/log"
)

const (
	LINE_LENGTH int  = 140
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
		LoopPipeLinearCoordinates[currentNode.YCoordinate*LINE_LENGTH+currentNode.XCoordinate] = currentNode
	}

	log.Debug().Msg("finished parsing loop")
	log.Trace().
		Interface("LoopNodes", loop.LoopNodes).
		Interface("LoopDirections", loop.LoopDirection).
		Send()

	for yCoord := 0; yCoord < len(PipeMaze); yCoord += 1 {
		for xCoord := 0; xCoord < LINE_LENGTH; xCoord += 1 {
			linearCoordinate := yCoord*LINE_LENGTH + xCoord
			if _, ok := LoopPipeLinearCoordinates[linearCoordinate]; !ok {
				PipeMaze[yCoord][xCoord] = GROUND_RUNE
			}
		}
	}

	enclosedNodeCount := 0
	loopCorner1 := regexp.MustCompile("F-*7")
	loopCorner2 := regexp.MustCompile("L-*J")
	loopBend1 := regexp.MustCompile("F-*J")
	loopBend2 := regexp.MustCompile("L-*7")
	for yCoord := 0; yCoord < len(PipeMaze); yCoord += 1 {
		originalLine := string(PipeMaze[yCoord])
		line := originalLine
		line = loopCorner1.ReplaceAllString(line, "")
		line = loopCorner2.ReplaceAllString(line, "")
		line = loopBend1.ReplaceAllString(line, "|")
		line = loopBend2.ReplaceAllString(line, "|")
		log.Debug().
			Int("YCoord", yCoord).
			Str("OriginalLine", originalLine).
			Str("EffectiveLine", line).
			Send()

		enclosedFlag := false
		for xCoord := 0; xCoord < len(line); xCoord += 1 {
			if line[xCoord] == '|' {
				enclosedFlag = !enclosedFlag
			}

			if line[xCoord] == '.' && enclosedFlag {
				enclosedNodeCount += 1
			}
		}
	}

	return enclosedNodeCount, nil
}

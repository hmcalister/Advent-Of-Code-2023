package part02

import (
	"bufio"
	"strings"

	"github.com/rs/zerolog/log"
)

const (
	START_LABEL_SUFFIX    = "A"
	TERMINAL_LABEL_SUFFIX = "Z"
)

type pathNodeData struct {
	Label          string
	LeftNodeLabel  string
	RightNodeLabel string
	IsTerminal     bool
}

type pathData struct {
	pathNodeMap       map[string]*pathNodeData
	allStartPathNodes []*pathNodeData
}

func createPath() *pathData {
	path := &pathData{}
	path.allStartPathNodes = make([]*pathNodeData, 0)
	path.pathNodeMap = make(map[string]*pathNodeData)

	return path
}

// Create and return a new path node data struct,
// as well as insert it into the pathNodeMap attribute
func (path *pathData) parseLineToPathNodeData(line string) {
	// A line is structured like the example: BKM = (CDC, PSH)
	// Labels are always exactly three letters long
	pathNodeLabel := line[0:3]

	// The left label exists from indices 7 to 10, and the right 12 to 15
	leftNodeLabel := line[7:10]
	rightNodeLabel := line[12:15]

	newPathNode := &pathNodeData{
		Label:          pathNodeLabel,
		LeftNodeLabel:  leftNodeLabel,
		RightNodeLabel: rightNodeLabel,
		IsTerminal:     strings.HasSuffix(pathNodeLabel, TERMINAL_LABEL_SUFFIX),
	}

	log.Trace().
		Str("ParsedLine", line).
		Interface("NewPathNode", newPathNode).
		Send()

	// Add the node to the node map and array
	if strings.HasSuffix(pathNodeLabel, START_LABEL_SUFFIX) {
		path.allStartPathNodes = append(path.allStartPathNodes, newPathNode)
	}
	path.pathNodeMap[pathNodeLabel] = newPathNode
}

func ProcessInput(fileScanner *bufio.Scanner) (int, error) {
	fileScanner.Scan()
	directionsLine := fileScanner.Text()
	log.Trace().Str("DirectionsLine", directionsLine).Send()

	directionsArray := make([]directionEnum, len(directionsLine))
	for i, r := range directionsLine {
		if r == 'L' {
			directionsArray[i] = DIRECTION_LEFT
		} else if r == 'R' {
			directionsArray[i] = DIRECTION_RIGHT
		} else {
			log.Panic().Msgf("encountered unknown rune while parsing directions line at index %v: %v", i, r)
		}
	}

	path := createPath()

	fileScanner.Scan()
	for fileScanner.Scan() {
		path.parseLineToPathNodeData(fileScanner.Text())
	}

	currentNodes := make([]*pathNodeData, len(path.allStartPathNodes))
	copy(currentNodes, path.allStartPathNodes)
	log.Debug().
		Int("NumCurrentNodes", len(currentNodes)).
		Send()

	numSteps := 0
	terminalNodeCount := 0
	var nextNode *pathNodeData
	for {
		terminalNodeCount = 0
		nextDirection := directionsArray[numSteps%len(directionsArray)]
		for currentNodeIndex, currentNode := range currentNodes {
			log.Trace().
				Interface("CurrentPathNode", currentNode).
				Int("NumSteps", numSteps).
				Str("NextDirection", nextDirection.String()).
				Send()

			if nextDirection == DIRECTION_LEFT {
				nextNode = path.pathNodeMap[currentNode.LeftNodeLabel]
			} else {
				nextNode = path.pathNodeMap[currentNode.RightNodeLabel]
			}
			currentNodes[currentNodeIndex] = nextNode

			if nextNode.IsTerminal {
				log.Trace().
					Interface("TerminalNodeFound", nextNode).
					Int("NumSteps", numSteps).
					Send()
				terminalNodeCount += 1
			}
		}
		if terminalNodeCount >= 2 {
			log.Debug().
				Int("NumSteps", numSteps).
				Int("TerminalNodeCount", terminalNodeCount).
				Send()
		}
		if terminalNodeCount == len(currentNodes) {
			break
		}

		numSteps += 1
	}

	return numSteps, nil
}

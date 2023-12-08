package part01

import (
	"bufio"

	"github.com/rs/zerolog/log"
)

const (
	START_LABEL    = "AAA"
	TERMINAL_LABEL = "ZZZ"
)

type pathNodeData struct {
	Label          string
	LeftNodeLabel  string
	RightNodeLabel string
	IsTerminal     bool
}

type pathData struct {
	pathNodeMap  map[string]*pathNodeData
	allPathNodes []*pathNodeData
}

func createPath() *pathData {
	path := &pathData{}
	path.allPathNodes = make([]*pathNodeData, 0)
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
		IsTerminal:     pathNodeLabel == TERMINAL_LABEL,
	}

	log.Trace().
		Str("ParsedLine", line).
		Interface("NewPathNode", newPathNode).
		Send()

	// Add the node to the node map and array
	path.allPathNodes = append(path.allPathNodes, newPathNode)
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

	currentPathNode := path.pathNodeMap[START_LABEL]
	numSteps := 0
	for {
		if currentPathNode.IsTerminal {
			log.Debug().
				Interface("TerminalNodeFound", currentPathNode).
				Int("NumSteps", numSteps).
				Send()
			break
		}
		nextDirection := directionsArray[numSteps%len(directionsArray)]

		log.Debug().
			Interface("CurrentPathNode", currentPathNode).
			Int("NumSteps", numSteps).
			Str("NextDireciton", nextDirection.String()).
			Send()

		if nextDirection == DIRECTION_LEFT {
			currentPathNode = path.pathNodeMap[currentPathNode.LeftNodeLabel]
		} else {
			currentPathNode = path.pathNodeMap[currentPathNode.RightNodeLabel]
		}
		numSteps += 1
	}

	return numSteps, nil
}

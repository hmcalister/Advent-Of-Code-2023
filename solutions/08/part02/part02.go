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

// greatest common divisor (GCD) via Euclidean algorithm
func gcd(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// find Least Common Multiple (LCM) via GCD
func lcm(a, b int, integers ...int) int {
	result := a * b / gcd(a, b)

	for i := 0; i < len(integers); i++ {
		result = lcm(result, integers[i])
	}

	return result
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
	log.Info().Int("DirectionArrayLength", len(directionsArray)).Send()

	path := createPath()
	fileScanner.Scan()
	for fileScanner.Scan() {
		path.parseLineToPathNodeData(fileScanner.Text())
	}

	allStartNodes := make([]*pathNodeData, len(path.allStartPathNodes))
	copy(allStartNodes, path.allStartPathNodes)
	log.Debug().
		Int("NumStartNodes", len(allStartNodes)).
		Interface("StartNodes", allStartNodes).
		Send()

	cumulativeLCM := 1
	var nextDirection directionEnum
	var currentNode *pathNodeData
	var nextNode *pathNodeData
	for startNodeIndex, startNode := range allStartNodes {

		log.Info().
			Int("StartNodeIndex", startNodeIndex).
			Interface("StartNode", startNode).
			Send()

		currentNode = startNode
		step := 0
		for !currentNode.IsTerminal {
			nextDirection = directionsArray[step%len(directionsArray)]
			if nextDirection == DIRECTION_LEFT {
				nextNode = path.pathNodeMap[currentNode.LeftNodeLabel]
			} else {
				nextNode = path.pathNodeMap[currentNode.RightNodeLabel]
			}

			log.Debug().
				Int("StartNodeIndex", startNodeIndex).
				Int("Step", step).
				Interface("NextNode", nextNode).
				Send()

			currentNode = nextNode
			step += 1
		}

		cumulativeLCM = lcm(cumulativeLCM, step)
		log.Info().
			Int("StartNodeIndex", startNodeIndex).
			Interface("StartNode", startNode).
			Interface("TerminalNode", currentNode).
			Int("NumSteps", step).
			Int("CumulativeLCM", cumulativeLCM).
			Send()
	}

	return cumulativeLCM, nil
}

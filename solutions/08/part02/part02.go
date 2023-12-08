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

type pathTraversalData struct {
	// Map from pathNodeLabel to index in our traversal
	pathNodesSeenMap         map[string]int
	pathTraversalNodes       []string
	terminalTraversalIndices []int
}

func createPathTraversalData() *pathTraversalData {
	pathTraversal := &pathTraversalData{}
	pathTraversal.pathNodesSeenMap = make(map[string]int)
	pathTraversal.pathTraversalNodes = make([]string, 0)
	pathTraversal.terminalTraversalIndices = make([]int, 0)
	return pathTraversal
}

// Push a new node onto the path traversal.
//
// If this node has already been seen, do not push and instead
// return the index in which this node occurs
func (pathTraversal *pathTraversalData) pushNextNode(node *pathNodeData) int {
	// Check if we have seen this node before
	if nodeLocation, ok := pathTraversal.pathNodesSeenMap[node.Label]; ok {
		return nodeLocation
	}

	pathTraversal.pathNodesSeenMap[node.Label] = len(pathTraversal.pathTraversalNodes)
	pathTraversal.pathTraversalNodes = append(pathTraversal.pathTraversalNodes, node.Label)
	if node.IsTerminal {
		pathTraversal.terminalTraversalIndices = append(pathTraversal.terminalTraversalIndices, len(pathTraversal.pathTraversalNodes))
	}
	return -1
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

	var nextDirection directionEnum
	var currentNode *pathNodeData
	var nextNode *pathNodeData
	allStartNodeTraversals := make([]*pathTraversalData, len(allStartNodes))
	for startNodeIndex, startNode := range allStartNodes {

		log.Info().
			Int("StartNodeIndex", startNodeIndex).
			Interface("StartNode", startNode).
			Send()

		currentNode = startNode
		allStartNodeTraversals[startNodeIndex] = createPathTraversalData()
		currentNodeTraversal := allStartNodeTraversals[startNodeIndex]
		currentNodeTraversal.pushNextNode(currentNode)
		step := 0
		for {
			nextDirection = directionsArray[step%len(directionsArray)]
			if nextDirection == 'L' {
				nextNode = path.pathNodeMap[currentNode.LeftNodeLabel]
			} else {
				nextNode = path.pathNodeMap[currentNode.RightNodeLabel]
			}

			nextNodeTraversalIndex := currentNodeTraversal.pushNextNode(nextNode)
			log.Debug().
				Int("StartNodeIndex", startNodeIndex).
				Int("Step", step).
				Interface("NextNode", nextNode).
				Int("NextNodeTraversalIndex", nextNodeTraversalIndex).
				Int("TotalPathLength", len(currentNodeTraversal.pathTraversalNodes)).
				Send()
			if nextNodeTraversalIndex != -1 {
				log.Info().
					Int("StartNodeIndex", startNodeIndex).
					Int("ClosedLoopIndex", nextNodeTraversalIndex).
					Int("TotalPathLength", len(currentNodeTraversal.pathTraversalNodes)).
					Int("LoopCycleLength", len(currentNodeTraversal.pathTraversalNodes)-nextNodeTraversalIndex).
					Interface("TerminalLoopIndices", currentNodeTraversal.terminalTraversalIndices).
					// Interface("LoopLabels", currentNodeTraversal.pathTraversalNodes).
					Send()
				break
			}

			step += 1
			currentNode = nextNode
		}
	}

	return 0, nil
}

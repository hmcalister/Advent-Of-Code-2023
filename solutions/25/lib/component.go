package lib

import (
	"bufio"
	"math/rand"
	"strings"

	"github.com/dominikbraun/graph"
	"github.com/rs/zerolog/log"
)

type ComponentGraph struct {
	Graph graph.Graph[string, string]
}

func ParseFileToComponentGraph(fileScanner *bufio.Scanner) *ComponentGraph {
	newComponentGraph := &ComponentGraph{
		Graph: graph.New(graph.StringHash),
	}

	var line string
	for fileScanner.Scan() {
		line = fileScanner.Text()

		colonIndex := strings.IndexRune(line, ':')
		currentComponent := line[:colonIndex]
		neighborComponents := strings.Fields(line[colonIndex+1:])

		log.Debug().
			Str("RawLine", line).
			Str("CurrentComponent", currentComponent).
			Interface("NeighborComponents", neighborComponents).
			Send()

		newComponentGraph.addComponent(currentComponent)
		for _, neighborComponent := range neighborComponents {
			newComponentGraph.addComponent(neighborComponent)
			newComponentGraph.Graph.AddEdge(currentComponent, neighborComponent)
		}
	}

	return newComponentGraph
}

func (compGraph *ComponentGraph) addComponent(componentName string) {
	err := compGraph.Graph.AddVertex(componentName)
	if err != nil {
		log.Trace().Str("Vertex", componentName).Msg("vertex already present in graph")
	}
}


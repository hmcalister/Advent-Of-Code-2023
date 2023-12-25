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

func (compGraph *ComponentGraph) MinimumCut() graph.Graph[string, string] {
	graphCopy, _ := compGraph.Graph.Clone()

	var numVertices int
	numVertices, _ = graphCopy.Order()

	for numVertices > 2 {
		allEdges, _ := graphCopy.Edges()
		contractingEdgeIndex := rand.Int() % len(allEdges)
		contractingEdge := allEdges[contractingEdgeIndex]

		mergeLeftVertex := contractingEdge.Source
		mergeRightVertex := contractingEdge.Target
		graphCopy.RemoveEdge(mergeLeftVertex, mergeRightVertex)
		adjacency, _ := graphCopy.AdjacencyMap()

		newVertexName := mergeLeftVertex + ", " + mergeRightVertex
		graphCopy.AddVertex(newVertexName)

		for neighbor := range adjacency[mergeLeftVertex] {
			graphCopy.AddEdge(newVertexName, neighbor)
			graphCopy.RemoveEdge(mergeLeftVertex, neighbor)
		}

		for neighbor := range adjacency[mergeRightVertex] {
			graphCopy.AddEdge(newVertexName, neighbor)
			graphCopy.RemoveEdge(mergeRightVertex, neighbor)
		}

		graphCopy.RemoveVertex(mergeLeftVertex)
		graphCopy.RemoveVertex(mergeRightVertex)

		numVertices, _ = graphCopy.Order()
	}

	return graphCopy
}

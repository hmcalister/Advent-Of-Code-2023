package lib

import (
	"github.com/dominikbraun/graph"
	"github.com/rs/zerolog/log"
)

type CondensedTrailData struct {
	TrailGraph      graph.Graph[string, Coordinate]
	startCoordinate Coordinate
	endCoordinate   Coordinate
}

type GraphTraversalData struct {
	CurrentVertex   string
	VisitedVertices map[string]interface{}
	TotalDistance   int
}

func (node GraphTraversalData) NextGraphTraversalData(nextVertex string, edgeWeight int) GraphTraversalData {
	nextNode := GraphTraversalData{
		CurrentVertex:   nextVertex,
		VisitedVertices: make(map[string]interface{}),
		TotalDistance:   node.TotalDistance + edgeWeight,
	}
	for k, v := range node.VisitedVertices {
		nextNode.VisitedVertices[k] = v
	}
	nextNode.VisitedVertices[node.CurrentVertex] = visitedCoordinatePresenceIndicator

	return nextNode
}

func ConvertTrailDataToCondensedTrailData(trailData *TrailData) *CondensedTrailData {
	// Define the new graph hash function
	coordinateHash := func(coord Coordinate) string {
		return coord.String()
	}

	condensedTrail := &CondensedTrailData{
		TrailGraph:      graph.New(coordinateHash, graph.Weighted()),
		startCoordinate: trailData.startCoordinate,
		endCoordinate:   trailData.endCoordinate,
	}

	// actually create the new graph from the TrailMap

	getValidMovementDirections := func(coord Coordinate, trailMap map[Coordinate]SurfaceTypeEnum) []DirectionEnum {
		validDirections := make([]DirectionEnum, 0)
		for _, direction := range []DirectionEnum{DIRECTION_UP, DIRECTION_RIGHT, DIRECTION_DOWN, DIRECTION_LEFT} {
			nextCoord := coord.Move(direction)
			surfaceType, isInTrailMap := trailMap[nextCoord]
			if isInTrailMap && surfaceType != SURFACE_FOREST {
				validDirections = append(validDirections, direction)
			}
		}

		return validDirections
	}

	// Define a struct to keep track of what paths we need to explore
	type graphConstructionStruct struct {
		PathNodeData

		previousVertexCoordinate   Coordinate
		distanceFromPreviousVertex int
	}

	graphConstructionQueue := make([]graphConstructionStruct, 0)

	// Add all graph vertexes
	condensedTrail.TrailGraph.AddVertex(trailData.startCoordinate)
	condensedTrail.TrailGraph.AddVertex(trailData.endCoordinate)

	for y := 0; y < trailData.mapHeight; y += 1 {
		for x := 0; x < trailData.mapWidth; x += 1 {
			currentCoord := Coordinate{x, y}
			surface := trailData.trailMap[currentCoord]
			if surface == SURFACE_FOREST {
				continue
			}
			validMoves := getValidMovementDirections(currentCoord, trailData.trailMap)
			if len(validMoves) >= 3 {
				log.Debug().
					Str("VertexCoordinate", currentCoord.String()).
					Str("CurrentSurface", surface.String()).
					Interface("ValidMoves", validMoves).
					Msg("VertexFound")
				condensedTrail.TrailGraph.AddVertex(currentCoord)

				for _, direction := range validMoves {
					graphConstructionQueue = append(graphConstructionQueue, graphConstructionStruct{
						PathNodeData: PathNodeData{
							currentCoordinate: currentCoord.Move(direction),
							visitedCoordinates: map[Coordinate]interface{}{
								currentCoord: visitedCoordinatePresenceIndicator,
							},
						},
						previousVertexCoordinate:   currentCoord,
						distanceFromPreviousVertex: 1,
					}) // Append traversal node
				} // for each direction
			} // if valid moves check
		} //for x
	} // for y

	var currentGraphTraversalData graphConstructionStruct
	for len(graphConstructionQueue) > 0 {
		currentGraphTraversalData, graphConstructionQueue = graphConstructionQueue[len(graphConstructionQueue)-1], graphConstructionQueue[:len(graphConstructionQueue)-1]

		// Check if this traversal has reached a new vertex
		if _, err := condensedTrail.TrailGraph.Vertex(currentGraphTraversalData.currentCoordinate.String()); err == nil {
			log.Debug().
				Str("PreviousVertexCoord", currentGraphTraversalData.previousVertexCoordinate.String()).
				Str("CurrentVertexCoord", currentGraphTraversalData.currentCoordinate.String()).
				Int("EdgeLength", currentGraphTraversalData.distanceFromPreviousVertex).
				Msg("EdgeFound")
			condensedTrail.TrailGraph.AddEdge(currentGraphTraversalData.previousVertexCoordinate.String(), currentGraphTraversalData.currentCoordinate.String(), graph.EdgeWeight(currentGraphTraversalData.distanceFromPreviousVertex))
			continue
		}

		for _, direction := range []DirectionEnum{DIRECTION_UP, DIRECTION_RIGHT, DIRECTION_DOWN, DIRECTION_LEFT} {
			nextCoord := currentGraphTraversalData.currentCoordinate.Move(direction)
			if _, alreadyVisited := currentGraphTraversalData.visitedCoordinates[nextCoord]; alreadyVisited {
				continue
			}

			nextSurfaceType, nextCoordInTrailMap := trailData.trailMap[nextCoord]
			if !nextCoordInTrailMap || nextSurfaceType == SURFACE_FOREST {
				continue
			}

			nextGraphTraversalData := graphConstructionStruct{
				currentGraphTraversalData.NextPathNode(direction),
				currentGraphTraversalData.previousVertexCoordinate,
				currentGraphTraversalData.distanceFromPreviousVertex + 1,
			}
			graphConstructionQueue = append(graphConstructionQueue, nextGraphTraversalData)
		}

	}

	return condensedTrail
}

func (condensedTrail *CondensedTrailData) FindPathNonSlippery() int {

	adjacencyMap, _ := condensedTrail.TrailGraph.AdjacencyMap()
	bestFinishPathLen := -1

	// It appears the start node and end node connect to exactly one other node
	// So use those other nodes as proxy start/end and just add the additional length

	additionalDistance := 0
	var proxyStartVertex string
	var proxyEndVertex string

	if len(adjacencyMap[condensedTrail.startCoordinate.String()]) != 1 {
		log.Fatal().Msg("start vertex has more than one neighbor")
	}
	for k, v := range adjacencyMap[condensedTrail.startCoordinate.String()] {
		proxyStartVertex = k
		additionalDistance += v.Properties.Weight
		log.Debug().Str("ProxyStartVertex", proxyStartVertex).Send()
	}
	if len(adjacencyMap[condensedTrail.endCoordinate.String()]) != 1 {
		log.Fatal().Msg("end vertex has more than one neighbor")
	}
	for k, v := range adjacencyMap[condensedTrail.endCoordinate.String()] {
		proxyEndVertex = k
		additionalDistance += v.Properties.Weight
		log.Debug().Str("ProxyEndVertex", proxyEndVertex).Send()
	}

	graphTraversalList := make([]GraphTraversalData, 0)
	graphTraversalList = append(graphTraversalList, GraphTraversalData{
		CurrentVertex:   proxyStartVertex,
		VisitedVertices: make(map[string]interface{}),
		TotalDistance:   additionalDistance,
	})

	var currentTraversalData GraphTraversalData
	for len(graphTraversalList) > 0 {
		currentTraversalData, graphTraversalList = graphTraversalList[len(graphTraversalList)-1], graphTraversalList[:len(graphTraversalList)-1]

		log.Debug().
			Str("CurrentCoord", currentTraversalData.CurrentVertex).
			Int("CurrentPathLen", currentTraversalData.TotalDistance).
			Int("QueueLength", len(graphTraversalList)).
			Int("BestFinishPathLength", bestFinishPathLen).
			Msg("PathFindCurrentNode")

		if currentTraversalData.CurrentVertex == proxyEndVertex {
			log.Debug().
				Int("FinishedPathNode", currentTraversalData.TotalDistance).
				Msg("CompletePathFound")

			if bestFinishPathLen < currentTraversalData.TotalDistance {
				bestFinishPathLen = currentTraversalData.TotalDistance
				log.Info().
					Int("BestFinishPathLength", bestFinishPathLen).
					Msg("NewBestPathFound")
			}
			continue
		}

		neighborsMap := adjacencyMap[currentTraversalData.CurrentVertex]
		for neighborVertex, edgeWeight := range neighborsMap {
			if _, ok := currentTraversalData.VisitedVertices[neighborVertex]; ok {
				continue
			}

			nextTraversal := currentTraversalData.NextGraphTraversalData(neighborVertex, edgeWeight.Properties.Weight)
			graphTraversalList = append(graphTraversalList, nextTraversal)
		}
	}

	return bestFinishPathLen
}

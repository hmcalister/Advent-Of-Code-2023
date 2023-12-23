package lib

import (
	"bufio"
	"container/heap"
	"errors"
	"fmt"
	"sort"

	"github.com/rs/zerolog/log"
)

type TrailData struct {
	trailMap        map[Coordinate]SurfaceTypeEnum
	startCoordinate Coordinate
	endCoordinate   Coordinate
	mapWidth        int
	mapHeight       int
}

func ParseFileToTrail(fileScanner *bufio.Scanner) *TrailData {
	trail := &TrailData{
		trailMap: make(map[Coordinate]SurfaceTypeEnum),
	}

	var line string
	for fileScanner.Scan() {
		line = fileScanner.Text()

		log.Debug().
			Int("LineNumber", trail.mapHeight).
			Str("RawLine", line).
			Int("NumRunes", len(line)).
			Send()

		for runeIndex, r := range line {
			currentCoordinate := Coordinate{
				X: runeIndex,
				Y: trail.mapHeight,
			}

			currentSurface, ok := runeToSurfaceTypeMap[r]
			if !ok {
				log.Fatal().Msgf("found unexpected surface rune(coordinate %v)  %v in line %v ", currentCoordinate.String(), r, line)
			}

			log.Trace().
				Str("CurrentCoordinate", currentCoordinate.String()).
				Str("CurrentSurface", currentSurface.String()).
				Str("RawRune", string(r)).
				Send()

			trail.trailMap[currentCoordinate] = currentSurface
		}

		trail.mapWidth = max(trail.mapWidth, len(line))
		trail.mapHeight += 1
	}

	for x := 0; x < trail.mapWidth; x += 1 {
		firstLineCoord := Coordinate{
			X: x,
			Y: 0,
		}
		if trail.trailMap[firstLineCoord] == SURFACE_PATH {
			trail.startCoordinate = firstLineCoord
		}

		lastLineCoordinate := Coordinate{
			X: x,
			Y: trail.mapHeight - 1,
		}
		if trail.trailMap[lastLineCoordinate] == SURFACE_PATH {
			trail.endCoordinate = lastLineCoordinate
		}
	}

	log.Debug().
		Str("StartCoordinate", trail.startCoordinate.String()).
		Str("EndCoordinate", trail.endCoordinate.String()).
		Send()

	return trail
}

func (trail *TrailData) VisualizePath(path PathNodeData) {
	line := make([]rune, trail.mapWidth)
	trailLine := make([]rune, trail.mapWidth)
	for y := 0; y < trail.mapHeight; y += 1 {
		for x := 0; x < trail.mapWidth; x += 1 {
			currentCoord := Coordinate{
				X: x,
				Y: y,
			}

			surfaceType, isInTrailMap := trail.trailMap[currentCoord]
			trailLine[x] = surfaceTypeToRuneMap[surfaceType]

			if _, isInPath := path.visitedCoordinates[currentCoord]; isInPath {
				line[x] = 'O'
			} else if isInTrailMap {
				line[x] = surfaceTypeToRuneMap[surfaceType]
			} else {
				log.Fatal().Msgf("failed to identify coordinate %v in either path or trailMap", currentCoord)
			}

		}
		log.Info().
			Str(fmt.Sprintf("Line %03v", y), string(line)).
			Str(fmt.Sprintf("TrailLine %03v", y), string(trailLine)).
			Msg("VisualizePath")
	}
}

// func (trail *TrailData) FindShortestPath() PathNodeData {

// }

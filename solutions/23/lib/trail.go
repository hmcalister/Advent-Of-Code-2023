package lib

import (
	"bufio"

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
	for y := 0; y < trail.mapHeight; y += 1 {
		line := make([]rune, trail.mapWidth)
		for x := 0; x < trail.mapWidth; x += 1 {
			currentCoord := Coordinate{
				X: x,
				Y: y,
			}

			if _, isInPath := path.VisitedCoordinates[currentCoord]; isInPath {
				line[x] = 'O'
			} else if surfaceType, isInTrailMap := trail.trailMap[currentCoord]; isInTrailMap {
				line[x] = surfaceTypeToRuneMap[surfaceType]
			} else {
				log.Fatal().Msgf("failed to identify coordinate %v in either path or trailMap", currentCoord)
			}
		}
		log.Info().
			Str("Line %05v: ", string(line)).
			Msg("VisualizePath")
	}
}

// func (trail *TrailData) FindShortestPath() PathNodeData {

// }

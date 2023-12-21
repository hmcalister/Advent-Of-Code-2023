package lib

import (
	"bufio"

	"github.com/rs/zerolog/log"
)

//go:generate stringer -type=SurfaceTypeEnum
type SurfaceTypeEnum int

const (
	SURFACE_PLOT SurfaceTypeEnum = iota
	SURFACE_ROCK SurfaceTypeEnum = iota
)

type GardenData struct {
	SurfaceData     map[coordinate]SurfaceTypeEnum
	StartCoordinate coordinate
	MapWidth        int
	MapHeight       int
}

func ParseFileToGardenData(fileScanner bufio.Scanner) GardenData {
	newGarden := GardenData{
		SurfaceData: make(map[coordinate]SurfaceTypeEnum),
	}

	yPosition := 0
	for fileScanner.Scan() {
		line := fileScanner.Text()
		for runeIndex, r := range line {
			currentCoordinate := coordinate{X: runeIndex, Y: yPosition}
			newGarden.MapWidth = max(newGarden.MapWidth, runeIndex)
			newGarden.MapHeight = max(newGarden.MapHeight, yPosition)
			log.Trace().
				Str("RawLine", line).
				Send()
			switch r {
			case '.':
				newGarden.SurfaceData[currentCoordinate] = SURFACE_PLOT
			case '#':
				newGarden.SurfaceData[currentCoordinate] = SURFACE_ROCK
			case 'S':
				newGarden.StartCoordinate = currentCoordinate
				newGarden.SurfaceData[currentCoordinate] = SURFACE_PLOT
			default:
				log.Fatal().Msgf("unexpected rune %v in line %v while parsing file to gardenData", r, line)
			}
		}
		yPosition += 1
	}
	return newGarden
}

func (garden GardenData) DebugLog() {
	log.Debug().
		Interface("StartCoordinate", garden.StartCoordinate).
		Int("NumCoordinatesMapped", len(garden.SurfaceData)).
		Int("MapWidth", garden.MapWidth).
		Int("MapHeight", garden.MapHeight).
		Msg("GardenDebug")
}

func (garden GardenData) NumReachableGardensInExactlyNumSteps(maxSteps int) int {
	PRESENCE_INDICATOR := struct{}{}
	DIRECTIONS := []DirectionEnum{DIRECTION_UP, DIRECTION_RIGHT, DIRECTION_DOWN, DIRECTION_LEFT}

	// The set of all coordinates that *can* be reached in the next step
	var nextPlots map[coordinate]interface{}
	var currentPlots map[coordinate]interface{}

	// The set of all gardens that were reached by the current step
	nextPlots = map[coordinate]interface{}{
		garden.StartCoordinate: PRESENCE_INDICATOR,
	}

	for stepNumber := 0; stepNumber < maxSteps; stepNumber += 1 {
		currentPlots = nextPlots
		nextPlots = make(map[coordinate]interface{})

		log.Debug().
			Int("CurrentStepNumber", stepNumber).
			Int("NumPlotsToConsider", len(currentPlots)).
			Send()

		for coord := range currentPlots {
			for _, direction := range DIRECTIONS {
				nextCoord := coord.Move(direction)

				nextStep, ok := garden.SurfaceData[nextCoord]
				if nextStep != SURFACE_PLOT || !ok {
					continue
				}
				log.Trace().
					Int("StepNumber", stepNumber).
					Interface("CurrentCoord", coord).
					Str("MovementDirection", direction.String()).
					Msg("FoundPlot")

				nextPlots[nextCoord] = PRESENCE_INDICATOR
			}
		}
	}

	return len(nextPlots)
}

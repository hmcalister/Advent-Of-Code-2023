package lib

import (
	"bufio"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

type DigLayoutData struct {
	// The coordinate of each trench turn
	trenchCoordinates []coordinate

	// The digLayout dimensions
	XLim int
	YLim int
}

// Parse an individual line from the input to the corresponding information
//
// # The direction the trench is to be dug
//
// # The number of spaces to be dug
func parseLineData(line string) (DirectionEnum, int) {
	fields := strings.Fields(line)

	// Take off parentheses and leading #
	colorString := fields[2][2 : len(fields[2])-1]

	distanceStr := colorString[:5]
	directionStr := colorString[5:]

	trenchDirection := directionDecoderMap[directionStr]

	numSpaces64, err := strconv.ParseInt(distanceStr, 16, 0)
	if err != nil {
		log.Fatal().Msgf("failed to parse hex encoded distance string %v in line %v", distanceStr, line)
	}
	numSpaces := int(numSpaces64)

	log.Trace().
		Str("RawLine", line).
		Str("ParsedDirection", trenchDirection.String()).
		Int("ParsedSpaces", numSpaces).
		Send()

	return trenchDirection, numSpaces
}

func NewDigLayoutFromFileScanner(fileScanner *bufio.Scanner) *DigLayoutData {
	var line string
	var currentCoordinate coordinate

	currentCoordinate = coordinate{0, 0}

	digLayout := &DigLayoutData{
		trenchCoordinates: make([]coordinate, 0),
		XLim:              0,
		YLim:              0,
	}

	// digLayout.trenchCoordinates = append(digLayout.trenchCoordinates, currentCoordinate)

	// Parse each line in the file, creating new trenches as we go
	for fileScanner.Scan() {
		line = fileScanner.Text()
		trenchDirection, trenchLength := parseLineData(line)
		log.Debug().
			Str("TrenchDirection", trenchDirection.String()).
			Int("TrenchLength", trenchLength).
			Interface("CoordinateInitial", currentCoordinate).
			Interface("CoordinateFinal", currentCoordinate.Move(trenchDirection, trenchLength)).
			Send()

		currentCoordinate = currentCoordinate.Move(trenchDirection, trenchLength)
		digLayout.updateLimitsAndRowStretches(currentCoordinate)
		digLayout.trenchCoordinates = append(digLayout.trenchCoordinates, currentCoordinate)
	}

	return digLayout
}

func (digLayout *DigLayoutData) updateLimitsAndRowStretches(newCoordinate coordinate) {
	digLayout.XLim = max(digLayout.XLim, newCoordinate.X+1)
	digLayout.YLim = max(digLayout.YLim, newCoordinate.Y+1)
}

func (digLayout *DigLayoutData) VisualizeDigLayout() {
	log.Info().
		Int("XLim", digLayout.XLim).
		Int("YLim", digLayout.YLim).
		Send()

	for trenchIndex, trenchCoordinate := range digLayout.trenchCoordinates {
		log.Debug().
			Int("TrenchIndex", trenchIndex).
			Interface("TrenchCoordinate", trenchCoordinate).
			Send()
	}
}

func (digLayout *DigLayoutData) CalculateTotalVolume() int {
	totalArea := 0

	// Use shoelace formula to determine total internal area
	for i := 0; i < len(digLayout.trenchCoordinates); i += 1 {
		coordOne := digLayout.trenchCoordinates[i]
		coordTwo := digLayout.trenchCoordinates[(i+1)%len(digLayout.trenchCoordinates)]

		// Add the determinate of the matrix:
		// | coordOne.X  coordTwo.X |
		// | coordOne.Y  coordTwo.Y |
		totalArea += coordOne.X*coordTwo.Y - coordOne.Y*coordTwo.X

		// Then also add the "area" due to the perimeter trenches
		//
		// It is okay to half this due to Picks theorem
		totalArea += abs(coordOne.X-coordTwo.X) + abs(coordOne.Y-coordTwo.Y)
	}

	// See Picks theorem
	return totalArea/2 + 1
}

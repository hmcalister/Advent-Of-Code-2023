package lib

import (
	"bufio"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

type DigLayoutData struct {
	// The stretch data for each row
	rowStretchCovers []rowStretchCoverData

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
	var previousCoordinate coordinate

	currentCoordinate = coordinate{0, 0}
	previousCoordinate = currentCoordinate

	digLayout := &DigLayoutData{
		rowStretchCovers: make([]rowStretchCoverData, 0),
		XLim:             0,
		YLim:             0,
	}

	// Parse each line in the file, creating new trenches as we go
	//
	// We also update the rowStretches as we go, so we can figure out how many interior points there are
	for fileScanner.Scan() {
		line = fileScanner.Text()
		trenchDirection, trenchLength := parseLineData(line)
		currentCoordinate = currentCoordinate.Move(trenchDirection, trenchLength)

		if trenchDirection == DIRECTION_RIGHT {
			newStrechData := stretchData{
				stretchStart: previousCoordinate.X,
				stretchLen:   currentCoordinate.X - previousCoordinate.X,
				interior:     true,
			}
			rowIndex := currentCoordinate.Y
			digLayout.rowStretchCovers[rowIndex].
		}
	}

	return digLayout
}

func (digLayout *DigLayoutData) updateLimitsAndRowStretches(newCoordinate coordinate) {
	digLayout.XLim = max(digLayout.XLim, newCoordinate.X+1)
	digLayout.YLim = max(digLayout.YLim, newCoordinate.Y+1)

	for y := len(digLayout.rowStretchCovers); y < digLayout.YLim; y += 1 {
		digLayout.rowStretchCovers = append(digLayout.rowStretchCovers, newRowStretchCover())
	}
}

func (digLayout *DigLayoutData) VisualizeDigLayout() {
	log.Info().
		Int("XLim", digLayout.XLim).
		Int("YLim", digLayout.YLim).
		Send()
}

func (digLayout *DigLayoutData) ExcavateInterior() {
	var isInterior bool

	interiorTrench := TrenchData{
		DigDirection: DIRECTION_NONE,
	}

	// Skip the first and last row and column, as these can never have interior points
	for y := digLayout.YMin; y < digLayout.YMax-1; y += 1 {
		pbar.Add(1)
		isInterior = false

		for x := digLayout.XMin; x < digLayout.XMax-1; x += 1 {
			currentCoordinate := coordinate{x, y}
			// Check if we have encountered a new trench edge

			if trench, ok := digLayout.DigMap[currentCoordinate]; ok {
				var ok bool
				var capTrench TrenchData
				if trench.DigDirection == DIRECTION_UP {
					// Skip past any remaining trenches
					for ; x < digLayout.XMax-1; x += 1 {
						currentCoordinate = coordinate{x, y}
						if _, ok := digLayout.DigMap[currentCoordinate]; !ok {
							x -= 1
							currentCoordinate = coordinate{x, y}
							break
						}
					}
					if capTrench, ok = digLayout.DigMap[currentCoordinate.Move(DIRECTION_UP)]; ok && capTrench.DigDirection == DIRECTION_UP {
						// ..................^....
						// ....^>>>>>>>>>>>>>>....
						// ....^..................
						isInterior = true
					}

					if capTrench, ok = digLayout.DigMap[currentCoordinate.Move(DIRECTION_DOWN)]; ok && capTrench.DigDirection == DIRECTION_UP {
						// ....^...................
						// ....^<<<<<<<<<<<<<<<....
						// ...................^...
						isInterior = true
					}
				} // trench.DigDirection == DIRECTION_UP

				if trench.DigDirection == DIRECTION_DOWN {
					// Skip past any remaining trenches
					for ; x < digLayout.XMax-1; x += 1 {
						currentCoordinate = coordinate{x, y}
						if _, ok := digLayout.DigMap[currentCoordinate]; !ok {
							x -= 1
							currentCoordinate = coordinate{x, y}
							break
						}
					}

					if capTrench, ok = digLayout.DigMap[currentCoordinate.Move(DIRECTION_UP)]; ok && capTrench.DigDirection == DIRECTION_DOWN {
						// ..................v....
						// ....v<<<<<<<<<<<<<v....
						// ....v..................
						isInterior = false
					}

					if capTrench, ok = digLayout.DigMap[currentCoordinate.Move(DIRECTION_DOWN)]; ok && capTrench.DigDirection == DIRECTION_DOWN {
						// ....v...................
						// ....v>>>>>>>>>>>>>>>....
						// ...................v...
						isInterior = false
					}
				} //trench.DigDirection == DIRECTION_DOWN

				continue
			}

			// If we are not on the interior, skip this coord
			if !isInterior {
				continue
			}

			digLayout.DigMap[currentCoordinate] = interiorTrench
		} // end xLoop
	} // end yLoop
}

func (digLayout *DigLayoutData) CalculateTotalVolume() int {
	return len(digLayout.DigMap)
}

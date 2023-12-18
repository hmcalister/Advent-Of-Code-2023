package lib

import (
	"bufio"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

type DigLayoutData struct {
	// A map from coordinate to the trench at that coordinate
	DigMap map[coordinate]TrenchData

	// The current depth of the digLayout
	CurrentDepth int

	// The digLayout dimensions
	XMin int
	XMax int
	YMin int
	YMax int
}

// Parse an individual line from the input to the corresponding information
//
// # The direction the trench is to be dug
//
// # The number of spaces to be dug
//
// And the color to paint the trench edges
func parseLineData(line string) (DirectionEnum, int, ColorData) {
	fields := strings.Fields(line)

	trenchDirection := directionDecoderMap[fields[0]]
	numSpaces, err := strconv.Atoi(fields[1])
	if err != nil {
		log.Fatal().Msgf("failed to parse number of spaces %v in line %v", fields[1], line)
	}
	// Take off parentheses
	colorString := fields[2][1 : len(fields[2])-1]
	color := ColorData{ColorString: colorString}

	log.Trace().
		Str("RawLine", line).
		Str("ParsedDirection", trenchDirection.String()).
		Int("ParsedSpaces", numSpaces).
		Interface("ParsedColor", color)

	return trenchDirection, numSpaces, color
}

func NewDigLayoutFromFileScanner(fileScanner *bufio.Scanner) *DigLayoutData {
	var line string
	var currentCoordinate coordinate
	var previousTrenchStretchEndCoordinate coordinate

	currentCoordinate = coordinate{0, 0}
	previousTrenchStretchEndCoordinate = currentCoordinate
	digLayout := &DigLayoutData{
		DigMap:       make(map[coordinate]TrenchData),
		CurrentDepth: 1,
		XMin:         0,
		XMax:         0,
		YMin:         0,
		YMax:         0,
	}

	digLayout.DigMap[currentCoordinate] = TrenchData{
		DigDirection: DIRECTION_RIGHT,
		EdgeColors:   make(map[DirectionEnum]ColorData),
	}

	// Each line in the file corresponds to a straight trench in the dig
	for fileScanner.Scan() {
		line = fileScanner.Text()
		trenchDirection, numSpaces, color := parseLineData(line)
		trench := newTrench(trenchDirection, digLayout.CurrentDepth, color)

		// Move along the trench and update the digmap as we go
		for i := 0; i < numSpaces; i += 1 {
			currentCoordinate = currentCoordinate.Move(trenchDirection)
			digLayout.DigMap[currentCoordinate] = trench
			log.Trace().Interface("NewTrench", currentCoordinate).Str("TrenchDirection", trenchDirection.String()).Send()

			digLayout.XMin = min(digLayout.XMin, currentCoordinate.X)
			digLayout.XMax = max(digLayout.XMax, currentCoordinate.X+1)
			digLayout.YMin = min(digLayout.YMin, currentCoordinate.Y)
			digLayout.YMax = max(digLayout.YMax, currentCoordinate.Y+1)
		}

		// Update the previous trench stretch with the new edge color
		previousStretchEnd := digLayout.DigMap[previousTrenchStretchEndCoordinate]
		digLayout.DigMap[previousTrenchStretchEndCoordinate] = previousStretchEnd.updateEdgeColors(trenchDirection, color)
	}

	return digLayout
}

func (digLayout *DigLayoutData) VisualizeDigLayout() {
	var trenchRune rune
	layoutVisualization := make([][]rune, digLayout.YMax-digLayout.YMin)
	for y := digLayout.YMin; y < digLayout.YMax; y += 1 {
		layoutVisualization[y-digLayout.YMin] = make([]rune, digLayout.XMax-digLayout.XMin)
		for x := digLayout.XMin; x < digLayout.XMax; x += 1 {
			layoutVisualization[y-digLayout.YMin][x-digLayout.XMin] = '.'
		}
	}
	for coordinate, trench := range digLayout.DigMap {
		x, y := coordinate.X-digLayout.XMin, coordinate.Y-digLayout.YMin
		switch trench.DigDirection {
		case DIRECTION_UP:
			trenchRune = '^'
		case DIRECTION_RIGHT:
			trenchRune = '>'
		case DIRECTION_DOWN:
			trenchRune = 'v'
		case DIRECTION_LEFT:
			trenchRune = '<'
		case DIRECTION_NONE:
			trenchRune = '#'
		}
		layoutVisualization[y][x] = trenchRune
	}
	for y := digLayout.YMin; y < digLayout.YMax; y += 1 {
		log.Info().Str("DigLayout", string(layoutVisualization[y-digLayout.YMin])).Send()
	}
}

func (digLayout *DigLayoutData) ExcavateInterior() {
	var isInterior bool

	interiorTrench := TrenchData{
		DigDirection: DIRECTION_NONE,
		Depth:        digLayout.CurrentDepth,
		EdgeColors:   make(map[DirectionEnum]ColorData),
	}

	// Skip the first and last row and column, as these can never have interior points
	for y := digLayout.YMin; y < digLayout.YMax-1; y += 1 {
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
		}
	}
}

func (digLayout *DigLayoutData) CalculateTotalVolume() int {
	totalVolume := 0
	for _, trench := range digLayout.DigMap {
		totalVolume += trench.Depth
	}

	return totalVolume
}

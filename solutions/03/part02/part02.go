package part02

import (
	"bufio"
	"slices"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	LINE_LENGTH = 141
	DIGITS      = "0123456789"
	GEAR_SYMBOL = '*'
)

var (
	COORD_OFFSETS = [3]int{-1, 0, 1}
)

type gearData struct {
	// The position of this gear, linearized
	Location int

	// The gear ratio (product of all nearby numbers)
	Ratio int

	// An array of all the unique neighbors of this gear
	UniqueNeighbors []int
}

type partNumberData struct {
	// A unique ID for this part number
	PartID int

	// The actual part number
	Number int
}

// An array of partNumbers
var partNumberArray []*partNumberData

// An array of gearData
var gearDataArray []*gearData

// A map from a coordinate to the corresponding part number (if one exists at that coordinate)
//
// Note that indices into this map should be of the form (LINE_LENGTH*y + x) to linearize the coordinates.
var partNumberMap map[int]*partNumberData

// Given a scanner over the puzzle input, calculate the sum of the part numbers.
//
// This is done by finding all numbers adjacent (incl. diagonally) with a symbol (non-period characters).
func ProcessInput(fileScanner *bufio.Scanner) (int, error) {
	result := 0
	var currentRune rune
	partNumberMap = make(map[int]*partNumberData)
	partNumberArray = make([]*partNumberData, 0)
	gearDataArray = make([]*gearData, 0)

	lineNumber := 0
	currentPartID := 0
	for fileScanner.Scan() {
		log.Debug().Int("LineNumber", lineNumber).Send()
		line := fileScanner.Text()

		for colIndex := 0; colIndex < len(line); colIndex += 1 {
			currentRune = rune(line[colIndex])

			// We are inspecting a new character. The options are:
			// - Gear Symbol: Add this coordinate to the gear list and carry on
			// - Digit: Parse this digit, add it to a temporary int, and continue
			// 		along the line until the entire digit is parsed.
			// 		Beware of end of line!

			if currentRune == GEAR_SYMBOL {
				log.Trace().
					Int("ColIndex", colIndex).
					Str("Symbol", string(GEAR_SYMBOL)).
					Send()
				gearDataArray = append(gearDataArray, &gearData{
					Location:        LINE_LENGTH*lineNumber + colIndex,
					Ratio:           1,
					UniqueNeighbors: make([]int, 0),
				})
			} else if strings.ContainsRune(DIGITS, currentRune) {
				currentData := &partNumberData{
					PartID: currentPartID,
				}
				currentPartID += 1

				for colIndex < len(line) {
					currentRune = rune(line[colIndex])
					if !strings.ContainsRune(DIGITS, currentRune) {
						break
					}

					currentDigit := strings.IndexRune(DIGITS, currentRune)
					currentData.Number = 10*currentData.Number + currentDigit
					log.Trace().
						Int("ColIndex", colIndex).
						Str("Symbol", string(currentRune)).
						Bool("RuneIsDigit", strings.ContainsRune(DIGITS, currentRune)).
						Int("ParsedInt", currentDigit).
						Int("CumulativeInt", currentData.Number).
						Send()
					partNumberMap[LINE_LENGTH*lineNumber+colIndex] = currentData

					colIndex += 1
				}
				colIndex -= 1
				log.Debug().
					Int("ColIndex", colIndex).
					Int("LinearCoordinate", LINE_LENGTH*lineNumber+colIndex).
					Int("FoundNumber", currentData.Number).
					Send()
			} else {
				log.Trace().
					Int("ColIndex", colIndex).
					Str("Symbol", string(currentRune)).
					Send()
			}

		}
		lineNumber += 1
	}

	// We now have the locations of all gears, as well as a hashmap from
	// location to partNumbers. We can now go to each gear and calculate
	// the gear ratio as well as the unique partNumbers around it.
	for _, gear := range gearDataArray {
		log.Debug().
			Int("CurrentGearLocation", gear.Location).
			Array("EffectiveCartesianCoordinates", zerolog.Arr().
				Int(gear.Location%LINE_LENGTH).
				Int(gear.Location/LINE_LENGTH)).
			Send()
		for _, delX := range COORD_OFFSETS {
			for _, delY := range COORD_OFFSETS {
				currentCoord := gear.Location + (LINE_LENGTH*delY + delX)
				partNumber, ok := partNumberMap[currentCoord]
				if ok && !slices.Contains(gear.UniqueNeighbors, partNumber.PartID) {
					// We have found a part that is not yet counted as a neighbor!
					gear.UniqueNeighbors = append(gear.UniqueNeighbors, partNumber.PartID)
					gear.Ratio *= partNumber.Number

					log.Debug().
						Int("FoundPartID", partNumber.PartID).
						Int("NewRatio", gear.Ratio).
						Send()
				}
			}
		}

		if len(gear.UniqueNeighbors) == 2 {
			result += gear.Ratio
		}
	}

	return result, nil
}

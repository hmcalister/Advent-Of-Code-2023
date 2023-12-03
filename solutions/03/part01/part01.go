package part01

import (
	"bufio"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	LINE_LENGTH = 141
)

var (
	COORD_OFFSETS = [3]int{-1, 0, 1}
	SYMBOLS       = "$@/*#=+-&%"
	DIGITS        = "0123456789"
)

type partNumberData struct {
	// The actual part number
	Number int

	// Boolean to see if this part number has already been counted
	Counted bool
}

// A map from a coordinate to the corresponding part number (if one exists at that coordinate)
//
// Note that indices into this map should be of the form (LINE_LENGTH*y + x) to linearize the coordinates.
var partNumberMap map[int]*partNumberData

// A list of positions in which symbols are found.
var symbolCoordinateList []int

// Given a scanner over the puzzle input, calculate the sum of the part numbers.
//
// This is done by finding all numbers adjacent (incl. diagonally) with a symbol (non-period characters).
func ProcessInput(fileScanner *bufio.Scanner) (int, error) {
	result := 0
	partNumberMap = make(map[int]*partNumberData)
	symbolCoordinateList = make([]int, 0)
	var currentRune rune

	lineNumber := 0
	// Walk over each line and find the symbols and part numbers,
	// storing each in the respective map and list.
	for fileScanner.Scan() {
		log.Debug().Int("LineNumber", lineNumber).Send()
		line := fileScanner.Text()

		for colIndex := 0; colIndex < len(line); colIndex += 1 {
			currentRune = rune(line[colIndex])

			// We are inspecting a new character. The options are:
			// - period: Do nothing, carry on
			// - Symbol: Add this coordinate to the symbol list and carry on
			// - Digit: Parse this digit, add it to a temporary int, and continue
			// 		along the line until the entire digit is parsed.
			// 		Beware of end of line!

			if currentRune == '.' {
				log.Trace().
					Int("ColIndex", colIndex).
					Str("Symbol", ".").
					Send()
				continue
			} else if strings.ContainsRune(SYMBOLS, currentRune) {
				log.Trace().
					Int("ColIndex", colIndex).
					Str("Symbol", string(currentRune)).
					Send()
				symbolCoordinateList = append(symbolCoordinateList, LINE_LENGTH*lineNumber+colIndex)
			} else if strings.ContainsRune(DIGITS, currentRune) {
				currentData := &partNumberData{}
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
				log.Panic().Msgf("found unexpected symbol: %v", string(currentRune))
			}
		}

		lineNumber += 1
	}
	// We now have the locations of all symbols, as well as a hashmap from
	// location to partNumbers. The rest is easy! Look at each symbol location,
	// walk over the neighbors (in all directions) and add the part number to the sum!
	//
	// Note the hashmap should allow us to check for coordinates that are
	// "outside" the schematic, i.e. we don't have to do any boundary checking!

	for _, symbolLocation := range symbolCoordinateList {
		for _, delX := range COORD_OFFSETS {
			for _, delY := range COORD_OFFSETS {
				currentCoord := symbolLocation + (LINE_LENGTH*delY + delX)
				partNumber, ok := partNumberMap[currentCoord]
				if ok && !partNumber.Counted {
					partNumber.Counted = true
					result += partNumber.Number
					log.Debug().
						Int("CurrentLinearCoordinate", currentCoord).
						Array("EffectiveCartesianCoordinates", zerolog.Arr().Int(currentCoord%LINE_LENGTH).Int(currentCoord/LINE_LENGTH)).
						Int("FoundPartNumber", partNumber.Number).
						Int("NewCount", result).
						Send()
				}
			}
		}
	}

	return result, nil
}

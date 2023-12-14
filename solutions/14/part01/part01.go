package part01

import (
	"bufio"

	"github.com/rs/zerolog/log"
)

const (
	EMPTY_SPACE  byte = '.'
	ROUNDED_ROCK byte = 'O'
	CUBE_ROCK    byte = '#'
)

type PlatformData struct {
	initialRows []string
	currentRows [][]byte
}

func NewPlatformData(rowStrings []string) PlatformData {
	rowData := make([][]byte, len(rowStrings))
	for rowIndex, row := range rowStrings {
		currRow := make([]byte, len(row))
		for colIndex, col := range row {
			currRow[colIndex] = byte(col)
		}
		rowData[rowIndex] = currRow
	}

	return PlatformData{
		initialRows: rowStrings,
		currentRows: rowData,
	}
}

func (platform PlatformData) LogCurrentRows() {
	for rowIndex, row := range platform.currentRows {
		log.Debug().
			Str("Row", string(row)).
			Int("RowIndex", rowIndex).
			Send()
	}
}

func (platform PlatformData) RollNorth() {
	for colIndex := 0; colIndex < len(platform.currentRows[0]); colIndex += 1 {
		stoppingIndex := -1
		for rowIndex := 0; rowIndex < len(platform.currentRows); rowIndex += 1 {
			currentState := platform.currentRows[rowIndex][colIndex]
			log.Trace().
				Interface("Coordinates", []int{rowIndex, colIndex}).
				Str("CurrentState", string(currentState)).
				Send()
			switch currentState {
			case EMPTY_SPACE:
				// If the current space is empty, we need not update anything and can simply continue
				continue
			case CUBE_ROCK:
				// If we encounter a new cube rock, we must update the stoppingIndex
				// Such that the next rounded rock will only roll to this position
				stoppingIndex = rowIndex
			case ROUNDED_ROCK:
				// If we encounter a rounded rock that can roll, it's time to update things!
				// The rock will roll to the stoppingIndex, as tracked above.
				// Importantly, the stoppingIndex is the index of the *block*,
				// so we only roll to one before it

				// This position is now empty, as the rock will roll
				platform.currentRows[rowIndex][colIndex] = EMPTY_SPACE

				// The rounded rock in this position will roll away to the stopping index
				platform.currentRows[stoppingIndex+1][colIndex] = ROUNDED_ROCK

				// The stoppingIndex of this column is now the roundedRock we just updated
				stoppingIndex = stoppingIndex + 1
			default:
				log.Fatal().Msgf("Encountered unknown character at (%v, %v): %v", colIndex, rowIndex, currentState)
			}
		}
	}
}

func (platform PlatformData) CalculateLoad() int {
	totalLoad := 0
	for rowIndex := 0; rowIndex < len(platform.currentRows); rowIndex += 1 {
		for colIndex := 0; colIndex < len(platform.currentRows[0]); colIndex += 1 {
			currentState := platform.currentRows[rowIndex][colIndex]
			if currentState == ROUNDED_ROCK {
				totalLoad += len(platform.currentRows) - rowIndex
			}
		}
	}

	return totalLoad
}

func ProcessInput(fileScanner *bufio.Scanner) (int, error) {
	rowStrings := make([]string, 0)
	for fileScanner.Scan() {
		line := fileScanner.Text()
		rowStrings = append(rowStrings, line)
		log.Debug().
			Str("ParsedLine", line).
			Send()
	}

	platform := NewPlatformData(rowStrings)
	platform.LogCurrentRows()
	log.Debug().Msg("Roll North")
	platform.RollNorth()
	platform.LogCurrentRows()
	totalLoad := platform.CalculateLoad()

	return totalLoad, nil
}

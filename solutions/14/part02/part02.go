package part02

import (
	"bufio"
	"hash"
	"hash/fnv"
	"slices"
	"strings"

	"github.com/rs/zerolog/log"
)

const (
	EMPTY_SPACE  byte = '.'
	ROUNDED_ROCK byte = 'O'
	CUBE_ROCK    byte = '#'
)

type PlatformData struct {
	rows         [][]byte
	cache        map[uint64]string
	hashIndices  []uint64
	hashFunction hash.Hash64
}

func convertStringsToByteArrays(strings []string) [][]byte {
	byteData := make([][]byte, len(strings))
	for rowIndex, row := range strings {
		currRow := make([]byte, len(row))
		for colIndex, col := range row {
			currRow[colIndex] = byte(col)
		}
		byteData[rowIndex] = currRow
	}

	return byteData
}

func convertByteArraysToStrings(byteArrays [][]byte) []string {
	strings := make([]string, len(byteArrays))
	for byteArrayIndex, byteArray := range byteArrays {
		strings[byteArrayIndex] = string(byteArray)
	}
	return strings
}

func NewPlatformData(rowStrings []string) *PlatformData {
	rowData := convertStringsToByteArrays(rowStrings)

	return &PlatformData{
		rows:         rowData,
		cache:        make(map[uint64]string),
		hashIndices:  make([]uint64, 0),
		hashFunction: fnv.New64a(),
	}
}

func (platform *PlatformData) LogCurrentRows() {
	for rowIndex, row := range platform.rows {
		log.Debug().
			Str("Row", string(row)).
			Int("RowIndex", rowIndex).
			Send()
	}
}

func (platform *PlatformData) convertRowsToString() string {
	return strings.Join(convertByteArraysToStrings(platform.rows), "\n")
}

func (platform *PlatformData) RollNorth() {
	log.Trace().Msg("RollNorth")
	for colIndex := 0; colIndex < len(platform.rows[0]); colIndex += 1 {
		stoppingIndex := -1
		for rowIndex := 0; rowIndex < len(platform.rows); rowIndex += 1 {
			currentState := platform.rows[rowIndex][colIndex]
			// log.Trace().
			// 	Interface("Coordinates", []int{rowIndex, colIndex}).
			// 	Str("CurrentState", string(currentState)).
			// 	Send()
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
				platform.rows[rowIndex][colIndex] = EMPTY_SPACE

				// The rounded rock in this position will roll away to the stopping index
				platform.rows[stoppingIndex+1][colIndex] = ROUNDED_ROCK

				// The stoppingIndex of this column is now the roundedRock we just updated
				stoppingIndex = stoppingIndex + 1
			default:
				log.Fatal().Msgf("Encountered unknown character at (%v, %v): %v", colIndex, rowIndex, currentState)
			}
		}
	}
}

func (platform *PlatformData) RollSouth() {
	log.Trace().Msg("RollSouth")
	for colIndex := 0; colIndex < len(platform.rows[0]); colIndex += 1 {
		stoppingIndex := len(platform.rows)
		for rowIndex := len(platform.rows) - 1; rowIndex >= 0; rowIndex -= 1 {
			currentState := platform.rows[rowIndex][colIndex]
			// log.Trace().
			// 	Interface("Coordinates", []int{rowIndex, colIndex}).
			// 	Str("CurrentState", string(currentState)).
			// 	Int("StoppingIndex", stoppingIndex).
			// 	Send()
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
				platform.rows[rowIndex][colIndex] = EMPTY_SPACE

				// The rounded rock in this position will roll away to the stopping index
				platform.rows[stoppingIndex-1][colIndex] = ROUNDED_ROCK

				// The stoppingIndex of this column is now the roundedRock we just updated
				stoppingIndex = stoppingIndex - 1
			default:
				log.Fatal().Msgf("Encountered unknown character at (%v, %v): %v", colIndex, rowIndex, currentState)
			}
		}
	}
}

func (platform *PlatformData) RollWest() {
	log.Trace().Msg("RollWest")
	for rowIndex := 0; rowIndex < len(platform.rows); rowIndex += 1 {
		stoppingIndex := -1
		for colIndex := 0; colIndex < len(platform.rows[0]); colIndex += 1 {
			currentState := platform.rows[rowIndex][colIndex]
			// log.Trace().
			// 	Interface("Coordinates", []int{rowIndex, colIndex}).
			// 	Str("CurrentState", string(currentState)).
			// 	Send()
			switch currentState {
			case EMPTY_SPACE:
				// If the current space is empty, we need not update anything and can simply continue
				continue
			case CUBE_ROCK:
				// If we encounter a new cube rock, we must update the stoppingIndex
				// Such that the next rounded rock will only roll to this position
				stoppingIndex = colIndex
			case ROUNDED_ROCK:
				// If we encounter a rounded rock that can roll, it's time to update things!
				// The rock will roll to the stoppingIndex, as tracked above.
				// Importantly, the stoppingIndex is the index of the *block*,
				// so we only roll to one before it

				// This position is now empty, as the rock will roll
				platform.rows[rowIndex][colIndex] = EMPTY_SPACE

				// The rounded rock in this position will roll away to the stopping index
				platform.rows[stoppingIndex+1][colIndex] = ROUNDED_ROCK

				// The stoppingIndex of this column is now the roundedRock we just updated
				stoppingIndex = stoppingIndex + 1
			default:
				log.Fatal().Msgf("Encountered unknown character at (%v, %v): %v", colIndex, rowIndex, currentState)
			}
		}
	}
}

func (platform *PlatformData) RollEast() {
	log.Trace().Msg("RollEast")
	for rowIndex := 0; rowIndex < len(platform.rows); rowIndex += 1 {
		stoppingIndex := len(platform.rows)
		for colIndex := len(platform.rows[0]) - 1; colIndex >= 0; colIndex -= 1 {
			currentState := platform.rows[rowIndex][colIndex]
			// log.Trace().
			// 	Interface("Coordinates", []int{rowIndex, colIndex}).
			// 	Str("CurrentState", string(currentState)).
			// 	Send()
			switch currentState {
			case EMPTY_SPACE:
				// If the current space is empty, we need not update anything and can simply continue
				continue
			case CUBE_ROCK:
				// If we encounter a new cube rock, we must update the stoppingIndex
				// Such that the next rounded rock will only roll to this position
				stoppingIndex = colIndex
			case ROUNDED_ROCK:
				// If we encounter a rounded rock that can roll, it's time to update things!
				// The rock will roll to the stoppingIndex, as tracked above.
				// Importantly, the stoppingIndex is the index of the *block*,
				// so we only roll to one before it

				// This position is now empty, as the rock will roll
				platform.rows[rowIndex][colIndex] = EMPTY_SPACE

				// The rounded rock in this position will roll away to the stopping index
				platform.rows[stoppingIndex-1][colIndex] = ROUNDED_ROCK

				// The stoppingIndex of this column is now the roundedRock we just updated
				stoppingIndex = stoppingIndex - 1
			default:
				log.Fatal().Msgf("Encountered unknown character at (%v, %v): %v", colIndex, rowIndex, currentState)
			}
		}
	}
}

func (platform *PlatformData) addToCache(stateBeforeHash uint64, stateAfter string) {
	platform.cache[stateBeforeHash] = stateAfter
	platform.hashIndices = append(platform.hashIndices, stateBeforeHash)

}

func (platform *PlatformData) PerformCycles(numberOfCycles int) {

	cycleIndex := 0
	for ; cycleIndex < numberOfCycles; cycleIndex += 1 {

		stateBefore := platform.convertRowsToString()
		platform.hashFunction.Reset()
		platform.hashFunction.Write([]byte(stateBefore))
		stateBeforeHash := platform.hashFunction.Sum64()

		log.Debug().
			Int("CycleIndex", cycleIndex).
			Int("Hash", int(stateBeforeHash)).
			Send()

		if _, ok := platform.cache[stateBeforeHash]; ok {
			log.Debug().
				Int("TotalCycles", len(platform.hashIndices)).
				Int("CurrentStateIndex", slices.Index(platform.hashIndices, stateBeforeHash)).
				Msg("CACHE HIT")
			// resultAsByteArrays := convertStringsToByteArrays(strings.Split(result, "\n"))
			// platform.rows = resultAsByteArrays
			break
		}

		platform.RollNorth()
		platform.RollWest()
		platform.RollSouth()
		platform.RollEast()
		stateAfter := platform.convertRowsToString()
		platform.addToCache(stateBeforeHash, stateAfter)
	}

	stateBefore := platform.convertRowsToString()
	platform.hashFunction.Reset()
	platform.hashFunction.Write([]byte(stateBefore))
	stateBeforeHash := platform.hashFunction.Sum64()

	loopStartIndex := slices.Index(platform.hashIndices, stateBeforeHash)
	loopLength := (len(platform.hashIndices) - loopStartIndex)
	numberOfCompleteLoops := (numberOfCycles - loopStartIndex) / loopLength
	additionalIterations := (numberOfCycles - loopStartIndex) - loopLength*numberOfCompleteLoops

	log.Debug().
		Int("NumberOfCycles", numberOfCycles).
		Int("LoopLength", loopLength).
		Int("NumberOfCompleteLoops", numberOfCompleteLoops).
		Int("AdditionalIterations", additionalIterations).
		Send()

	for i := 0; i < additionalIterations; i += 1 {
		stateBefore := platform.convertRowsToString()
		platform.hashFunction.Reset()
		platform.hashFunction.Write([]byte(stateBefore))
		stateBeforeHash := platform.hashFunction.Sum64()

		platform.rows = convertStringsToByteArrays(strings.Split(platform.cache[stateBeforeHash], "\n"))
	}
}

func (platform *PlatformData) CalculateLoad() int {
	totalLoad := 0
	for rowIndex := 0; rowIndex < len(platform.rows); rowIndex += 1 {
		for colIndex := 0; colIndex < len(platform.rows[0]); colIndex += 1 {
			currentState := platform.rows[rowIndex][colIndex]
			if currentState == ROUNDED_ROCK {
				totalLoad += len(platform.rows) - rowIndex
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

	platform.PerformCycles(1_000_000_000)

	// platform.LogCurrentRows()
	totalLoad := platform.CalculateLoad()

	return totalLoad, nil
}

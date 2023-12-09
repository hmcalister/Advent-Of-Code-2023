package part01

import (
	"bufio"
	"slices"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

type historyData struct {
	HistoricValues []int
}

func parseLineToHistoryData(line string) historyData {
	fields := strings.Fields(line)
	history := historyData{
		HistoricValues: make([]int, len(fields)),
	}

	var err error
	for i, f := range fields {
		history.HistoricValues[i], err = strconv.Atoi(f)
		if err != nil {
			log.Panic().Msgf("failed to parse history value %v at index %v in line %v", f, i, line)
		}
	}

	return history
}

func (history historyData) findNextValueInSequence() int {
	var currentSequence []int
	var nextSequence []int
	var nextSequenceAllZero bool
	lastItemInSequences := make([]int, 0)

	currentSequence = history.HistoricValues

	// While we do not have a zero sequence
	nextSequenceAllZero = false
	for !nextSequenceAllZero {
		nextSequence = make([]int, len(currentSequence)-1)
		nextSequenceAllZero = true
		for i := 0; i < len(currentSequence)-1; i += 1 {
			nextSequence[i] = currentSequence[i+1] - currentSequence[i]
			if nextSequence[i] != 0 {
				nextSequenceAllZero = false
			}
		}
		log.Trace().
			Bool("NextSequenceAllZeros", nextSequenceAllZero).
			Interface("NextSequence", nextSequence).
			Send()

		lastItemInSequences = append(lastItemInSequences, currentSequence[len(currentSequence)-1])
		currentSequence = nextSequence
	}

	log.Trace().
		Interface("LastItemInSequences", lastItemInSequences).
		Send()

	nextTerm := 0
	slices.Reverse(lastItemInSequences)
	for _, s := range lastItemInSequences {
		log.Trace().
			Int("CurrentNextTerm", nextTerm).
			Int("NextLastItem", s).
			Int("NextTerm", nextTerm+s).
			Send()
		nextTerm += s
	}

	return nextTerm
}

func ProcessInput(fileScanner *bufio.Scanner) (int, error) {
	result := 0

	for fileScanner.Scan() {
		line := fileScanner.Text()
		history := parseLineToHistoryData(line)
		log.Debug().
			Str("ParsedLine", line).
			Interface("ParsedHistory", history).
			Send()

		result += history.findNextValueInSequence()
	}

	return result, nil
}

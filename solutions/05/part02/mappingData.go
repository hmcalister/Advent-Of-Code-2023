package part02

import (
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

// Structure to record the source range, given as [rangeStart:rangeEnd)
// to destination range, given as [rangeStart+offset:rangeEnd+offset)
// of a map.
type mappingData struct {
	rangeStart int
	rangeEnd   int
	offset     int
}

// Parse a single line of puzzleInput to a mapping data struct.
//
// Will error if line is not exactly three whitespace separated numbers.
func parseLineToMappingData(line string) *mappingData {
	var err error
	parts := strings.Fields(line)
	destinationRangeStart, err := strconv.Atoi(parts[0])
	if err != nil {
		log.Fatal().Msgf("error parsing destinationRange:%v", err)
	}

	sourceRangeStart, err := strconv.Atoi(parts[1])
	if err != nil {
		log.Fatal().Msgf("error parsing sourceRangeStart:%v", err)
	}

	rangeLength, err := strconv.Atoi(parts[2])
	if err != nil {
		log.Fatal().Msgf("error parsing rangeLen:%v", err)
	}

	md := &mappingData{
		rangeStart: sourceRangeStart,
		rangeEnd:   sourceRangeStart + rangeLength,
		offset:     destinationRangeStart - sourceRangeStart,
	}

	// log.Trace().
	// 	Str("InitString", line).
	// 	Interface("ParsedMapRangeStart", md.rangeStart).
	// 	Interface("ParsedMapRangeEnd", md.rangeEnd).
	// 	Interface("ParsedMapOffset", md.offset).
	// 	Send()
	return md
}

// Checks if the given testValue is in this mappingData's range.
func (md *mappingData) IsInRange(testValue int) bool {
	return md.rangeStart <= testValue && testValue < md.rangeEnd
}

// Maps the given value according to the mapping.
//
// Note, ensure the given value is actually in the mappingData range! (See md.IsInRange)
func (md *mappingData) MapValue(initialValue int) int {
	return initialValue + md.offset
}

package lib

import (
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

// Details about a mapping from one domain to another
//
// Mappings that the form of a+offset if a is in [start, end)
type mapData struct {
	sourceStart int
	sourceEnd   int
	offset      int
}

func (md mapData) ValueInMapSource(value int) bool {
	return md.sourceStart <= value && value < md.sourceEnd
}

func (md mapData) ValueInMapDest(value int) bool {
	return md.MapValue(md.sourceStart) <= value && value < md.MapValue(md.sourceEnd)
}

func (md mapData) MapValue(value int) int {
	return value + md.offset
}

func (md mapData) InverseMapValue(value int) int {
	return value - md.offset
}

// Line must be of the form of three ints separated by whitespace
//
// The order of the ints are destinationStart, sourceStart, RangeLength
func parseLineToMapping(line string) mapData {
	fields := strings.Fields(line)

	destinationStartStr := fields[0]
	destinationStart, err := strconv.Atoi(destinationStartStr)
	if err != nil {
		log.Panic().Msgf("failed to parse destination start %v (line %v)", destinationStartStr, line)
	}

	sourceStartStr := fields[1]
	sourceStart, err := strconv.Atoi(sourceStartStr)
	if err != nil {
		log.Panic().Msgf("failed to parse source start %v (line %v)", sourceStartStr, line)
	}

	rangeLenStr := fields[2]
	rangeLen, err := strconv.Atoi(rangeLenStr)
	if err != nil {
		log.Panic().Msgf("failed to parse range length %v (line %v)", rangeLenStr, line)
	}

	return mapData{
		sourceStart: sourceStart,
		sourceEnd:   sourceStart + rangeLen,
		offset:      destinationStart - sourceStart,
	}
}

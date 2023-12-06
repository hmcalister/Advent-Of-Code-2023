package part02

import (
	"hmcalister/aoc05/lib"
	"math"

	"github.com/rs/zerolog/log"
)

type seedRangeData struct {
	SeedRangeStart  int
	SeedRangeLength int
}

func checkSeedRange(seedRange seedRangeData, domainMapper lib.DomainMapper) int {
	log.Info().Interface("CheckingSeedRange", seedRange).Send()

	// We only need to check the beginning of each map, since the maps themselves are monotonically increasing.

	currentSeedValue := seedRange.SeedRangeStart

	// Find the map that contains the current seed

	domainMapperRangeStarts := domainMapper.GetAllRangeStarts()

	mapIndex := 0
	for domainMapperRangeStarts[mapIndex] < currentSeedValue {
		mapIndex += 1
	}
	mapIndex -= 1

	// Now just check maps until the start is beyond our target range

	minSeedVal := math.MaxInt
	for currentSeedValue < seedRange.SeedRangeStart+seedRange.SeedRangeLength {
		mappedSeedValue := domainMapper.MapValue(currentSeedValue)
		if mappedSeedValue < minSeedVal {
			minSeedVal = mappedSeedValue
		}

		log.Debug().
			Int("CheckedSeed", currentSeedValue).
			Int("MapIndex", mapIndex).
			Int("MapFirstValue", domainMapperRangeStarts[mapIndex]).
			Int("MapLastValue", domainMapperRangeStarts[mapIndex+1]).
			Int("MappedSeedValue", mappedSeedValue).
			Send()

		mapIndex += 1
		currentSeedValue = domainMapperRangeStarts[mapIndex]

	}

	return minSeedVal
}

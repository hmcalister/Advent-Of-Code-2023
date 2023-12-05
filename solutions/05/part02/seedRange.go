package part02

import (
	"math"

	"github.com/rs/zerolog/log"
)

type seedRangeData struct {
	SeedRangeStart  int
	SeedRangeLength int
}

func checkSeedRange(seedRange *seedRangeData, domainMappersArray []*domainMapper, mapBoundaries []int) int {
	log.Info().Interface("CheckingSeedRange", seedRange).Send()

	// We only need to check the beginning of each map, since the maps themselves are monotonically increasing.

	currentSeedValue := seedRange.SeedRangeStart

	// Find the map that contains the current seed

	mapIndex := 0
	for mapBoundaries[mapIndex] < currentSeedValue {
		mapIndex += 1
	}

	// Now just check maps until the start is beyond our target range

	minSeedVal := math.MaxInt
	for currentSeedValue < seedRange.SeedRangeStart+seedRange.SeedRangeLength {
		mappedSeedValue := currentSeedValue
		for _, dm := range domainMappersArray {
			mappedSeedValue = dm.MapValue(mappedSeedValue)
		}
		if mappedSeedValue < minSeedVal {
			minSeedVal = mappedSeedValue
		}

		log.Debug().
			Int("CheckedSeed", currentSeedValue).
			Int("MapIndex", mapIndex).
			Int("MapFirstValue", mapBoundaries[mapIndex]).
			Int("MapLastValue", mapBoundaries[mapIndex+1]).
			Int("MappedSeedValue", mappedSeedValue).
			Send()

		mapIndex += 1
		currentSeedValue = mapBoundaries[mapIndex]

	}

	return minSeedVal
}

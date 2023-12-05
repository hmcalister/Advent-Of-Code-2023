package part02

import (
	"math"

	"github.com/rs/zerolog/log"
)

type seedRangeData struct {
	SeedRangeStart  int
	SeedRangeLength int
}

func checkSeedRange(seedRange *seedRangeData, composedDomainMapper *domainMapper) int {
	log.Info().Interface("CheckingSeedRange", seedRange).Send()

	// We only need to check the beginning of each map, since the maps themselves are monotonically increasing.

	currentSeedValue := seedRange.SeedRangeStart

	// Find the map that contains the current seed

	mapIndex := 0
	targetMap := composedDomainMapper.MapDataArray[mapIndex]
	for targetMap.sourceRangeStart+targetMap.rangeLength < currentSeedValue {
		mapIndex += 1
		targetMap = composedDomainMapper.MapDataArray[mapIndex]
	}

	// Now just check maps until the start is beyond our target range

	minSeedVal := math.MaxInt
	for currentSeedValue < seedRange.SeedRangeStart+seedRange.SeedRangeLength {
		mappedSeedValue := composedDomainMapper.MapValue(currentSeedValue)
		if mappedSeedValue < minSeedVal {
			minSeedVal = mappedSeedValue
		}
		log.Debug().
			Int("CheckedSeed", currentSeedValue).
			Int("MapIndex", mapIndex).
			Int("MapFirstValue", composedDomainMapper.MapDataArray[mapIndex].sourceRangeStart).
			Int("MapLastValue", composedDomainMapper.MapDataArray[mapIndex].sourceRangeStart+composedDomainMapper.MapDataArray[mapIndex].rangeLength).
			Int("MappedSeedValue", mappedSeedValue).
			Send()

		mapIndex += 1
		currentSeedValue = composedDomainMapper.MapDataArray[mapIndex].sourceRangeStart

	}

	return minSeedVal
}

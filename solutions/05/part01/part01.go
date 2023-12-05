package part01

import (
	"bufio"
	"math"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

func ProcessInput(fileScanner *bufio.Scanner) (int, error) {
	// Handle seeds
	fileScanner.Scan()
	seedLine := fileScanner.Text()
	seedLine = seedLine[7:]
	seedValuesStrs := strings.Fields(seedLine)
	seedValues := make([]int, len(seedValuesStrs))
	for i, seedValueStr := range seedValuesStrs {
		val, err := strconv.Atoi(seedValueStr)
		if err != nil {
			log.Fatal().Msgf("error parsing seed value:%v", err)
		}
		seedValues[i] = val
		log.Debug().
			Int("SeedValueIndex", i).
			Str("SeedValueStr", seedValueStr).
			Int("SeedValue", val).
			Send()
	}
	fileScanner.Scan()

	allDomainMappers := parseFileToDomainMappers(fileScanner)
	log.Debug().Int("NumDomainMappers", len(allDomainMappers)).Send()

	minSeedVal := math.MaxInt
	// Feed each seed through the maps and see where they end up
	for i, seed := range seedValues {
		seedMappedValue := seed
		for _, dm := range allDomainMappers {
			seedMappedValue = dm.MapValue(seedMappedValue)
		}
		log.Debug().
			Int("SeedValueIndex", i).
			Int("SeedValue", seed).
			Int("MappedSeedValue", seedMappedValue).
			Send()

		if seedMappedValue < minSeedVal {
			minSeedVal = seedMappedValue
			log.Debug().Msgf("new lowest seed found: %v (%v)", i, seedMappedValue)
		}
	}

	return minSeedVal, nil
}

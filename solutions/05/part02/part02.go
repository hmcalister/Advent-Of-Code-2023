package part02

import (
	"bufio"
	"hmcalister/aoc05/lib"
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
	fileScanner.Scan()

	allDomainMappers := lib.GetIdentityMapper()
	for fileScanner.Scan() {
		allDomainMappers = lib.ComposeDomainMappers(allDomainMappers, lib.ParseSectionToDomainMapper(fileScanner))
	}

	minMappedValue := math.MaxInt
	for i := 0; i < len(seedValuesStrs); i += 2 {
		seedValueStartStr := seedValuesStrs[i]
		seedValueRangeStr := seedValuesStrs[i+1]

		seedValueStart, err := strconv.Atoi(seedValueStartStr)
		if err != nil {
			log.Fatal().Msgf("error parsing seed value start:%v", err)
		}
		seedValueRange, err := strconv.Atoi(seedValueRangeStr)
		if err != nil {
			log.Fatal().Msgf("error parsing seed value start:%v", err)
		}

		rangeValue := checkSeedRange(seedRangeData{
			SeedRangeStart:  seedValueStart,
			SeedRangeLength: seedValueRange,
		}, allDomainMappers)
		if rangeValue < minMappedValue {
			minMappedValue = rangeValue
			log.Info().Msgf("New best location found: %v", minMappedValue)
		}
	}

	return minMappedValue, nil
}

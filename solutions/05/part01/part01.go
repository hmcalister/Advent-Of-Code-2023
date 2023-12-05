package part01

import (
	"bufio"
	"math"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

type mappingData struct {
	rangeLength      int
	sourceRangeStart int
	destRangeStart   int
}

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
		rangeLength,
		sourceRangeStart,
		destinationRangeStart,
	}

	log.Trace().
		Str("InitString", line).
		Int("NumFields", len(parts)).
		Int("rangeLength", rangeLength).
		Int("sourceRangeStart", sourceRangeStart).
		Int("destinationRangeStart", destinationRangeStart).
		Send()
	return md
}

func (md *mappingData) IsInRange(testValue int) bool {
	return md.sourceRangeStart <= testValue && testValue < md.sourceRangeStart+md.rangeLength
}

func (md *mappingData) MapValue(initialValue int) int {
	return initialValue + (md.destRangeStart - md.sourceRangeStart)
}

type domainMapper struct {
	domainMappingID string
	MapDataArray    []*mappingData
}

func (dm *domainMapper) MapValue(val int) int {
	for mdID, md := range dm.MapDataArray {
		if md.IsInRange(val) {
			log.Trace().Int("ApplyingMapID", mdID).Send()
			return md.MapValue(val)
		}
	}
	return val
}

type endToEndMapper struct {
	domainMappers []*domainMapper
}

func (e2em *endToEndMapper) MapValue(val int) int {
	log.Trace().Int("InitValue", val).Send()
	for _, dm := range e2em.domainMappers {
		val = dm.MapValue(val)
		log.Trace().
			Str("ApplyingDomainID", dm.domainMappingID).
			Int("MappedValue", val).
			Send()
	}
	return val
}

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

	consecutiveMaps := make([]*domainMapper, 0)
	// Parse each mapping in turn
	for fileScanner.Scan() {
		// Get the initial line
		domainMappingID := fileScanner.Text()
		log.Debug().Str("MappingID", domainMappingID).Send()

		currentMapData := make([]*mappingData, 0)
		for fileScanner.Scan() {
			line := fileScanner.Text()
			if len(line) == 0 {
				break
			}
			log.Trace().Str("NextLine", line).Send()
			currentMapData = append(currentMapData, parseLineToMappingData(line))
		}
		consecutiveMaps = append(consecutiveMaps, &domainMapper{domainMappingID: domainMappingID, MapDataArray: currentMapData})
	}
	fullMapper := &endToEndMapper{domainMappers: consecutiveMaps}

	minSeedVal := math.MaxInt
	// Feed each seed through the maps and see where they end up
	for i, seed := range seedValues {
		log.Debug().
			Int("SeedValueIndex", i).
			Int("SeedValue", seed).
			Send()
		seedMappedValue := fullMapper.MapValue(seed)
		if seedMappedValue < minSeedVal {
			minSeedVal = seedMappedValue
			log.Debug().Msgf("new lowest seed found: %v (%v)", i, seedMappedValue)
		}
	}

	return minSeedVal, nil
}

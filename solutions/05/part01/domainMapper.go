package part01

import (
	"bufio"
	"math"
	"slices"
	"sort"

	"github.com/rs/zerolog/log"
)

// A collection of maps detailing the entire transformation from one domain to the next.
type domainMapper struct {
	MapDataArray []*mappingData
	boundaries   []int
}

// Map a value from the source domain to the destination domain.
func (dm *domainMapper) MapValue(val int) int {
	for _, md := range dm.MapDataArray {
		if md.IsInRange(val) {
			return md.MapValue(val)
		}
	}
	return val
}

func parseSectionToDomainMapper(fileScanner *bufio.Scanner) *domainMapper {
	// Get the initial line, with the domain mapping ID
	domainMappingID := fileScanner.Text()
	log.Debug().Str("MappingID", domainMappingID).Send()

	currentMapData := make([]*mappingData, 0)
	// The remaining lines  (up to an empty line) are mapData
	for fileScanner.Scan() {
		line := fileScanner.Text()
		if len(line) == 0 {
			break
		}
		currentMapData = append(currentMapData, parseLineToMappingData(line))
	}
	newDomainManager := &domainMapper{MapDataArray: currentMapData}
	sort.Slice(newDomainManager.MapDataArray, func(i, j int) bool {
		return newDomainManager.MapDataArray[i].rangeStart < newDomainManager.MapDataArray[j].rangeStart
	})

	boundaries := make([]int, 0)
	boundaries = append(boundaries, 0)
	for _, mapData := range newDomainManager.MapDataArray {
		boundaries = append(boundaries, mapData.rangeStart)
		boundaries = append(boundaries, mapData.rangeEnd)
	}
	boundaries = append(boundaries, math.MaxInt)
	newDomainManager.boundaries = slices.Compact(boundaries)

	return newDomainManager
}

func parseFileToAllDomainMappers(fileScanner *bufio.Scanner) []*domainMapper {
	// We start with a domain mapper that performs the identity mapping.
	// As we parse each domain mapper, we will compose with this mapper.

	allDomainMappers := make([]*domainMapper, 0)
	for fileScanner.Scan() {
		allDomainMappers = append(allDomainMappers, parseSectionToDomainMapper(fileScanner))
	}
	return allDomainMappers
}

// Given two domain mappers, compose them and return a single domain mapper.
//
// Note this method assumes the global mappings have mapA feed directly into mapB.
func composeDomainMappers(domainAMapper *domainMapper, domainBMapper *domainMapper) *domainMapper {
	// The idea here is to check each map of mapA (sorted by start) in turn,
	// and split these up by whatever mapB does to that range.

	// Generate all of the boundaries, together!
	composedBoundaries := make([]int, len(domainAMapper.boundaries)+len(domainBMapper.boundaries))
	copy(composedBoundaries, domainAMapper.boundaries)
	copy(composedBoundaries[len(domainAMapper.boundaries):], domainBMapper.boundaries)
	slices.Sort(composedBoundaries)
	composedBoundaries = slices.Compact(composedBoundaries)
	log.Trace().
		Interface("DomainAMapperBoundaries", domainAMapper.boundaries).
		Interface("DomainBMapperBoundaries", domainBMapper.boundaries).
		Interface("ComposedBoundaries", composedBoundaries).
		Send()

	composedMaps := make([]*mappingData, 0)
	for boundaryIndex := 0; boundaryIndex < len(composedBoundaries)-1; boundaryIndex += 1 {
		mapStart := composedBoundaries[boundaryIndex]
		mapEnd := composedBoundaries[boundaryIndex+1]
		mapOffset := domainBMapper.MapValue(domainAMapper.MapValue(mapStart)) - mapStart

		if len(composedMaps) > 0 && composedMaps[len(composedMaps)-1].offset == mapOffset {
			composedMaps[len(composedMaps)-1].rangeEnd = mapEnd
			continue
		}

		composedMaps = append(composedMaps, &mappingData{
			rangeStart: mapStart,
			rangeEnd:   mapEnd,
			offset:     mapOffset,
		})
	}

	return &domainMapper{
		MapDataArray: composedMaps,
		boundaries:   composedBoundaries,
	}
}

func parseFileToComposedMapper(fileScanner *bufio.Scanner) *domainMapper {
	// We start with a domain mapper that performs the identity mapping.
	// As we parse each domain mapper, we will compose with this mapper.

	composedDomainMapper := &domainMapper{
		MapDataArray: []*mappingData{{
			rangeStart: 0,
			rangeEnd:   math.MaxInt,
			offset:     0,
		}},
	}

	for fileScanner.Scan() {
		nextDomainMapper := parseSectionToDomainMapper(fileScanner)
		composedDomainMapper = composeDomainMappers(composedDomainMapper, nextDomainMapper)
	}
	for i, mapData := range composedDomainMapper.MapDataArray {
		log.Debug().
			Int("MapIndex", i).
			Interface("MapRangeStart", mapData.rangeStart).
			Interface("MapRangeEnd", mapData.rangeEnd).
			Interface("MapOffset", mapData.offset).
			Send()
	}

	return composedDomainMapper
}

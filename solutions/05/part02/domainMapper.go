package part02

import (
	"bufio"
	"math"
	"slices"
	"sort"

	"github.com/rs/zerolog/log"
)

// A collection of maps detailing the entire transformation from one domain to the next.
type domainMapper struct {
	MapDataArray  []*mappingData
	mapBoundaries []int
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
		log.Trace().Str("NextLine", line).Send()
		currentMapData = append(currentMapData, parseLineToMappingData(line))
	}
	newDomainManager := &domainMapper{MapDataArray: currentMapData}
	sort.Slice(newDomainManager.MapDataArray, func(i, j int) bool {
		return newDomainManager.MapDataArray[i].sourceRangeStart < newDomainManager.MapDataArray[j].sourceRangeStart
	})

	boundaries := make([]int, 0)
	boundaries = append(boundaries, 0)
	for _, mapData := range newDomainManager.MapDataArray {
		boundaries = append(boundaries, mapData.sourceRangeStart)
		boundaries = append(boundaries, mapData.sourceRangeStart+mapData.rangeLength)
	}
	boundaries = append(boundaries, math.MaxInt)
	newDomainManager.mapBoundaries = slices.Compact(boundaries)

	return newDomainManager
}

// Given two domain mappers, compose them and return a single domain mapper.
//
// Note this method assumes the global mappings have mapA feed directly into mapB.
func composeDomainMappers(domainAMapper *domainMapper, domainBMapper *domainMapper) *domainMapper {
	// The idea here is to check each map of mapA (sorted by start) in turn,
	// and split these up by whatever mapB does to that range.

	// Generate all of the boundaries, together!
	composedBoundaries := make([]int, len(domainAMapper.mapBoundaries)+len(domainBMapper.mapBoundaries))
	copy(composedBoundaries, domainAMapper.mapBoundaries)
	copy(composedBoundaries[len(domainAMapper.mapBoundaries):], domainBMapper.mapBoundaries)
	slices.Sort(composedBoundaries)
	composedBoundaries = slices.Compact(composedBoundaries)

	composedMaps := make([]*mappingData, 0)
	for boundaryIndex := 0; boundaryIndex < len(composedBoundaries)-1; boundaryIndex += 1 {
		mapStart := composedBoundaries[boundaryIndex]
		mapLength := composedBoundaries[boundaryIndex+1] - mapStart
		mapOffset := domainBMapper.MapValue(domainAMapper.MapValue(mapStart)) - mapStart

		if len(composedMaps) > 0 && composedMaps[len(composedMaps)-1].offset == mapOffset {
			composedMaps[len(composedMaps)-1].rangeLength += mapLength
			continue
		}

		composedMaps = append(composedMaps, &mappingData{
			sourceRangeStart: mapStart,
			rangeLength:      mapLength,
			offset:           mapOffset,
		})
	}

	log.Debug().Interface("ComposedMapBoundaries", composedBoundaries).Send()

	return &domainMapper{
		MapDataArray:  composedMaps,
		mapBoundaries: composedBoundaries,
	}
}

func parseFileToComposedMapper(fileScanner *bufio.Scanner) *domainMapper {
	// We start with a domain mapper that performs the identity mapping.
	// As we parse each domain mapper, we will compose with this mapper.

	composedDomainMapper := &domainMapper{
		MapDataArray: []*mappingData{{
			sourceRangeStart: 0,
			rangeLength:      math.MaxInt,
			offset:           0,
		}},
	}

	for fileScanner.Scan() {
		nextDomainMapper := parseSectionToDomainMapper(fileScanner)
		composedDomainMapper = composeDomainMappers(composedDomainMapper, nextDomainMapper)
	}
	for i, mapData := range composedDomainMapper.MapDataArray {
		log.Debug().
			Int("MapIndex", i).
			Interface("MapSourceRangeStart", mapData.sourceRangeStart).
			Interface("MapRangeLength", mapData.rangeLength).
			Interface("MapOffset", mapData.offset).
			Send()
	}

	return composedDomainMapper
}

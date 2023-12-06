package lib

import (
	"bufio"
	"math"
	"slices"
	"sort"

	"github.com/rs/zerolog/log"
)

type DomainMapper struct {
	maps []mapData
}

func GetIdentityMapper() DomainMapper {
	return DomainMapper{
		maps: []mapData{{
			sourceStart: 0,
			sourceEnd:   math.MaxInt,
			offset:      0,
		}},
	}
}

func (dm DomainMapper) MapValue(value int) int {
	for _, m := range dm.maps {
		if m.ValueInMapSource(value) {
			return m.MapValue(value)
		}
	}

	return value
}

func (dm DomainMapper) InverseMapValue(value int) int {
	for _, m := range dm.maps {
		if m.ValueInMapDest(value) {
			return m.InverseMapValue(value)
		}
	}

	return value
}

func (dm DomainMapper) GetAllRangeStarts() []int {
	rangeStarts := make([]int, len(dm.maps))
	for i, m := range dm.maps {
		rangeStarts[i] = m.sourceStart
	}
	return rangeStarts
}

func ParseSectionToDomainMapper(fileScanner *bufio.Scanner) DomainMapper {
	line := fileScanner.Text()
	log.Debug().
		Str("DomainMapperID", line).
		Send()

	maps := make([]mapData, 0)
	for fileScanner.Scan() {
		line = fileScanner.Text()
		if len(line) == 0 {
			break
		}

		log.Trace().
			Str("ParsingLineToMap", line).
			Send()

		maps = append(maps, parseLineToMapping(line))
	}

	sort.Slice(maps, func(i, j int) bool {
		return maps[i].sourceStart < maps[j].sourceStart
	})

	return DomainMapper{maps}
}

func ComposeDomainMappers(domainAMapper DomainMapper, domainBMapper DomainMapper) DomainMapper {
	composedBoundaries := make([]int, 0)
	composedBoundaries = append(composedBoundaries, 0)
	for _, AMap := range domainAMapper.maps {
		composedBoundaries = append(composedBoundaries, AMap.sourceStart)
		composedBoundaries = append(composedBoundaries, AMap.sourceEnd)
	}
	for _, BMap := range domainBMapper.maps {
		composedBoundaries = append(composedBoundaries, domainAMapper.InverseMapValue(BMap.sourceStart))
		composedBoundaries = append(composedBoundaries, domainAMapper.InverseMapValue(BMap.sourceEnd))
	}
	composedBoundaries = append(composedBoundaries, math.MaxInt)
	slices.Sort(composedBoundaries)
	composedBoundaries = slices.Compact(composedBoundaries)

	log.Trace().
		Interface("ComposedBoundaries", composedBoundaries).
		Send()

	composedMaps := make([]mapData, 0)
	for i := 0; i < len(composedBoundaries)-1; i += 1 {
		sourceStart := composedBoundaries[i]
		sourceEnd := composedBoundaries[i+1]
		offset := domainBMapper.MapValue(domainAMapper.MapValue(sourceStart)) - sourceStart
		log.Trace().
			Int("BoundaryIndex", i).
			Int("SourceStart", sourceStart).
			Int("SourceEnd", sourceEnd).
			Int("Offset", offset).
			Send()

		if len(composedMaps) > 0 && composedMaps[len(composedMaps)-1].offset == offset {
			log.Trace().Msgf("combining maps with same offset %v (redundant map sourceStart %v)", offset, sourceStart)
			composedMaps[len(composedMaps)-1].sourceEnd = sourceEnd
			continue
		}

		composedMaps = append(composedMaps, mapData{
			sourceStart,
			sourceEnd,
			offset,
		})
	}

	return DomainMapper{maps: composedMaps}
}

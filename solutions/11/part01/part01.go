package part01

import (
	"bufio"

	"github.com/rs/zerolog/log"
)

const (
	EMPTY_SPACE_RUNE rune = '.'
	GALAXY_RUNE      rune = '#'
)

type galaxyData struct {
	GalaxyID    int
	XCoordinate int
	YCoordinate int
}

type cosmologicalMapData struct {
	GalaxiesByRow    []int
	GalaxiesByColumn []int
	Galaxies         []galaxyData
}

func newCosmologicalMap() *cosmologicalMapData {
	return &cosmologicalMapData{
		GalaxiesByRow:    make([]int, 0),
		GalaxiesByColumn: make([]int, 0),
	}
}

func (cosmologicalMap *cosmologicalMapData) addNewRow(line string) {
	log.Debug().
		Str("ParsingNextLine", line).
		Send()

	numCols := len(line)
	if len(cosmologicalMap.GalaxiesByColumn) < numCols {
		log.Debug().
			Int("NewLineLength", numCols).
			Int("CurrentMapByColumnsLen", len(cosmologicalMap.GalaxiesByColumn)).
			Msg("Resizing")
		newByColumnArr := make([]int, numCols)
		copy(newByColumnArr, cosmologicalMap.GalaxiesByColumn)
		cosmologicalMap.GalaxiesByColumn = newByColumnArr
	}

	totalGalaxiesInRow := 0
	for i, r := range line {
		if r == GALAXY_RUNE {
			galaxy := galaxyData{
				GalaxyID:    len(cosmologicalMap.Galaxies),
				XCoordinate: i,
				YCoordinate: len(cosmologicalMap.GalaxiesByRow),
			}
			cosmologicalMap.Galaxies = append(cosmologicalMap.Galaxies, galaxy)
			cosmologicalMap.GalaxiesByColumn[i] += 1
			totalGalaxiesInRow += 1

			log.Debug().
				Interface("FoundGalaxy", galaxy).
				Int("GalaxiesByColumnCount", cosmologicalMap.GalaxiesByColumn[i]).
				Int("GalaxiesByRowCount", totalGalaxiesInRow).
				Send()
		}
	}
	cosmologicalMap.GalaxiesByRow = append(cosmologicalMap.GalaxiesByRow, totalGalaxiesInRow)
}

func (cosmologicalMap *cosmologicalMapData) shortestDistanceBetweenPoints(x1, y1, x2, y2 int) int {
	if y1 > y2 {
		y2, y1 = y1, y2
	}
	yDel := y2 - y1
	for y := y1 + 1; y < y2; y += 1 {
		if cosmologicalMap.GalaxiesByRow[y] == 0 {
			yDel += 1
			log.Trace().Msgf("empty row at index %v", y)
		}
	}

	if x1 > x2 {
		x2, x1 = x1, x2
	}
	xDel := x2 - x1
	for x := x1 + 1; x < x2; x += 1 {
		if cosmologicalMap.GalaxiesByColumn[x] == 0 {
			xDel += 1
			log.Trace().Msgf("empty column at index %v", x)
		}
	}

	return yDel + xDel
}

func (cosmologicalMap *cosmologicalMapData) calculateShortestPairwiseDistances() int {
	totalPairwiseDistances := 0
	totalPairs := 0
	for i := 0; i < len(cosmologicalMap.Galaxies)-1; i += 1 {
		galaxyOne := cosmologicalMap.Galaxies[i]
		for j := i + 1; j < len(cosmologicalMap.Galaxies); j += 1 {
			galaxyTwo := cosmologicalMap.Galaxies[j]
			thisPairDistance := cosmologicalMap.shortestDistanceBetweenPoints(galaxyOne.XCoordinate, galaxyOne.YCoordinate, galaxyTwo.XCoordinate, galaxyTwo.YCoordinate)
			totalPairwiseDistances += thisPairDistance
			totalPairs += 1

			log.Debug().
				Int("GalaxyOneIndex", i).
				Int("GalaxyTwoIndex", j).
				Interface("GalaxyOne", galaxyOne).
				Interface("GalaxyTwo", galaxyTwo).
				Int("PairwiseDistance", thisPairDistance).
				Send()
		}
	}

	log.Debug().Int("TotalPairs", totalPairs).Send()

	return totalPairwiseDistances
}

func ProcessInput(fileScanner *bufio.Scanner) (int, error) {
	cosmologicalMap := newCosmologicalMap()
	for fileScanner.Scan() {
		line := fileScanner.Text()
		cosmologicalMap.addNewRow(line)
	}

	totalShortestDistances := cosmologicalMap.calculateShortestPairwiseDistances()

	return totalShortestDistances, nil
}

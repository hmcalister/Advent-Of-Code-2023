package lib

import (
	"bufio"

	"github.com/rs/zerolog/log"
)

type StormData struct {
	hailstoneCollection []HailstoneData
}

func ParseFileToStorm(fileScanner *bufio.Scanner) StormData {
	hailstoneCollection := make([]HailstoneData, 0)

	for fileScanner.Scan() {
		line := fileScanner.Text()
		nextHailstone := parseLineToHailstone(line)

		log.Debug().
			Str("RawLine", line).
			Interface("ParsedHailstone", nextHailstone).
			Msg("ParsingPuzzleInput")

		hailstoneCollection = append(hailstoneCollection, nextHailstone)
	}

	return StormData{
		hailstoneCollection: hailstoneCollection,
	}
}

func (storm StormData) PathIntersectionInXY(minimumPositionBound, maximumPositionBound float64) int {
	validCollisionCount := 0
	for hailstoneOneIndex := 0; hailstoneOneIndex < len(storm.hailstoneCollection); hailstoneOneIndex += 1 {
		hailstoneOne := storm.hailstoneCollection[hailstoneOneIndex]
		for hailstoneTwoIndex := hailstoneOneIndex + 1; hailstoneTwoIndex < len(storm.hailstoneCollection); hailstoneTwoIndex += 1 {
			hailstoneTwo := storm.hailstoneCollection[hailstoneTwoIndex]

			log.Debug().
				Int("HailstoneOneIndex", hailstoneOneIndex).
				Int("HailstoneTwoIndex", hailstoneTwoIndex).
				Msg("PathIntersectionInXY")

			pathIntersection, err := hailstoneOne.FindPathIntersectionPositionInXY(hailstoneTwo)
			if err != nil {
				log.Trace().
					Int("HailstoneOneIndex", hailstoneOneIndex).
					Int("HailstoneTwoIndex", hailstoneTwoIndex).
					Str("HailstoneOne", hailstoneOne.String()).
					Str("HailstoneTwo", hailstoneTwo.String()).
					Str("FoundErr", err.Error()).
					Msg("PathIntersectionInXY")
				continue
			}

			log.Trace().
				Int("HailstoneOneIndex", hailstoneOneIndex).
				Int("HailstoneTwoIndex", hailstoneTwoIndex).
				Str("HailstoneOne", hailstoneOne.String()).
				Str("HailstoneTwo", hailstoneTwo.String()).
				Str("PathIntersection", vectorToString(pathIntersection)).
				Msg("PathIntersectionInXY")

			if minimumPositionBound <= pathIntersection.AtVec(0) && pathIntersection.AtVec(0) <= maximumPositionBound &&
				minimumPositionBound <= pathIntersection.AtVec(1) && pathIntersection.AtVec(1) <= maximumPositionBound {
				log.Debug().Msg("ValidCollisionFound")
				validCollisionCount += 1
			}

		}
	}

	return validCollisionCount
}

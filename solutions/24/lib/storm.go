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

func (storm StormData) DetectCollisionsInXY(minimumPositionBound, maximumPositionBound float64) int {
	validCollisionCount := 0
	for hailstoneOneIndex := 0; hailstoneOneIndex < len(storm.hailstoneCollection); hailstoneOneIndex += 1 {
		hailstoneOne := storm.hailstoneCollection[hailstoneOneIndex]
		for hailstoneTwoIndex := hailstoneOneIndex + 1; hailstoneTwoIndex < len(storm.hailstoneCollection); hailstoneTwoIndex += 1 {
			hailstoneTwo := storm.hailstoneCollection[hailstoneTwoIndex]

			log.Debug().
				Int("HailstoneOneIndex", hailstoneOneIndex).
				Int("HailstoneTwoIndex", hailstoneTwoIndex).
				// Interface("HailstoneOne", hailstoneOne).
				// Interface("HailstoneTwo", hailstoneTwo).
				Msg("DetectCollisionInXY")

			collisionTime, err := hailstoneOne.FindCollisionTimeInXY(hailstoneTwo)
			if err != nil {
				log.Debug().
					Str("FoundErr", err.Error()).
					Msg("DetectCollisionInXY")
				continue
			}

			collisionPosition := hailstoneOne.calculatePositionAtTime(collisionTime)
			log.Debug().
				Float64("CollisionTime", collisionTime).
				Interface("CollisionPosition", collisionPosition).
				Msg("DetectCollisionInXY")

			if minimumPositionBound <= collisionPosition.X && collisionPosition.X <= maximumPositionBound &&
				minimumPositionBound <= collisionPosition.Y && collisionPosition.Y <= maximumPositionBound {
				log.Debug().Msg("ValidCollisionFound")
				validCollisionCount += 1
			}

		}
	}

	return validCollisionCount
}

package part02

import (
	"bufio"
	"hmcalister/aox22/lib"
	"sort"

	"github.com/rs/zerolog/log"
)

const (
	PRESENCE_INDICATOR int = 1
)

func checkSupportsInMap(supports []int, disintegratedMap map[int]int) bool {
	for _, s := range supports {
		if _, ok := disintegratedMap[s]; !ok {
			return false
		}
	}

	return true
}

func ProcessInput(fileScanner *bufio.Scanner) (int, error) {
	pile := lib.ParseFileToBrickPile(fileScanner)
	pile.SimulateBrickFall()

	allBricksByHeight := pile.Bricks

	// Get all the support information organized by the Z height of bricks
	//
	// Note we are sorting allSupports by allBricks.Z
	allSupportsByHeight := pile.Supports
	sort.Slice(allSupportsByHeight, func(i, j int) bool {
		return allBricksByHeight[i].Start.Z < allBricksByHeight[j].Start.Z
	})

	// Get all the bricks in the pile, organized by Z height
	sort.Slice(allBricksByHeight, func(i, j int) bool {
		return allBricksByHeight[i].Start.Z < allBricksByHeight[j].Start.Z
	})

	disintegrationIndicesByBrick := make([]map[int]int, len(allBricksByHeight))

	totalOtherDisintegrations := 0
	for currentBrickIndex := len(allBricksByHeight) - 1; currentBrickIndex >= 0; currentBrickIndex -= 1 {
		currentBrick := allBricksByHeight[currentBrickIndex]
		log.Debug().
			Int("CurrentBrickIndex", currentBrickIndex).
			Str("CurrentBrick", currentBrick.String()).
			Send()

		// Make a new disintegration map for this brick and add the brick to it
		currentBrickDisintegrationMap := make(map[int]int)
		currentBrickDisintegrationMap[currentBrickIndex] = PRESENCE_INDICATOR

		for consideringBrickIndex := currentBrickIndex + 1; consideringBrickIndex < len(allBricksByHeight); consideringBrickIndex += 1 {
			// Before we do anything, just check if the considering brick has already been removed indirectly
			if _, ok := currentBrickDisintegrationMap[consideringBrickIndex]; ok {
				continue
			}

			consideringBrickSupports := allSupportsByHeight[consideringBrickIndex]

			// If the considering brick  will be disintegrated by the current bricks removal, note this
			if checkSupportsInMap(consideringBrickSupports, currentBrickDisintegrationMap) {
				currentBrickDisintegrationMap[consideringBrickIndex] = PRESENCE_INDICATOR

				// We can also add all the items from the considering bricks map, since it will be destroyed
				for indirectDisintegrationIndex := range disintegrationIndicesByBrick[consideringBrickIndex] {
					currentBrickDisintegrationMap[indirectDisintegrationIndex] = PRESENCE_INDICATOR
				}
			}
		}

		log.Debug().Int("DisintegrationCount", len(currentBrickDisintegrationMap)).Send()
		totalOtherDisintegrations += len(currentBrickDisintegrationMap) - 1

		// Add the now filled disintegration map to the full array
		disintegrationIndicesByBrick[currentBrickIndex] = currentBrickDisintegrationMap
	}

	return totalOtherDisintegrations, nil
}

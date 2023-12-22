package part02

import (
	"bufio"
	"hmcalister/aox22/lib"
)

const (
	PRESENCE_INDICATOR bool = true
)

func checkSupportsInMap(supports []int, disintegratedMap map[int]interface{}) bool {
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

	totalOtherDisintegrations := 0
	for currentBrickIndex, currentBrick := range pile.Bricks {
		currentBrickDisintegrationMap := make(map[int]interface{})
		currentBrickDisintegrationMap[currentBrickIndex] = PRESENCE_INDICATOR

		for consideringBrickIndex, consideringBrick := range pile.Bricks {
			if consideringBrick.Start.Z <= currentBrick.Start.Z {
				continue
			}

			consideringBrickSupports := pile.Supports[consideringBrickIndex]

			// If the considering brick  will be disintegrated by the current bricks removal, note this
			if checkSupportsInMap(consideringBrickSupports, currentBrickDisintegrationMap) {
				currentBrickDisintegrationMap[consideringBrickIndex] = PRESENCE_INDICATOR
			}
		}

		totalOtherDisintegrations += len(currentBrickDisintegrationMap) - 1
	}

	return totalOtherDisintegrations, nil
}

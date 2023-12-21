package part01

import (
	"bufio"
	"hmcalister/aoc21/lib"
)

func ProcessInput(fileScanner *bufio.Scanner) (int, error) {
	garden := lib.ParseFileToGardenData(*fileScanner)
	garden.DebugLog()
	numPlots := garden.NumReachableGardensInExactlyNumSteps(64)

	return numPlots, nil
}

package part02

import (
	"bufio"
	"hmcalister/aoc20/lib"
)

func ProcessInput(fileScanner *bufio.Scanner) (int, error) {
	moduleConfig := lib.ParseFileToModuleConfiguration(fileScanner)

	minButtonPressesFor_RX_LOW := moduleConfig.FindLowestButtonPushesToAchieve_RX_LOW()

	return minButtonPressesFor_RX_LOW, nil
}

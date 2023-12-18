package part02

import (
	"bufio"
	"hmcalister/aoc18/part02/lib"
)

func ProcessInput(fileScanner *bufio.Scanner) (int, error) {
	digLayout := lib.NewDigLayoutFromFileScanner(fileScanner)
	digLayout.VisualizeDigLayout()
	totalVolume := digLayout.CalculateTotalVolume()

	return totalVolume, nil
}

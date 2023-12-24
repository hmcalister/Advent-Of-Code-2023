package part01

import (
	"bufio"
	"hmcalister/aoc24/lib"
)

func ProcessInput(fileScanner *bufio.Scanner) (int, error) {
	storm := lib.ParseFileToStorm(fileScanner)
	storm.DetectCollisionsInXY(0, 1)

	return 0, nil
}

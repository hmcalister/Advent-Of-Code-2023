package part01

import (
	"bufio"
	"hmcalister/aoc24/lib"
)

func ProcessInput(fileScanner *bufio.Scanner) (int, error) {
	storm := lib.ParseFileToStorm(fileScanner)
	numCollisions := storm.PathIntersectionInXY(200000000000000, 400000000000000)

	return numCollisions, nil
}

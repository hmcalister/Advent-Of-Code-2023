package part02

import (
	"bufio"
	"hmcalister/aoc17/lib"
)

func ProcessInput(fileScanner *bufio.Scanner) (int, error) {
	layout := lib.NewLayoutFromFileScanner(fileScanner, 3, 10)
	goalDist := layout.PathFind()
	return goalDist, nil
}

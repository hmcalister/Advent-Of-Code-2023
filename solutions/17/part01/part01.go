package part01

import (
	"bufio"
	"hmcalister/aoc17/lib"
)

func ProcessInput(fileScanner *bufio.Scanner) (int, error) {
	layout := lib.NewLayoutFromFileScanner(fileScanner, 0, 3)
	goalDist := layout.PathFind()
	return goalDist, nil
}

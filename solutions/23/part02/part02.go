package part02

import (
	"bufio"
	"hmcalister/aoc23/lib"
)

func ProcessInput(fileScanner *bufio.Scanner) (int, error) {
	trail := lib.ParseFileToTrail(fileScanner)
	condensedTrail := lib.ConvertTrailDataToCondensedTrailData(trail)
	// file, _ := os.Create("./graphVis.gv")
	// draw.DOT(condensedTrail.TrailGraph, file)

	longestPath := condensedTrail.FindPathNonSlippery()

	return longestPath, nil
}

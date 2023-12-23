package part01

import (
	"bufio"
	"hmcalister/aoc23/lib"
)

func ProcessInput(fileScanner *bufio.Scanner) (int, error) {
	trail := lib.ParseFileToTrail(fileScanner)
	path, err := trail.FindPathSlippery()
	if err != nil {
		return -1, err
	}
	trail.VisualizePath(path)

	return path.PathLength(), nil
}

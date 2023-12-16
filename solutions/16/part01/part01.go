package part01

import (
	"bufio"
	"hmcalister/aoc16/lib"
	"slices"

	"github.com/rs/zerolog/log"
)

func ProcessInput(fileScanner *bufio.Scanner) (int, error) {
	layout := lib.NewLayoutData(fileScanner)
	layout.ShowLayout()
	layout.ProcessLayout()
	layout.ShowEnergizedCells()

	slices.Sort(layout.EnergizedLinearCoordinates)
	layout.EnergizedLinearCoordinates = slices.Compact(layout.EnergizedLinearCoordinates)
	log.Debug().Interface("EnergizedLinearCoords", layout.EnergizedLinearCoordinates).Send()

	return len(layout.EnergizedLinearCoordinates), nil
}

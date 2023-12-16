package part01

import (
	"bufio"
	"hmcalister/aoc16/lib"

	"github.com/rs/zerolog/log"
)

func ProcessInput(fileScanner *bufio.Scanner) (int, error) {
	layout := lib.NewLayoutData(fileScanner, &lib.LightRay{
		Direction: lib.DIRECTION_EAST,
		XCoord:    0,
		YCoord:    0,
	})
	layout.ShowLayout()
	layout.ProcessLayout()
	log.Debug().Interface("EnergizedLinearCoords", layout.EnergizedLinearCoordinates).Send()
	layout.ShowEnergizedCells()

	return len(layout.EnergizedLinearCoordinates), nil
}

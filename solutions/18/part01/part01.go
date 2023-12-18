package part01

import (
	"bufio"
	"hmcalister/aoc18/part01/lib"

	"github.com/rs/zerolog/log"
)

func ProcessInput(fileScanner *bufio.Scanner) (int, error) {
	digLayout := lib.NewDigLayoutFromFileScanner(fileScanner)

	log.Info().Msg("Trench Before Excavation")
	digLayout.VisualizeDigLayout()

	digLayout.ExcavateInterior()

	log.Info().Msg("Trench AfterExcavation")
	digLayout.VisualizeDigLayout()

	totalVolume := digLayout.CalculateTotalVolume()

	return totalVolume, nil
}

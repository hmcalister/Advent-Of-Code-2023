package part01

import (
	"bufio"
	"hmcalister/aoc17/lib"

	"github.com/rs/zerolog/log"
)

func ProcessInput(fileScanner *bufio.Scanner) (int, error) {
	costMap := lib.GetCostMapFromFileScanner(fileScanner)
	layout := lib.NewLayout(costMap)
	goalState := layout.AStarPathFind()
	if goalState == nil {
		log.Fatal().Msg("no path found")
	}
	return goalState.GScore, nil
}

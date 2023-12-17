package part01

import (
	"bufio"
	"hmcalister/aoc17/lib"

	"github.com/rs/zerolog/log"
)

func ProcessInput(fileScanner *bufio.Scanner) (int, error) {
	layout := lib.NewLayout(fileScanner)
	goalState := layout.DijkstraPathFind()
	if goalState == nil {
		log.Fatal().Msg("no path found")
	}
	return goalState.TraversedDistance, nil
}

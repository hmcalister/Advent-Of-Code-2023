package part01

import (
	"bufio"
	"hmcalister/aoc12/lib"

	"github.com/rs/zerolog/log"
)

func ProcessInput(fileScanner *bufio.Scanner) (int, error) {
	result := 0
	for fileScanner.Scan() {
		line := fileScanner.Text()
		row := lib.ParseLineToSpringRowData(line)
		rowArrangements := lib.CalculatePossibleArrangements(row.RowLine, row.ContiguousDamagedGroupData)
		result += rowArrangements

		log.Debug().
			Str("Line", line).
			Interface("SpringRow", row).
			Int("PossibleArrangements", rowArrangements).
			Int("CumulativeResult", result).
			Send()
	}
	return result, nil
}

package part01

import (
	"bufio"
	"hmcalister/aoc12/lib"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

func parseLineToSpringRowData(line string) lib.SpringRowData {
	fields := strings.Fields(line)
	log.Trace().
		Str("ParsedLine", line).
		Interface("Fields", fields).
		Msg("Parsing Line To Data")

	contiguousDamagedGroupDataStrs := strings.Split(fields[1], ",")
	contiguousDamagedGroupData := make([]int, len(contiguousDamagedGroupDataStrs))
	for i, str := range contiguousDamagedGroupDataStrs {
		parsedInt, err := strconv.Atoi(str)
		if err != nil {
			log.Fatal().Msgf("failed to parsed contiguousDamagedGroup string %v to integer on line %v", str, line)
		}
		log.Trace().
			Str("ContiguousGroupString", str).
			Int("ParsedInt", parsedInt).
			Send()
		contiguousDamagedGroupData[i] = parsedInt
	}

	rowLine := strings.Trim(fields[0], string(lib.OPERATIONAL_SPRING_RUNE))
	// rowLine = rowLine + string(OPERATIONAL_SPRING_RUNE)
	return lib.SpringRowData{
		RowLine:                    rowLine,
		ContiguousDamagedGroupData: contiguousDamagedGroupData,
	}
}

func ProcessInput(fileScanner *bufio.Scanner) (int, error) {
	result := 0
	possibleArrangementsCalculator := lib.NewPossibleArrangementsCalculator()
	for fileScanner.Scan() {
		line := fileScanner.Text()
		row := parseLineToSpringRowData(line)
		rowArrangements := possibleArrangementsCalculator.CalculatePossibleArrangements(row.RowLine, row.ContiguousDamagedGroupData)
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

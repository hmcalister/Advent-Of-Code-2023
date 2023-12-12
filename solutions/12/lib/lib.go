package lib

import (
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

const (
	DAMAGED_SPRING_RUNE     rune = '#'
	OPERATIONAL_SPRING_RUNE rune = '.'
	UNKNOWN_SPRING_RUNE     rune = '?'
)

type springRowData struct {
	RowLine                    string
	ContiguousDamagedGroupData []int
}

func parseLineToSpringRowData(line string) springRowData {
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

	rowLine := strings.Trim(fields[0], string(OPERATIONAL_SPRING_RUNE))
	// rowLine = rowLine + string(OPERATIONAL_SPRING_RUNE)
	return springRowData{
		RowLine:                    rowLine,
		ContiguousDamagedGroupData: contiguousDamagedGroupData,
	}
}

func calculatePossibleArrangements(line string, remainingGroups []int) int {
	log.Trace().
		Str("Line", line).
		Interface("RemainingGroups", remainingGroups).
		Msg("CalculatePossibleArrangements Start")

	// Strip away all the starting operational springs
	line = strings.TrimLeft(line, string(OPERATIONAL_SPRING_RUNE))

	if len(line) == 0 {
		log.Trace().Msg("no line remaining")
		if len(remainingGroups) == 0 {
			return 1
		} else {
			return 0
		}
	}

	// If our line now starts with an unknown spring, simply consider both cases of the unknown
	if line[0] == byte(UNKNOWN_SPRING_RUNE) {
		log.Trace().Msg("considering both variations of unknown start spring")
		unknownAsDamagedSpringLine := strings.Replace(line, string(UNKNOWN_SPRING_RUNE), string(DAMAGED_SPRING_RUNE), 1)
		unknownAsOperationalSpringLine := strings.Replace(line, string(UNKNOWN_SPRING_RUNE), string(OPERATIONAL_SPRING_RUNE), 1)
		return calculatePossibleArrangements(unknownAsDamagedSpringLine, remainingGroups) + calculatePossibleArrangements(unknownAsOperationalSpringLine, remainingGroups)
	}

	// Otherwise, our line must start with a damaged spring! ----------------------------------------------------------

	// If we have no groups left it is always possible in one way, no matter the line remaining
	if len(remainingGroups) == 0 {
		log.Trace().Msg("no groups remaining")
		return 0
	}

	// If the length of our string is too short to accommodate even the first group, it is impossible
	if len(line) < remainingGroups[0] {
		log.Trace().Msg("line too short for first group")
		return 0
	}

	// We now assume the first remaining group occurs at the very start of the line
	// If there are any undamaged springs in the first group, it is invalid
	if strings.ContainsRune(line[:remainingGroups[0]], OPERATIONAL_SPRING_RUNE) {
		log.Trace().Msg("first group contains operational springs")
		return 0
	}

	// If the number of groups remaining is greater than one, we have to consider subproblems
	if len(remainingGroups) > 1 {
		// If the line cannot fit the remaining group OR the first group does not terminate as expected, this is invalid
		if len(line) < remainingGroups[0]+1 || line[remainingGroups[0]] == byte(DAMAGED_SPRING_RUNE) {
			log.Trace().Msg("first group not terminated with unknown or operational spring")
			return 0
		}

		// Otherwise, we can try the subproblem by removing the first group and the associated start of line!
		return calculatePossibleArrangements(line[remainingGroups[0]+1:], remainingGroups[1:])
	} else {
		// There is exactly one group remaining!
		return calculatePossibleArrangements(line[remainingGroups[0]:], remainingGroups[1:])
	}
}

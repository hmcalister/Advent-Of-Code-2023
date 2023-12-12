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

type SpringRowData struct {
	RowLine                    string
	ContiguousDamagedGroupData []int
}

// A calculator to keep track of all possible arrangements - with caching!
type PossibleArrangementsCalculator struct {
	possibleARrangementsCache map[string]int
}

func NewPossibleArrangementsCalculator() *PossibleArrangementsCalculator {
	return &PossibleArrangementsCalculator{
		possibleARrangementsCache: make(map[string]int),
	}
}

func (calc *PossibleArrangementsCalculator) CalculatePossibleArrangements(line string, remainingGroups []int) int {
	inputHash := calc.convertLineAndGroupsToString(line, remainingGroups)
	if result, ok := calc.possibleARrangementsCache[inputHash]; ok {
		return result
	} else {
		result := calc.calculatePossibleArrangementsRecursive(line, remainingGroups)
		calc.possibleARrangementsCache[inputHash] = result
		return result
	}
}

func (calc *PossibleArrangementsCalculator) convertLineAndGroupsToString(line string, remainingGroups []int) string {
	var strBuilder strings.Builder
	for i, v := range remainingGroups {
		if i > 0 {
			strBuilder.WriteString(",")
		}
		strBuilder.WriteString(strconv.Itoa(v))
	}
	remainingGroupsStr := strBuilder.String()
	return line + " " + remainingGroupsStr
}

func (calc *PossibleArrangementsCalculator) calculatePossibleArrangementsRecursive(line string, remainingGroups []int) int {
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
		return calc.CalculatePossibleArrangements(unknownAsDamagedSpringLine, remainingGroups) + calc.CalculatePossibleArrangements(unknownAsOperationalSpringLine, remainingGroups)
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
		return calc.CalculatePossibleArrangements(line[remainingGroups[0]+1:], remainingGroups[1:])
	} else {
		// There is exactly one group remaining!
		return calc.CalculatePossibleArrangements(line[remainingGroups[0]:], remainingGroups[1:])
	}
}

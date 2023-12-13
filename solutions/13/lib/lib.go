package lib

import (
	"bufio"
	"errors"

	"github.com/rs/zerolog/log"
)

const (
	ASH_RUNE  = '.'
	ROCK_RUNE = '#'
)

func reverseString(s string) string {
	reversedRunes := make([]rune, len(s))
	for i, r := range s {
		reversedRunes[len(s)-i-1] = r
	}
	return string(reversedRunes)
}

func parseSectionToPattern(fileScanner *bufio.Scanner, patternID int) (PatternData, error) {
	var line string
	currentPatternRows := make([]string, 0)
	// Take from the file scanner until the next blank line
	for {
		if !fileScanner.Scan() {
			log.Trace().Msg("End of File found")
			return PatternData{}, errors.New("end of file")
		}
		line = fileScanner.Text()
		// Check if we have reached the end of a pattern
		if len(line) == 0 {
			log.Trace().Msg("Blank line found")
			break
		}

		log.Trace().
			Str("NextLine", line).
			Send()
		currentPatternRows = append(currentPatternRows, line)
	}

	currentPattern := newPatternData(patternID, currentPatternRows)

	log.Trace().Msgf("Finished Parsing Pattern %v", currentPattern.PatternID)
	return currentPattern, nil
}

func ParseFileToPatterns(fileScanner *bufio.Scanner) []PatternData {

	filePatterns := make([]PatternData, 0)

	for {
		log.Trace().
			Int("PatternID", len(filePatterns)).
			Msg("Start Parsing Pattern")

		newPattern, err := parseSectionToPattern(fileScanner, len(filePatterns))
		if err != nil {
			break
		}
		filePatterns = append(filePatterns, newPattern)
	}

	return filePatterns
}

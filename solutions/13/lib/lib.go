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

func findReflectionIndices(s string) []int {
	reversedString := reverseString(s)
	reflectionIndices := make([]int, 0)

	log.Trace().
		Str("TestString", s).
		Str("ReversedString", reversedString).
		Int("I_Limit", len(s)-1).
		Send()

	// We can check for a reflection by taking progressively smaller
	// substrings of the reversed string, and asking if testString ends with it
	for i := 1; i <= len(s)-1; i += 1 {
		strLen := min(i, len(s)-i)
		forwardStart := i - strLen
		forwardEnd := i
		reverseStart := len(s) - i - strLen
		reverseEnd := len(s) - i
		forwardPartial := s[forwardStart:forwardEnd]
		reversePartial := reversedString[reverseStart:reverseEnd]
		log.Trace().
			Int("TestingReflectionIndex", i).
			Str("ForwardPartial", forwardPartial).
			Str("ReversePartial", reversePartial).
			Send()
		if reversePartial == forwardPartial {
			log.Trace().Msg("Match Found")
			reflectionIndices = append(reflectionIndices, i)
		}
	}

	return reflectionIndices
}

type PatternData struct {
	PatternID int
	Rows      []string
	Columns   []string
}

func newPatternData(ID int, rows []string) PatternData {
	columns := make([]string, len(rows[0]))
	for columnIndex := range columns {
		columnRunes := make([]rune, len(rows))
		for rowIndex, row := range rows {
			columnRunes[rowIndex] = rune(row[columnIndex])
		}
		columns[columnIndex] = string(columnRunes)
	}

	return PatternData{
		PatternID: ID,
		Rows:      rows,
		Columns:   columns,
	}
}

func (pattern PatternData) FindReflections() (int, int) {
	rowReflectionIndexCounts := make(map[int]int)
	// Find the possible reflections across all rows
	for rowIndex, row := range pattern.Rows {
		log.Trace().
			Int("PatternID", pattern.PatternID).
			Int("RowIndex", rowIndex).
			Str("Row", row).
			Msg("Start Reflection Finding")
		indices := findReflectionIndices(row)
		log.Debug().
			Int("PatternID", pattern.PatternID).
			Int("RowIndex", rowIndex).
			Str("Row", row).
			Interface("ReflectionIndices", indices).
			Msg("Row Reflection Finding Results")
		for _, index := range indices {
			rowReflectionIndexCounts[index] += 1
		}
	}

	// Only take those indices that are present across all rows
	rowReflectionIndices := make([]int, 0)
	for k, v := range rowReflectionIndexCounts {
		if v == len(pattern.Rows) {
			rowReflectionIndices = append(rowReflectionIndices, k)
		}
	}

	columnReflectionIndexCounts := make(map[int]int)
	// Find the possible reflections across all columns
	for columnIndex, column := range pattern.Columns {
		log.Trace().
			Int("PatternID", pattern.PatternID).
			Int("ColumnIndex", columnIndex).
			Str("Column", column).
			Msg("Start Reflection Finding")
		indices := findReflectionIndices(column)
		log.Debug().
			Int("PatternID", pattern.PatternID).
			Int("ColumnIndex", columnIndex).
			Str("Column", column).
			Interface("ReflectionIndices", indices).
			Msg("Column Reflection Finding Results")
		for _, index := range indices {
			columnReflectionIndexCounts[index] += 1
		}
	}

	// Only take those indices that are present across all rows
	columnReflectionIndices := make([]int, 0)
	for k, v := range columnReflectionIndexCounts {
		if v == len(pattern.Columns) {
			columnReflectionIndices = append(columnReflectionIndices, k)
		}
	}

	log.Debug().
		Interface("RowReflectionIndices", rowReflectionIndices).
		Interface("ColumnReflectionIndices", columnReflectionIndices).
		Send()

	if len(rowReflectionIndices) == 1 {
		return rowReflectionIndices[0], 0

	} else if len(columnReflectionIndices) == 1 {
		return 0, columnReflectionIndices[0]
	} else {
		log.Fatal().Msgf("Failed to find unique reflection indices for pattern %v", pattern.PatternID)
		return 0, 0
	}
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

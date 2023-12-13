package lib

import "github.com/rs/zerolog/log"

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

func (pattern PatternData) FindReflectionsIndices() (int, int) {
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

func (pattern PatternData) FindSmudgedReflectionsIndices() (int, int) {
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
		if v == len(pattern.Rows)-1 {
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
		if v == len(pattern.Columns)-1 {
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

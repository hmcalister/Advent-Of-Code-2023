package part01

import (
	"bufio"
	"hmcalister/aoc13/lib"

	"github.com/rs/zerolog/log"
)

func ProcessInput(fileScanner *bufio.Scanner) (int, error) {
	filePatterns := lib.ParseFileToPatterns(fileScanner)
	log.Debug().Int("NumberOfPatterns", len(filePatterns)).Send()

	result := 0
	for _, pattern := range filePatterns {
		log.Debug().Interface("Pattern", pattern).Send()
		rowReflectionIndex, columnReflectionIndex := pattern.FindReflections()

		result += rowReflectionIndex + 100*columnReflectionIndex
	}

	return result, nil
}

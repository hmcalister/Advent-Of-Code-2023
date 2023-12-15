package part01

import (
	"bufio"
	"strings"

	"github.com/rs/zerolog/log"
)

func HASHAlgorithm(s string) int {
	h := 0

	for _, currentRune := range s {
		h += int(currentRune)
		h *= 17
		h = h % 256
	}

	return h
}

func ProcessInput(fileScanner *bufio.Scanner) (int, error) {
	fullText := ""
	for fileScanner.Scan() {
		line := fileScanner.Text()
		fullText += line
	}

	fullTextFields := strings.Split(fullText, ",")

	result := 0
	for fieldIndex, field := range fullTextFields {
		fieldHash := HASHAlgorithm(field)

		log.Debug().
			Int("FieldIndex", fieldIndex).
			Str("Field", field).
			Int("Hash", fieldHash).
			Send()

		result += fieldHash
	}

	return result, nil
}

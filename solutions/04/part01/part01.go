package part01

import (
	"bufio"
	"slices"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

type scratchCardData struct {
	CardID         int
	WinningNumbers []int
	FoundNumbers   []int
	Score          int
}

// Convert a list of integers as a string (space separated) to an integer array
func stringToIntArray(s string) []int {
	var err error

	intStrs := strings.Fields(s)
	log.Trace().Str("StrToParse", s).Int("NumFieldsFound", len(intStrs)).Send()

	parsedInts := make([]int, len(intStrs))
	for i, intStr := range intStrs {
		parsedInts[i], err = strconv.Atoi(intStr)
		if err != nil {
			log.Fatal().Msgf("error when parsing integer %v", err)
		}
	}

	return parsedInts
}

func createScratchCard(line string) *scratchCardData {
	var err error
	cardData := &scratchCardData{}

	// The line starts with "Card X:", so we find the card number and strip this away
	colonIndex := strings.IndexRune(line, ':')
	cardIDStr := strings.TrimSpace(line[5:colonIndex])
	cardData.CardID, err = strconv.Atoi(cardIDStr)
	log.Trace().
		Str("CardIDStr", cardIDStr).
		Send()

	if err != nil {
		log.Fatal().Msgf("error when parsing cardID: %v", err)
	}
	line = line[colonIndex+1:]

	log.Trace().
		Int("CardID", cardData.CardID).
		Str("RemainingConfigStr", line).
		Send()

	// Next step, separate the two halves, the winning numbers and found numbers
	barIndex := strings.IndexRune(line, '|')
	winningNumbersStr := strings.TrimSpace(line[:barIndex])
	foundNumbersStr := strings.TrimSpace(line[barIndex+1:])

	log.Trace().
		Int("CardID", cardData.CardID).
		Int("BarIndex", barIndex).
		Str("WinningNumberStr", winningNumbersStr).
		Str("FoundNumberStr", foundNumbersStr).
		Send()

	cardData.WinningNumbers = stringToIntArray(winningNumbersStr)
	cardData.FoundNumbers = stringToIntArray(foundNumbersStr)

	for _, n := range cardData.FoundNumbers {
		if slices.Contains(cardData.WinningNumbers, n) {
			if cardData.Score == 0 {
				cardData.Score = 1
			} else {
				cardData.Score *= 2
			}
			log.Trace().
				Int("CardID", cardData.CardID).
				Int("FoundWinningNumber", n).
				Int("NewScore", cardData.Score).
				Send()
		}
	}

	log.Debug().
		Int("CardID", cardData.CardID).
		Int("FinalScore", cardData.Score).
		Send()

	return cardData
}

// Given a scanner over the input file, calculate the total number of points
// that the scratchcards have earned, and return it
func ProcessInput(fileScanner *bufio.Scanner) (int, error) {
	result := 0

	for fileScanner.Scan() {
		line := fileScanner.Text()
		cardData := createScratchCard(line)
		result += cardData.Score
	}
	return result, nil
}

package part01

import (
	"bufio"
	"hmcalister/aoc07/part01/lib"
	"sort"

	"github.com/rs/zerolog/log"
)

func ProcessInput(fileScanner *bufio.Scanner) (int, error) {
	allHands := make([]lib.HandData, 0)
	for fileScanner.Scan() {
		line := fileScanner.Text()
		allHands = append(allHands, lib.ParseLineToHandData(line))
	}

	log.Debug().
		Int("NumHandsFound", len(allHands)).
		Send()

	sort.Slice(allHands, func(i, j int) bool {
		return allHands[i].HandStrength < allHands[j].HandStrength
	})

	result := 0
	for handPosition, hand := range allHands {
		handWinnings := hand.BidAmount * (handPosition + 1)
		log.Debug().
			Int("HandPosition", handPosition).
			Int("HandWinnings", handWinnings).
			Interface("Hand", hand).
			Send()
		result += handWinnings
	}

	return result, nil
}

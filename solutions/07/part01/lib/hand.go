package lib

import (
	"math"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

//go:generate stringer -type=HandTypeEnum
type HandTypeEnum int

const (
	CARD_STRENGTH = "23456789TJQKA"

	HIGH_CARD       HandTypeEnum = iota
	ONE_PAIR        HandTypeEnum = iota
	TWO_PAIR        HandTypeEnum = iota
	THREE_OF_A_KIND HandTypeEnum = iota
	FULL_HOUSE      HandTypeEnum = iota
	FOUR_OF_A_KIND  HandTypeEnum = iota
	FIVE_OF_A_KIND  HandTypeEnum = iota
)

type HandData struct {
	BidAmount     int
	CardStrengths []int
	HandType      HandTypeEnum
	HandStrength  float64
}

func calculateHandType(cardStrengths []int) HandTypeEnum {
	strengthCounts := make(map[int]int)
	for _, strength := range cardStrengths {
		strengthCounts[strength] += 1
	}

	log.Trace().Interface("CardCounts", strengthCounts).Send()

	threeOfAKindNums := make([]int, 0)
	twoOfAKindNums := make([]int, 0)
	// Check for hand types that can be determined in the first pass
	for num, count := range strengthCounts {
		if count == 5 {
			return FIVE_OF_A_KIND
		}

		if count == 4 {
			return FOUR_OF_A_KIND
		}

		if count == 3 {
			threeOfAKindNums = append(threeOfAKindNums, num)
		}

		if count == 2 {
			twoOfAKindNums = append(twoOfAKindNums, num)
		}
	}

	if len(threeOfAKindNums) == 1 && len(twoOfAKindNums) == 1 {
		return FULL_HOUSE
	} else if len(threeOfAKindNums) == 1 {
		return THREE_OF_A_KIND
	} else if len(twoOfAKindNums) == 2 {
		return TWO_PAIR
	} else if len(twoOfAKindNums) == 1 {
		return ONE_PAIR
	} else {
		return HIGH_CARD
	}
}

func calculateHandPartialStrength(cardStrengths []int) float64 {
	partialStrength := 0.0
	for position, strength := range cardStrengths {
		partialStrength += math.Pow(float64(len(CARD_STRENGTH)), -float64(position+1)) * float64(strength)
	}

	return partialStrength
}

func ParseLineToHandData(line string) HandData {
	fields := strings.Fields(line)
	log.Debug().
		Str("ParsingLine", line).
		Int("NumFields", len(fields)).
		Send()

	bidAmount, err := strconv.Atoi(fields[1])
	if err != nil {
		log.Panic().Msgf("failed to parse bid amount %v in line %v", fields[1], line)
	}

	cardStrengths := make([]int, len(fields[0]))
	for i, currentCard := range fields[0] {
		cardStrength := strings.IndexRune(CARD_STRENGTH, currentCard)
		if cardStrength == -1 {
			log.Panic().Msgf("encountered unknown card value %v in line %v", currentCard, line)
		}
		cardStrengths[i] = cardStrength
	}

	handType := calculateHandType(cardStrengths)
	handStrength := calculateHandPartialStrength(cardStrengths) + float64(handType)

	parsedHandData := HandData{
		bidAmount,
		cardStrengths,
		handType,
		handStrength,
	}

	log.Debug().
		Interface("ParsedHandData", parsedHandData).
		Str("HandType", parsedHandData.HandType.String()).
		Send()

	return parsedHandData
}

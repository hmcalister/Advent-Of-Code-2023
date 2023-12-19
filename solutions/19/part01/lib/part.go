package lib

import (
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

const (
	ExtremelyCoolString string = "x"
	MusicalString       string = "m"
	AerodynamicString   string = "a"
	ShinyString         string = "s"
)

type PartData struct {
	ExtremelyCoolRating int
	MusicalRating       int
	AerodynamicRating   int
	ShinyRating         int
}

func ParseLineToPartData(line string) PartData {
	line = line[1 : len(line)-1]
	fields := strings.FieldsFunc(line, func(r rune) bool {
		return r == ','
	})

	newPart := PartData{}

	for _, field := range fields {
		propertyValue, err := strconv.Atoi(field[2:])
		if err != nil {
			log.Fatal().Msgf("failed to parse part property value %v in line %v", field[2:], line)
		}

		switch field[:1] {
		case ExtremelyCoolString:
			newPart.ExtremelyCoolRating = propertyValue
		case MusicalString:
			newPart.MusicalRating = propertyValue
		case AerodynamicString:
			newPart.AerodynamicRating = propertyValue
		case ShinyString:
			newPart.ShinyRating = propertyValue
		}
	}

	return newPart
}

func (part *PartData) SumRatings() int {
	return part.ExtremelyCoolRating + part.MusicalRating + part.AerodynamicRating + part.ShinyRating
}

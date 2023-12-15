package part02

import (
	"bufio"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

type LensData struct {
	Identifier  string
	FocalLength int
}

type BoxArrayData struct {
	Boxes [][]LensData
}

func NewBoxArray() *BoxArrayData {
	boxes := make([][]LensData, 256)
	for boxIndex := range boxes {
		boxes[boxIndex] = make([]LensData, 0)
	}
	return &BoxArrayData{
		Boxes: boxes,
	}
}

func (boxArray *BoxArrayData) removeLabel(label string) {
	labelHash := HASHAlgorithm(label)
	targetBox := boxArray.Boxes[labelHash]
	for lensIndex, lens := range targetBox {
		if lens.Identifier == label {
			for i := lensIndex; i < len(targetBox)-1; i += 1 {
				targetBox[i] = targetBox[i+1]
			}
			boxArray.Boxes[labelHash] = targetBox[:len(targetBox)-1]
			log.Trace().
				Int("BoxIndex", labelHash).
				Str("LensRemoved", lens.Identifier).
				Msg("Lens Removed")
			return
		}
	}

	log.Trace().Msg("No Lens Removed")
}

func (boxArray *BoxArrayData) addLens(label string, focalLen int) {
	newLens := LensData{
		Identifier:  label,
		FocalLength: focalLen,
	}

	labelHash := HASHAlgorithm(label)
	targetBox := boxArray.Boxes[labelHash]
	for lensIndex, lens := range targetBox {
		if lens.Identifier == newLens.Identifier {
			log.Trace().
				Int("BoxIndex", labelHash).
				Interface("ReplacedLens", lens).
				Interface("NewLens", newLens).
				Msg("Lens Replaced")
			targetBox[lensIndex] = newLens
			return
		}
	}

	log.Trace().
		Int("BoxIndex", labelHash).
		Interface("Lens", newLens).
		Msg("Lens Added")
	targetBox = append(targetBox, newLens)
	boxArray.Boxes[labelHash] = targetBox
}

func (boxArray *BoxArrayData) processField(field string) {
	log.Trace().Str("Field", field).Msg("Processing Field")

	var label string
	if strings.ContainsRune(field, '-') {
		label = field[:strings.IndexRune(field, '-')]
		boxArray.removeLabel(label)
	} else if strings.ContainsRune(field, '=') {
		label = field[:strings.IndexRune(field, '=')]
		focalLenStr := field[strings.IndexRune(field, '=')+1:]
		focalLen, err := strconv.Atoi(focalLenStr)
		if err != nil {
			log.Fatal().Msgf("could not parse given focal len %v from field %v to integer", focalLenStr, field)
		}
		boxArray.addLens(label, focalLen)
	} else {
		log.Fatal().Msgf("field %v did not contain expected '-' or '='", field)
	}
}

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
	fileScanner.Scan()
	fullText := fileScanner.Text()
	fullTextFields := strings.Split(fullText, ",")

	boxArray := NewBoxArray()

	for _, field := range fullTextFields {
		boxArray.processField(field)
	}

	totalPower := 0
	for boxIndex, box := range boxArray.Boxes {
		boxPower := 0
		for lensIndex, lens := range box {
			lensPower := (boxIndex + 1) * (1 + lensIndex) * lens.FocalLength
			log.Debug().
				Int("BoxIndex", boxIndex).
				Int("LensIndex", lensIndex).
				Str("LensLabel", lens.Identifier).
				Int("LensPower", lensPower).
				Send()
			// log.Trace().
			// 	Int("BoxContrib", (boxIndex+1)).
			// 	Int("SlotContrib", (1+lensIndex)).
			// 	Int("FocalContrib", lens.FocalLength).
			// 	Send()
			boxPower += lensPower
		}
		if boxPower != 0 {
			log.Debug().
				Int("BoxIndex", boxIndex).
				Int("BoxPower", boxPower).
				Send()
		}
		totalPower += boxPower
	}

	return totalPower, nil
}

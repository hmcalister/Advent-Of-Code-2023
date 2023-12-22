package part01

import (
	"bufio"
	"hmcalister/aox22/lib"

	"github.com/rs/zerolog/log"
)

func ProcessInput(fileScanner *bufio.Scanner) (int, error) {
	pile := lib.ParseFileToBrickPile(fileScanner)
	pile.SimulateBrickFall()

	disintegrateBrickIndicesMap := make(map[int]bool)
	for i := 0; i < len(pile.Bricks); i += 1 {
		disintegrateBrickIndicesMap[i] = true
	}
	for i := 0; i < len(pile.Bricks); i += 1 {
		log.Info().
			Int("BrickIndex", i).
			Str("Brick", pile.Bricks[i].String()).
			Interface("SupportingBricks", pile.Supports[i]).
			Send()
		if len(pile.Supports[i]) == 1 {
			disintegrateBrickIndicesMap[pile.Supports[i][0]] = false
		}
	}

	disintegrateBrickIndicesArr := make([]int, 0)
	for i := 0; i < len(pile.Bricks); i += 1 {
		if disintegrateBrickIndicesMap[i] {
			disintegrateBrickIndicesArr = append(disintegrateBrickIndicesArr, i)
		}
	}
	log.Info().Interface("DisintegrateBrickIndices", disintegrateBrickIndicesArr).Send()

	return len(disintegrateBrickIndicesArr), nil
}

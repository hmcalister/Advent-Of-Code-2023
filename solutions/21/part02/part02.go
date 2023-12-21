package part02

import (
	"bufio"
	"hmcalister/aoc21/lib"

	"github.com/openacid/slimarray/polyfit"
	"github.com/rs/zerolog/log"
)

func ProcessInput(fileScanner *bufio.Scanner) (int, error) {
	garden := lib.ParseFileToGardenData(*fileScanner)
	garden.DebugLog()

	// Since grid is square and start row/col has no rocks, {f(n), f(n+width), f(n+2*width),...} is quadratic
	//
	// So we only need to find n (the total number of steps modulo the width of the grid) and those three values.
	// Then we can construct a polynomial that fits those three values and calculate f(n+X*width) for an appropriate X

	numSteps := 26501365
	mapSize := garden.MapWidth
	n := numSteps % mapSize

	log.Debug().
		Int("NumSteps", numSteps).
		Int("MapSize", mapSize).
		Int("n", n).
		Interface("ProbeValues", []int{n, n + mapSize, n + 2*mapSize}).
		Send()

	// log.Debug().Int("Test", garden.NumReachableGardensInExactlyNumSteps(100)).Send()
	results := garden.ProbeWithValues([]int{n, n + mapSize, n + 2*mapSize})

	xVals := []float64{float64(n), float64(n + mapSize), float64(n + 2*mapSize)}
	yVals := []float64{float64(results[0]), float64(results[1]), float64(results[2])}
	polynomial := polyfit.NewFit(xVals, yVals, 2).Solve()

	X := float64(n + (numSteps/mapSize)*mapSize)
	result := polynomial[0] + X*polynomial[1] + X*X*polynomial[2]

	log.Debug().Float64("Result", result).Send()

	return 0, nil
}

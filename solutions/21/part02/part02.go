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

	// garden.NumReachableGardensInExactlyNumSteps(70)
	results := garden.FindNewPlotsAtValues([]int{n, n + mapSize, n + 2*mapSize})

	xVals := []float64{0.0, 1.0, 2.0}
	yVals := []float64{float64(results[0]), float64(results[1]), float64(results[2])}
	polynomial := polyfit.NewFit(xVals, yVals, 2).Solve()

	// polynomial := []float64{4059, 15823, 15397}
	X := float64(numSteps / mapSize)
	log.Debug().
		Interface("yVals", yVals).
		Interface("Coefficients", polynomial).
		Float64("X", X).
		Send()
	result := polynomial[0] + X*polynomial[1] + X*X*polynomial[2]

	log.Debug().Float64("Result", result).Send()

	// 630129824772393

	return 0, nil
}

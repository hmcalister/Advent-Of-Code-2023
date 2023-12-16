package part02

import (
	"bufio"
	"hmcalister/aoc16/lib"
	"math"

	"github.com/schollz/progressbar/v3"
)

func ProcessInput(fileScanner *bufio.Scanner) (int, error) {
	var pbar *progressbar.ProgressBar

	layoutRunes := lib.CreateLayoutData(fileScanner)
	yLim := len(layoutRunes)
	xLim := len(layoutRunes[0])

	highestEnergizedVal := math.MinInt

	// Check north edge
	pbar = progressbar.Default(int64(xLim), "North Edge")
	for x := 0; x < xLim; x += 1 {
		layout := lib.NewLayoutData(layoutRunes, &lib.LightRay{
			Direction: lib.DIRECTION_SOUTH,
			XCoord:    x,
			YCoord:    0,
		})
		layout.ProcessLayout()
		highestEnergizedVal = max(highestEnergizedVal, len(layout.EnergizedLinearCoordinates))
		pbar.Add(1)
	}

	// Check east edge
	pbar = progressbar.Default(int64(yLim), "East Edge")
	for y := 0; y < yLim; y += 1 {
		layout := lib.NewLayoutData(layoutRunes, &lib.LightRay{
			Direction: lib.DIRECTION_WEST,
			XCoord:    xLim - 1,
			YCoord:    y,
		})
		layout.ProcessLayout()
		highestEnergizedVal = max(highestEnergizedVal, len(layout.EnergizedLinearCoordinates))
		pbar.Add(1)
	}

	// Check south edge
	pbar = progressbar.Default(int64(xLim), "South Edge")
	for x := 0; x < xLim; x += 1 {
		layout := lib.NewLayoutData(layoutRunes, &lib.LightRay{
			Direction: lib.DIRECTION_NORTH,
			XCoord:    x,
			YCoord:    yLim - 1,
		})
		layout.ProcessLayout()
		highestEnergizedVal = max(highestEnergizedVal, len(layout.EnergizedLinearCoordinates))
		pbar.Add(1)
	}

	// Check west edge
	pbar = progressbar.Default(int64(yLim), "West Edge")
	for y := 0; y < yLim; y += 1 {
		layout := lib.NewLayoutData(layoutRunes, &lib.LightRay{
			Direction: lib.DIRECTION_EAST,
			XCoord:    0,
			YCoord:    y,
		})
		layout.ProcessLayout()
		highestEnergizedVal = max(highestEnergizedVal, len(layout.EnergizedLinearCoordinates))
		pbar.Add(1)
	}

	return highestEnergizedVal, nil
}

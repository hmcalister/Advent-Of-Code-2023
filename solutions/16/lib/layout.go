package lib

import (
	"bufio"

	"github.com/rs/zerolog/log"
)

type LayoutData struct {
	Layout                     [][]LayoutRuneEnum
	LineLength                 int
	EnergizedLinearCoordinates []int
	ProcessedLightRayMap       map[string]interface{}
	UnprocessedLightRays       []*LightRay
	directionMap               map[DirectionEnum]map[LayoutRuneEnum][]DirectionEnum
}

func NewLayoutData(fileScanner *bufio.Scanner) *LayoutData {
	var line string
	var lineLength int

	layout := make([][]LayoutRuneEnum, 0)

	for fileScanner.Scan() {
		line = fileScanner.Text()
		lineLength = len(line)
		layout = append(layout, make([]LayoutRuneEnum, lineLength))

		for runeIndex, r := range line {
			layoutRune := LayoutRuneEnum(r)
			log.Trace().
				Int("LineIndex", len(layout)-1).
				Int("RuneIndex", runeIndex).
				Str("RuneFound", layoutRune.String()).
				Send()
			layout[len(layout)-1][runeIndex] = layoutRune
		}
	}

	return &LayoutData{
		Layout:                     layout,
		LineLength:                 lineLength,
		EnergizedLinearCoordinates: make([]int, 0),
		ProcessedLightRayMap:       make(map[string]interface{}),
		UnprocessedLightRays: []*LightRay{{
			Direction: DIRECTION_EAST,
			XCoord:    0,
			YCoord:    0,
		}},
		directionMap: CreateDirectionMap(),
	}
}

func (layout *LayoutData) CartesianToLinearCoordinate(x, y int) int {
	return y*layout.LineLength + x
}

func (layout *LayoutData) LinearToCartesianCoordinate(l int) (int, int) {
	return l % layout.LineLength, l / layout.LineLength
}

func (layout *LayoutData) processLightRay(currentRay *LightRay) {
	currentRune := layout.Layout[currentRay.YCoord][currentRay.XCoord]
	nextDirections := layout.directionMap[currentRay.Direction][currentRune]
	for _, dir := range nextDirections {
		nextLightRay := &LightRay{
			Direction: dir,
			XCoord:    currentRay.XCoord,
			YCoord:    currentRay.YCoord,
		}
		nextLightRay.MarchRay()

		// Ensure light ray is actually in the bounds of the layout
		if nextLightRay.XCoord < 0 || nextLightRay.XCoord >= layout.LineLength || nextLightRay.YCoord < 0 || nextLightRay.YCoord >= len(layout.Layout) {
			continue
		}

		layout.UnprocessedLightRays = append(layout.UnprocessedLightRays, nextLightRay)
	}
}

func (layout *LayoutData) ProcessLayout() {
	var currentRay *LightRay
	for len(layout.UnprocessedLightRays) > 0 {
		currentRay, layout.UnprocessedLightRays = layout.UnprocessedLightRays[0], layout.UnprocessedLightRays[1:]
		currentRayStr := currentRay.String()

		// Check if this exact light ray has already been processed
		if _, ok := layout.ProcessedLightRayMap[currentRayStr]; ok {
			continue
		}
		layout.ProcessedLightRayMap[currentRayStr] = currentRay
		layout.EnergizedLinearCoordinates = append(layout.EnergizedLinearCoordinates, layout.CartesianToLinearCoordinate(currentRay.XCoord, currentRay.YCoord))
		layout.processLightRay(currentRay)
	}
}

func (layout *LayoutData) ShowLayout() {
	for _, line := range layout.Layout {
		log.Info().
			Str("Layout", string(line)).
			Send()
	}
}

func (layout *LayoutData) ShowEnergizedCells() {
	cells := make([][]LayoutRuneEnum, len(layout.Layout))
	for lineIndex := range cells {
		newLine := make([]LayoutRuneEnum, layout.LineLength)
		for cellIndex := 0; cellIndex < layout.LineLength; cellIndex += 1 {
			newLine[cellIndex] = priv_NON_ENERGIZED_RUNE
		}
		cells[lineIndex] = newLine
	}

	for _, energizedLinearCoord := range layout.EnergizedLinearCoordinates {
		x, y := layout.LinearToCartesianCoordinate(energizedLinearCoord)
		cells[y][x] = priv_ENERGIZED_RUNE
	}

	for _, line := range cells {
		log.Info().
			Str("EnergizedCells", string(line)).
			Send()
	}
}

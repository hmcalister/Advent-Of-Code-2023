package lib

import (
	"bufio"
	"errors"
	"slices"

	"github.com/rs/zerolog/log"
)

type BrickPileData struct {
	// The actual bricks in the pile. Will start as the snapshot of bricks,
	// then (once SimulateBrickFall is called) will update to a fallen brick stack.
	Bricks []BrickData

	// For each brick in the Bricks array, gives a list of other BrickIndices
	// that support this brick.
	Supports [][]int

	// A map from coordinate to brick index
	CoordinateMap map[Coordinate]int
}

func ParseFileToBrickPile(fileScanner *bufio.Scanner) BrickPileData {
	pile := BrickPileData{
		Bricks:        make([]BrickData, 0),
		Supports:      make([][]int, 0),
		CoordinateMap: make(map[Coordinate]int),
	}

	var line string
	for fileScanner.Scan() {
		line = fileScanner.Text()
		brick := ParseLineToBrick(line)
		brick.ID = len(pile.Bricks)

		pile.Bricks = append(pile.Bricks, brick)
		pile.Supports = append(pile.Supports, make([]int, 0))

		for _, coord := range brick.enumerateAllCoordinates() {
			if _, ok := pile.CoordinateMap[coord]; ok {
				log.Fatal().Msgf("brick collision at coordinate %v", coord)
			}
			pile.CoordinateMap[coord] = brick.ID
		}
	}

	return pile
}

func (pile BrickPileData) brickIsSupported(brickIndex int) bool {
	return pile.Bricks[brickIndex].Start.Z == 0 || len(pile.Supports[brickIndex]) != 0
}

func (pile BrickPileData) findUnimpededUnsupportedBrickIndex() (int, error) {
	targetBrickIndex := -1

	// Select the first brick that is unsupported
	for brickIndex := range pile.Supports {
		if !pile.brickIsSupported(brickIndex) {
			targetBrickIndex = brickIndex
			break
		}
	}

	if targetBrickIndex == -1 {
		log.Trace().Msg("no unsupported bricks found")
		return -1, errors.New("no unsupported brick found")
	}

	log.Trace().Int("InitialUnsupportedBrickIndex", targetBrickIndex).Msg("FindUnimpededUnsupportedBrick")

	// For each X,Y coordinate of this brick, starting at the ground and moving up to the brick:
	// 1) a supported brick is found, (continue to the next coordinate)
	// 2) an unsupported brick is found (change targetBrickIndex and try again)

TryTargetBrickCoordinates:
	for {
		targetBrick := pile.Bricks[targetBrickIndex]
		allCoordinates := targetBrick.enumerateXYCoordinates()
		for _, currentCoordinate := range allCoordinates {
			for z := 0; z < targetBrick.Start.Z; z += 1 {
				currentCoordinate.Z = z

				if brickIndex, ok := pile.CoordinateMap[currentCoordinate]; ok && !pile.brickIsSupported(brickIndex) {
					// We have found a lower unsupported brick
					targetBrickIndex = brickIndex
					log.Trace().Int("NextUnsupportedBrickIndex", targetBrickIndex).Msg("FindUnimpededUnsupportedBrick")
					continue TryTargetBrickCoordinates
				} // If lower brick found
			} // For each Z moving upwards
		} // For each XY of brick coordinate

		// We have checked every coordinate from Z=0 upwards and found nothing, so this must be an unimpeded unsupported brick
		return targetBrickIndex, nil
	} // Try Target Brick Coords
}

func (pile BrickPileData) moveBrickToSupportedPosition(brickIndex int) {
	targetBrick := pile.Bricks[brickIndex]

	log.Trace().
		Int("BrickIndex", brickIndex).
		Str("Brick", targetBrick.String()).
		Msg("MoveBrickToSupportedPosition")

	log.Trace().Msg("Deleting Old Brick Coordinates From CoordinateMap")
	brickVolumeCoordinates := targetBrick.enumerateAllCoordinates()
	// Remove all the old brick coordinates from the map
	for _, coord := range brickVolumeCoordinates {
		delete(pile.CoordinateMap, coord)
	}

	log.Trace().Msg("Finding New Brick Z Coordinate")
	brickXYCoordinates := targetBrick.enumerateXYCoordinates()
	brickNewZ := targetBrick.Start.Z
newZCoordSearchLoop:
	for brickNewZ >= 0 {
		for _, coord := range brickXYCoordinates {
			coord.Z = brickNewZ
			// If we DO find a brick, we have found the limit to how far we can fall,
			// so increment Z once and leave the loop
			if foundBrickIndex, ok := pile.CoordinateMap[coord]; ok {
				log.Trace().
					Int("CurrentZ", brickNewZ).
					Interface("SharedXY", coord).
					Int("FoundBrickIndex", foundBrickIndex).
					Str("FoundBrick", pile.Bricks[foundBrickIndex].String()).
					Msg("Found Brick At Current Z")
				break newZCoordSearchLoop
			}
		}
		log.Trace().Int("ZProbeValue", brickNewZ).Msg("No Bricks Found at Z")

		// We did not find any bricks, we must decrement Z and try again
		brickNewZ -= 1
	}

	brickNewZ += 1
	log.Trace().Int("NewZ", brickNewZ).Send()

	// Update the actual brick data
	targetBrick.Start.Z = brickNewZ
	pile.Bricks[brickIndex] = targetBrick
	log.Trace().Str("NewBrick", targetBrick.String()).Send()

	log.Trace().Msg("Add New Brick Coordinates To CoordinateMap")
	// Update the CoordinateMap
	brickVolumeCoordinates = targetBrick.enumerateAllCoordinates()
	for _, coord := range brickVolumeCoordinates {
		log.Trace().Interface("NewCoord", coord).Send()
		pile.CoordinateMap[coord] = brickIndex
	}

	log.Trace().Msg("Update Support Array")
	// Update the support array
	supportingBricks := make([]int, 0)
	for _, coord := range brickXYCoordinates {
		coord.Z = targetBrick.Start.Z - 1
		if supportingIndex, ok := pile.CoordinateMap[coord]; ok && !slices.Contains(supportingBricks, supportingIndex) {
			supportingBricks = append(supportingBricks, supportingIndex)
		}
	}
	log.Trace().Interface("SupportingBricks", supportingBricks).Send()
	pile.Supports[brickIndex] = supportingBricks
}

func (pile BrickPileData) SimulateBrickFall() {
	// Until all bricks are supported (and hence cannot fall), we must:
	// 1) Identify a brick that can fall unimpeded by other bricks
	// 2) Place that brick in position and update it
	// 3) Mark that brick as supported, and note which bricks support it
	// 4) Repeat

	// Step 1) Identify a brick that can fall unimpeded.
	//
	// Idea: Choose a brick. Check if there are any unsupported bricks directly below it.
	// If there are, update our choice to that brick. Repeat until there are no bricks below us.
	// This must converge as there are finite bricks, so we cannot continue to
	// find bricks below us forever, and bricks are strictly rectangular.

	// SimulateBrickFallLoop:
	for {
		log.Debug().Msg("Brick Fall Loop Start")
		// Start by finding an unsupported brick
		unsupportedBrickIndex, err := pile.findUnimpededUnsupportedBrickIndex()
		if err != nil {
			// We have no more unsupported bricks!
			return
		}

		log.Debug().Msgf("Moving Brick %v", unsupportedBrickIndex)
		// Move this brick down under gravity
		pile.moveBrickToSupportedPosition(unsupportedBrickIndex)
	} // SimulateBrickFallLoop

}

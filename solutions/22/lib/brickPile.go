package lib

import (
	"bufio"
	"errors"

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
		CoordinateMap: make(map[Coordinate]int),
	}

	var line string
	for fileScanner.Scan() {
		line = fileScanner.Text()
		brick := ParseLineToBrick(line)
		brick.ID = len(pile.Bricks)

		pile.Bricks = append(pile.Bricks, brick)
		pile.Supports = append(pile.Supports, make([]int, 0))

		for _, coord := range brick.enumerateCoordinates() {
			if _, ok := pile.CoordinateMap[coord]; ok {
				log.Fatal().Msgf("brick collision at coordinate %v", coord)
			}
			pile.CoordinateMap[coord] = brick.ID
		}
	}

	return pile
}

func (pile BrickPileData) findUnimpededUnsupportedBrickIndex() (int, error) {
	targetBrickIndex := -1

	// Select the first brick that is unsupported
	for brickIndex, brickSupports := range pile.Supports {
		if len(brickSupports) == 0 {
			targetBrickIndex = brickIndex
			break
		}
	}

	if targetBrickIndex == -1 {
		return -1, errors.New("no unsupported brick found")
	}

	// For each X,Y coordinate of this brick, starting at the ground and moving up to the brick:
	// 1) a supported brick is found, (continue to the next coordinate)
	// 2) an unsupported brick is found (change targetBrickIndex and try again)

TryTargetBrickCoordinates:
	for {
		var currentCoordinate Coordinate
		targetBrick := pile.Bricks[targetBrickIndex]
		for dx := 0; dx < targetBrick.Lengths.X; dx += 1 {
			for dy := 0; dy < targetBrick.Lengths.Y; dy += 1 {
				for z := 0; z < targetBrick.Start.Z; z += 1 {
					currentCoordinate = Coordinate{
						X: targetBrick.Start.X + dx,
						Y: targetBrick.Start.Y + dy,
						Z: z,
					}

					if brickIndex, ok := pile.CoordinateMap[currentCoordinate]; ok {
						// We have found a lower brick
						targetBrickIndex = brickIndex
						continue TryTargetBrickCoordinates
					} // If lower brick found
				} // For each Z moving upwards
			} // For each dy
		} // For each dx

		// We have checked every coordinate from Z=0 upwards and found nothing, so this must be an unimpeded unsupported brick
		return targetBrickIndex, nil
	} // Try Target Brick Coords
}

func (pile BrickPileData) SimulateBrickFall() {
	// Until all bricks are supported (and hence cannot fall), we must:
	// 1) Identify a brick that can fall unimpeded by other bricks
	// 2) Place that brick in position and update it in the FallenBrickPile
	// 3) Mark that brick as supported, and note which bricks support it
	// 4) Repeat

	// Step 1) Identify a brick that can fall unimpeded.
	//
	// Idea: Choose a brick. Check if there are any unsupported bricks directly below it.
	// If there are, update our choice to that brick. Repeat until there are no bricks below us.
	// This must converge as there are finite bricks, so we cannot continue to
	// find bricks below us forever, and bricks are strictly rectangular.

	// SimulateBrickFallLoop
	for {
		// Start by finding an unsupported brick
		unsupportedBrickIndex, err := pile.findUnimpededUnsupportedBrickIndex()
		if err != nil {
			// We have no more unsupported bricks!
			return
		}

	} // SimulateBrickFallLoop

}

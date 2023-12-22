package lib

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
)

type BrickData struct {
	ID      int
	Start   Coordinate
	Lengths Coordinate
}

func ParseLineToBrick(line string) BrickData {
	lineFields := strings.FieldsFunc(line, func(r rune) bool {
		return r == '~'
	})
	brickStarts := StringToCoordinate(lineFields[0])
	brickEnds := StringToCoordinate(lineFields[1])
	parsedBrick := createBrickFromCoordinates(brickStarts, brickEnds)

	log.Debug().
		Str("RawLine", line).
		Interface("StringFields", lineFields).
		Str("ParsedBrick", parsedBrick.String()).
		Send()

	return parsedBrick
}

func createBrickFromCoordinates(brickStart, brickEnd Coordinate) BrickData {
	brickLens := Coordinate{
		X: brickEnd.X - brickStart.X + 1,
		Y: brickEnd.Y - brickStart.Y + 1,
		Z: brickEnd.Z - brickStart.Z + 1,
	}

	if brickLens.X < 0 || brickLens.Y < 0 || brickLens.Z < 0 {
		log.Fatal().Interface("BrickLens", brickLens).Msg("found negative brick len")
	}

	return BrickData{
		Start:   brickStart,
		Lengths: brickLens,
	}
}

func (brick BrickData) enumerateXYCoordinates() []Coordinate {
	area := brick.Lengths.X * brick.Lengths.Y
	allCoordinates := make([]Coordinate, area)
	count := 0
	for dx := 0; dx < brick.Lengths.X; dx += 1 {
		for dy := 0; dy < brick.Lengths.Y; dy += 1 {
			allCoordinates[count] = Coordinate{
				X: brick.Start.X + dx,
				Y: brick.Start.Y + dy,
				Z: brick.Start.Z,
			}
			count += 1
		}
	}

	return allCoordinates
}

func (brick BrickData) enumerateAllCoordinates() []Coordinate {
	volume := brick.Lengths.X * brick.Lengths.Y * brick.Lengths.Z
	allCoordinates := make([]Coordinate, volume)
	count := 0
	for dx := 0; dx < brick.Lengths.X; dx += 1 {
		for dy := 0; dy < brick.Lengths.Y; dy += 1 {
			for dz := 0; dz < brick.Lengths.Z; dz += 1 {
				allCoordinates[count] = Coordinate{
					X: brick.Start.X + dx,
					Y: brick.Start.Y + dy,
					Z: brick.Start.Z + dz,
				}
				count += 1
			}
		}
	}

	return allCoordinates
}

func (brick BrickData) String() string {
	volume := brick.Lengths.X * brick.Lengths.Y * brick.Lengths.Z
	return fmt.Sprintf("%v, %v (Vol: %v)", brick.Start.String(), brick.Lengths.String(), volume)
}

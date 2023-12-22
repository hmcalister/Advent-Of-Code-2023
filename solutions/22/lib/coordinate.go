package lib

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

type Coordinate struct {
	X int
	Y int
	Z int
}

func (coord Coordinate) String() string {
	return fmt.Sprintf("(%v, %v, %v)", coord.X, coord.Y, coord.Z)
}

func StringToCoordinate(str string) Coordinate {
	fields := strings.Split(str, ",")
	x, err := strconv.Atoi(fields[0])
	if err != nil {
		log.Fatal().Msgf("failed to convert field %v to integer in string %v", fields[0], str)
	}
	y, err := strconv.Atoi(fields[1])
	if err != nil {
		log.Fatal().Msgf("failed to convert field %v to integer in string %v", fields[1], str)
	}
	z, err := strconv.Atoi(fields[2])
	if err != nil {
		log.Fatal().Msgf("failed to convert field %v to integer in string %v", fields[2], str)
	}
	return Coordinate{
		x,
		y,
		z,
	}
}

package lib

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Vector struct {
	X float64
	Y float64
	Z float64
}

func (v Vector) String() string {
	return fmt.Sprintf("(%v, %v, %v)", v.X, v.Y, v.Z)
}

func parseFieldsToVector(fields []string) (Vector, error) {
	if len(fields) != 3 {
		return Vector{}, errors.New("cannot create vector with incorrect number of fields")
	}

	x, err := strconv.Atoi(strings.TrimSpace(fields[0]))
	if err != nil {
		return Vector{}, errors.New("failed to parse x coordinate")
	}

	y, err := strconv.Atoi(strings.TrimSpace(fields[1]))
	if err != nil {
		return Vector{}, errors.New("failed to parse y coordinate")
	}

	z, err := strconv.Atoi(strings.TrimSpace(fields[2]))
	if err != nil {
		return Vector{}, errors.New("failed to parse z coordinate")
	}

	return Vector{float64(x), float64(y), float64(z)}, nil
}

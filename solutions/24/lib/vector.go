package lib

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"gonum.org/v1/gonum/mat"
)

func vectorToString(v *mat.VecDense) string {
	// f := mat.Formatted(v, mat.FormatPython())
	return fmt.Sprintf("(%v, %v, %v)", v.AtVec(0), v.AtVec(1), v.AtVec(2))
}

func parseFieldsToVector(fields []string) (*mat.VecDense, error) {
	if len(fields) != 3 {
		return nil, errors.New("cannot create *mat.VecDense with incorrect number of fields")
	}

	x, err := strconv.Atoi(strings.TrimSpace(fields[0]))
	if err != nil {
		return nil, errors.New("failed to parse x coordinate")
	}

	y, err := strconv.Atoi(strings.TrimSpace(fields[1]))
	if err != nil {
		return nil, errors.New("failed to parse y coordinate")
	}

	z, err := strconv.Atoi(strings.TrimSpace(fields[2]))
	if err != nil {
		return nil, errors.New("failed to parse z coordinate")
	}

	return mat.NewVecDense(3, []float64{float64(x), float64(y), float64(z)}), nil
}

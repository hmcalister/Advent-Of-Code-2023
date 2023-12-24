package lib

import (
	"errors"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"gonum.org/v1/gonum/mat"
)

type HailstoneData struct {
	Position *mat.VecDense
	Velocity *mat.VecDense
}

func (hailstone HailstoneData) String() string {
	return fmt.Sprintf("%v, %v", vectorToString(hailstone.Position), vectorToString(hailstone.Velocity))
}

func parseLineToHailstone(line string) HailstoneData {
	ampersandIndex := strings.IndexRune(line, '@')
	if ampersandIndex == -1 {
		log.Fatal().Msgf("failed to find ampersand in line %v", line)
	}

	positionFields := strings.Split(line[:ampersandIndex], ", ")
	velocityFields := strings.Split(line[ampersandIndex+1:], ", ")

	positionVec, err := parseFieldsToVector(positionFields)
	if err != nil {
		log.Fatal().
			Str("Line", line).
			Interface("PositionFields", positionFields).
			Msg(err.Error())
	}
	velocityVec, err := parseFieldsToVector(velocityFields)
	if err != nil {
		log.Fatal().
			Str("Line", line).
			Interface("VelocityFields", velocityFields).
			Msg(err.Error())
	}

	return HailstoneData{
		Position: positionVec,
		Velocity: velocityVec,
	}
}

func (hailstone HailstoneData) FindPathIntersectionPositionInXY(secondHailstone HailstoneData) (*mat.VecDense, error) {
	// Find the point in space that the hailstone paths intersect.
	// If the paths do not intersect, return an error stating this.
	//
	// We can find this by solving the simultaneous equations in all of X, Y, and Z:
	// X1+t1*V1 = X2+t1*V2
	// ...
	//
	// Once we solve for t1 and t2, we can find the position easily

	A := mat.NewDense(2, 2, []float64{
		hailstone.Velocity.AtVec(0), -secondHailstone.Velocity.AtVec(0),
		hailstone.Velocity.AtVec(1), -secondHailstone.Velocity.AtVec(1),
	})
	b := mat.NewVecDense(2, []float64{
		secondHailstone.Position.AtVec(0) - hailstone.Position.AtVec(0),
		secondHailstone.Position.AtVec(1) - hailstone.Position.AtVec(1),
	})

	simultaneousEquationSolns := mat.NewVecDense(2, nil)
	err := simultaneousEquationSolns.SolveVec(A, b)
	if err != nil {
		return nil, err
	}
	log.Trace().
		Float64("HailstoneOneTimeIntersection", simultaneousEquationSolns.AtVec(0)).
		Float64("HailstoneTwoTimeIntersection", simultaneousEquationSolns.AtVec(1)).
		Send()

	if simultaneousEquationSolns.AtVec(0) < 0 || simultaneousEquationSolns.AtVec(1) < 0 {
		return nil, errors.New("paths cross with negative time")
	}

	pathIntersectionPosition := hailstone.calculatePositionAtTime(simultaneousEquationSolns.AtVec(0))

	return pathIntersectionPosition, nil
}

func (hailstone HailstoneData) FindPathIntersectionPosition(secondHailstone HailstoneData) (*mat.VecDense, error) {
	// Find the point in space that the hailstone paths intersect.
	// If the paths do not intersect, return an error stating this.
	//
	// We can find this by solving the simultaneous equations in all of X, Y, and Z:
	// X1+t1*V1 = X2+t1*V2
	// ...
	//
	// Once we solve for t1 and t2, we can find the position easily

	A := mat.NewDense(3, 2, []float64{
		hailstone.Velocity.AtVec(0), secondHailstone.Velocity.AtVec(0),
		hailstone.Velocity.AtVec(1), secondHailstone.Velocity.AtVec(1),
		hailstone.Velocity.AtVec(2), secondHailstone.Velocity.AtVec(2),
	})
	b := mat.NewVecDense(3, []float64{
		secondHailstone.Position.AtVec(0) - hailstone.Position.AtVec(0),
		secondHailstone.Position.AtVec(1) - hailstone.Position.AtVec(1),
		secondHailstone.Position.AtVec(2) - hailstone.Position.AtVec(2),
	})

	simultaneousEquationSolns := mat.NewVecDense(2, nil)
	err := simultaneousEquationSolns.SolveVec(A, b)
	if err != nil {
		return nil, err
	}

	pathIntersectionPosition := hailstone.calculatePositionAtTime(simultaneousEquationSolns.AtVec(0))

	return pathIntersectionPosition, nil
}

func (hailstone HailstoneData) FindCollisionTime(secondHailstone HailstoneData) (float64, error) {
	// Find the collision time for X and for Y, and if they are equal return it
	// Otherwise, return error
	//
	// We can find collision time by checking X1+V1*t = X2+V2*t,
	// And hence (X1-X2)/(V2-V1) = t

	if secondHailstone.Velocity.AtVec(0) == hailstone.Velocity.AtVec(0) {
		if secondHailstone.Position.AtVec(0) == hailstone.Position.AtVec(0) {
			return 0.0, nil
		} else {
			return 0.0, errors.New("no collision possible when x velocities equal")
		}
	}

	if secondHailstone.Velocity.AtVec(1) == hailstone.Velocity.AtVec(1) {
		if secondHailstone.Position.AtVec(1) == hailstone.Position.AtVec(1) {
			return 0.0, nil
		} else {
			return 0.0, errors.New("no collision possible when y velocities equal")
		}
	}

	if secondHailstone.Velocity.AtVec(2) == hailstone.Velocity.AtVec(2) {
		if secondHailstone.Position.AtVec(2) == hailstone.Position.AtVec(2) {
			return 0.0, nil
		} else {
			return 0.0, errors.New("no collision possible when y velocities equal")
		}
	}

	xCollisionTime := (hailstone.Position.AtVec(0) - secondHailstone.Position.AtVec(0)) / (secondHailstone.Velocity.AtVec(0) - hailstone.Velocity.AtVec(0))
	yCollisionTime := (hailstone.Position.AtVec(1) - secondHailstone.Position.AtVec(1)) / (secondHailstone.Velocity.AtVec(1) - hailstone.Velocity.AtVec(1))
	zCollisionTime := (hailstone.Position.AtVec(2) - secondHailstone.Position.AtVec(2)) / (secondHailstone.Velocity.AtVec(2) - hailstone.Velocity.AtVec(2))

	log.Trace().
		Str("Hailstone", hailstone.String()).
		Str("SecondHailstone", secondHailstone.String()).
		Float64("XCollisionTime", xCollisionTime).
		Float64("YCollisionTime", yCollisionTime).
		Float64("ZCollisionTime", zCollisionTime).
		Send()

	if xCollisionTime == yCollisionTime && xCollisionTime == zCollisionTime {
		return xCollisionTime, nil
	} else {
		return 0.0, errors.New("collision does not occur")
	}
}

func (hailstone HailstoneData) calculatePositionAtTime(time float64) *mat.VecDense {
	newPosition := mat.NewVecDense(3, nil)
	newPosition.AddScaledVec(hailstone.Position, time, hailstone.Velocity)
	return newPosition
}

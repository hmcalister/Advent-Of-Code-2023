package lib

import (
	"errors"
	"strings"

	"github.com/rs/zerolog/log"
)

type HailstoneData struct {
	Position Vector
	Velocity Vector
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

func (hailstone HailstoneData) FindCollisionTimeInXY(secondHailstone HailstoneData) (float64, error) {
	// Find the collision time for X and for Y, and if they are equal return it
	// Otherwise, return error
	//
	// We can find collision time by checking X1+V1*t = X2+V2*t,
	// And hence (X1-X2)/(V2-V1) = t

	if secondHailstone.Velocity.X == hailstone.Velocity.X {
		if secondHailstone.Position.X == hailstone.Position.X {
			return 0.0, nil
		} else {
			return 0.0, errors.New("no collision possible when x velocities equal")
		}
	}

	if secondHailstone.Velocity.Y == hailstone.Velocity.Y {
		if secondHailstone.Position.Y == hailstone.Position.Y {
			return 0.0, nil
		} else {
			return 0.0, errors.New("no collision possible when y velocities equal")
		}
	}

	xCollisionTime := (hailstone.Position.X - secondHailstone.Position.X) / (secondHailstone.Velocity.X - hailstone.Velocity.X)
	yCollisionTime := (hailstone.Position.Y - secondHailstone.Position.Y) / (secondHailstone.Velocity.Y - hailstone.Velocity.Y)

	log.Trace().
		Interface("Hailstone", hailstone).
		Interface("SecondHailstone", secondHailstone).
		Float64("XCollisionTime", xCollisionTime).
		Float64("YCollisionTime", yCollisionTime).
		Send()

	if xCollisionTime == yCollisionTime {
		return xCollisionTime, nil
	} else {
		return 0.0, errors.New("collision does not occur")
	}
}

func (hailstone HailstoneData) calculatePositionAtTime(time float64) Vector {
	return Vector{
		X: hailstone.Position.X + time*hailstone.Velocity.X,
		Y: hailstone.Position.Y + time*hailstone.Velocity.Y,
		Z: hailstone.Position.Z + time*hailstone.Velocity.Z,
	}
}

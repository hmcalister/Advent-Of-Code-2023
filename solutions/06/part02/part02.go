package part02

import (
	"bufio"
	"math"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

type raceDetails struct {
	raceID       int
	timeAllowed  int
	bestDistance int
}

func (rd raceDetails) calculateError() int {
	// First, find the roots of quadratic of when we can win the race
	lowerRoot := (float64(rd.timeAllowed) - math.Sqrt(float64(rd.timeAllowed*rd.timeAllowed-4*rd.bestDistance))) / 2
	upperRoot := (float64(rd.timeAllowed) + math.Sqrt(float64(rd.timeAllowed*rd.timeAllowed-4*rd.bestDistance))) / 2

	lowestAlpha := max(int(math.Floor(lowerRoot))+1, 0)
	upperAlpha := min(int(math.Ceil(upperRoot))-1, rd.timeAllowed)

	log.Debug().
		Float64("LowerRoot", lowerRoot).
		Float64("UpperRoot", upperRoot).
		Int("LowestAlpha", lowestAlpha).
		Int("UpperAlpha", upperAlpha).
		Send()
	return upperAlpha - lowestAlpha + 1
}

func ProcessInput(fileScanner *bufio.Scanner) (int, error) {
	var line string
	// Handle the Time line
	fileScanner.Scan()
	line = fileScanner.Text()
	line = strings.TrimSpace(line[5:])
	timeString := strings.Join(strings.Fields(line), "")

	log.Debug().
		Str("Line", line).
		Str("TimeString", timeString).
		Send()

	// Handle the distances line
	fileScanner.Scan()
	line = fileScanner.Text()
	line = strings.TrimSpace(line[9:])
	distanceString := strings.Join(strings.Fields(line), "")

	log.Debug().
		Str("Line", line).
		Str("DistanceString", distanceString).
		Send()

	time, err := strconv.Atoi(timeString)
	if err != nil {
		log.Fatal().Msgf("failed to parse time string (%v)", timeString)
	}

	distance, err := strconv.Atoi(distanceString)
	if err != nil {
		log.Fatal().Msgf("failed to parse distance string (%v)", distanceString)
	}

	rd := raceDetails{
		raceID:       0,
		timeAllowed:  time,
		bestDistance: distance,
	}
	result := rd.calculateError()
	log.Debug().
		Int("RaceID", rd.raceID).
		Int("CalculatedError", result).
		Send()

	return result, nil
}

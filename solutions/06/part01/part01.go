package part01

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
	log.Trace().
		Float64("LowerRoot", lowerRoot).
		Float64("UpperRoot", upperRoot).
		Int("LowestAlpha", lowestAlpha).
		Int("UpperAlpha", upperAlpha).
		Send()

	return upperAlpha - lowestAlpha + 1
}

func parseDataToRaceDetails(timeStrArr []string, distanceStrArr []string) []raceDetails {
	raceDetailsArr := make([]raceDetails, len(timeStrArr))
	for i := 0; i < len(timeStrArr); i += 1 {
		time, err := strconv.Atoi(timeStrArr[i])
		if err != nil {
			log.Fatal().Msgf("failed to parse time string for raceID %v (%v)", i, timeStrArr[i])
		}

		distance, err := strconv.Atoi(distanceStrArr[i])
		if err != nil {
			log.Fatal().Msgf("failed to parse distance string for raceID %v (%v)", i, distanceStrArr[i])
		}

		raceDetailsArr[i] = raceDetails{
			raceID:       i,
			timeAllowed:  time,
			bestDistance: distance,
		}
		log.Info().
			Int("ParsedRaceID", raceDetailsArr[i].raceID).
			Int("ParsedRaceTime", raceDetailsArr[i].timeAllowed).
			Int("ParsedRaceDistance", raceDetailsArr[i].bestDistance).
			Send()

	}

	return raceDetailsArr
}

func ProcessInput(fileScanner *bufio.Scanner) (int, error) {
	var line string
	// Handle the Time line
	fileScanner.Scan()
	line = fileScanner.Text()
	line = strings.TrimSpace(line[5:])
	timeStrings := strings.Fields(line)

	// Handle the distances line
	fileScanner.Scan()
	line = fileScanner.Text()
	line = strings.TrimSpace(line[9:])
	distanceStrings := strings.Fields(line)

	result := 1
	raceDetailsArray := parseDataToRaceDetails(timeStrings, distanceStrings)
	for _, rd := range raceDetailsArray {
		e := rd.calculateError()
		log.Debug().
			Int("RaceID", rd.raceID).
			Int("CalculatedError", e).
			Send()
		result *= e
	}

	return result, nil
}

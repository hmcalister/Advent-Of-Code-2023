package part01

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	MAX_RED   = 12
	MAX_GREEN = 13
	MAX_BLUE  = 14
)

type trialData struct {
	redCube   int
	greenCube int
	blueCube  int
}

// Given a string of a trial, put the data into an observationDataStruct
func makeTrialData(trialString string) (*trialData, error) {
	td := trialData{}

	// Each trial has a number of observation, separated by commas
	observations := strings.Split(trialString, ", ")
	for observationIndex, observation := range observations {
		observation = strings.Trim(observation, " ")
		observationNumber, _ := strconv.Atoi(observation[:strings.IndexByte(observation, ' ')])

		if strings.HasSuffix(observation, " red") {
			td.redCube = observationNumber
		} else if strings.HasSuffix(observation, " green") {
			td.greenCube = observationNumber
		} else if strings.HasSuffix(observation, " blue") {
			td.blueCube = observationNumber
		} else {
			return nil, fmt.Errorf("error in observation %v", observationIndex)
		}
	}

	return &td, nil
}

// Check if the given observation data is valid
func (td *trialData) isValid() bool {
	return td.redCube <= MAX_RED && td.greenCube <= MAX_GREEN && td.blueCube <= MAX_BLUE
}

func processLine(line string) (bool, error) {
	// Remove the "Game X:" prefix
	line = line[strings.IndexByte(line, ':')+1:]

	// Each trial is separated by a semicolon
	trials := strings.Split(line, "; ")

	for trialIndex, trial := range trials {
		td, err := makeTrialData(trial)
		if err != nil {
			return false, errors.Join(fmt.Errorf("error in trial %v", trialIndex), err)
		}
		if !td.isValid() {
			return false, nil
		}
	}

	return true, nil
}

// Given a scanner over an input file, return the sum of the GameIDs that
// satisfy the conditions of the number of red, green, and blue cubes in each
// trial (separated by semicolons).
func ProcessInput(fileScanner *bufio.Scanner) (int, error) {
	result := 0

	for fileScanner.Scan() {
		line := fileScanner.Text()
		gameID, _ := strconv.Atoi(line[5:strings.IndexByte(line, ':')])
		linePossible, err := processLine(line)
		if err != nil {
			return -1, fmt.Errorf("failed to process Game %v", gameID)
		}

		if linePossible {
			result += gameID
		}
	}

	return result, nil
}

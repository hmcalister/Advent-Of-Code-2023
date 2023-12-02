package part02

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type trialData struct {
	redCube   int
	greenCube int
	blueCube  int
}

type gameData struct {
	maxRedCube   int
	maxGreenCube int
	maxBlueCube  int
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

// Given some trial information, add this to the current game data information
//
// i.e. update the maximum known number of cubes
func (gd *gameData) addTrialInformation(td *trialData) {
	gd.maxRedCube = max(gd.maxRedCube, td.redCube)
	gd.maxGreenCube = max(gd.maxGreenCube, td.greenCube)
	gd.maxBlueCube = max(gd.maxBlueCube, td.blueCube)
}

// Calculate the power of the game data struct
func (gd *gameData) calculateGamePower() int {
	return gd.maxRedCube * gd.maxGreenCube * gd.maxBlueCube
}

// Given a line from the puzzle input, calculate the product of the
// maximum number of each red, green, and blue cubes.
func processLine(line string) (int, error) {
	// Remove the "Game X:" prefix
	line = line[strings.IndexByte(line, ':')+1:]

	// Each trial is separated by a semicolon
	trials := strings.Split(line, "; ")

	gd := &gameData{}

	for trialIndex, trial := range trials {
		td, err := makeTrialData(trial)
		if err != nil {
			return -1, errors.Join(fmt.Errorf("error in trial %v", trialIndex), err)
		}
		gd.addTrialInformation(td)
	}

	return gd.calculateGamePower(), nil
}

// Given a scanner over an input file, return the sum of the GameIDs that
// satisfy the conditions of the number of red, green, and blue cubes in each
// trial (separated by semicolons).
func ProcessInput(fileScanner *bufio.Scanner) (int, error) {
	result := 0

	for fileScanner.Scan() {
		line := fileScanner.Text()
		gameID, _ := strconv.Atoi(line[5:strings.IndexByte(line, ':')])
		gamePower, err := processLine(line)
		if err != nil {
			return -1, errors.Join(fmt.Errorf("failed to process Game %v", gameID), err)
		}

		result += gamePower
	}

	return result, nil
}

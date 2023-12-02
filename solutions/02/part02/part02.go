package part02

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/sync/errgroup"
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
			return nil, fmt.Errorf("error in observation %v: %v", observationIndex, observation)
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
			return -1, errors.Join(fmt.Errorf("error in trial %v: %v", trialIndex, trial), err)
		}
		gd.addTrialInformation(td)
	}

	return gd.calculateGamePower(), nil
}

// Given a scanner over an input file, return the sum of the GameIDs that
// satisfy the conditions of the number of red, green, and blue cubes in each
// trial (separated by semicolons).
func ProcessInput(fileScanner *bufio.Scanner) (int, error) {
	// Context for the error group that will process each line in parallel
	ctx := context.Background()
	errGroup, _ := errgroup.WithContext(ctx)

	// Channel for results of each game / line
	resultsChan := make(chan int)

	// Sum of all game powers, computed in goroutine below
	result := 0

	// WaitGroup to tell when summation has finished
	var resultSumWaitGroup sync.WaitGroup

	// Start the summation goroutine, use resultSumWaitGroup to tell when it is done
	resultSumWaitGroup.Add(1)
	go func() {
		for receivedResult := range resultsChan {
			result += receivedResult
		}
		resultSumWaitGroup.Done()
	}()

	// For each line in the file
	for fileScanner.Scan() {
		// Grab the text of this line
		line := fileScanner.Text()

		// Grab the game ID, we can assume the file is nicely formatted such that this doesn't error
		gameID, _ := strconv.Atoi(line[5:strings.IndexByte(line, ':')])

		// Start a goroutine in our errgroup to process this line
		errGroup.Go(func() error {
			// Variable shadowing to copy variables to within this closure
			gameID, line := gameID, line
			gamePower, err := processLine(line)
			if err != nil {
				return errors.Join(fmt.Errorf("failed to process Game %v", gameID), err)
			}

			resultsChan <- gamePower
			return nil
		})
	}

	// Wait for the parallel line processing to finish
	if err := errGroup.Wait(); err != nil {
		return -1, err
	}
	// Close the results channel, now result summation goroutine knows there is no more data
	close(resultsChan)
	// Finally, wait for results summation goroutine to finish (in case there was a bottleneck here)
	resultSumWaitGroup.Wait()

	return result, nil
}

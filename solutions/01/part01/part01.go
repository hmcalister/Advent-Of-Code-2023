package part01

import (
	"bufio"
	"fmt"
	"strconv"
	"unicode"
)

// Given a line of input, find and return the first digit of that line as a string
func FirstDigit(line string) (string, error) {
	for i := 0; i < len(line); i += 1 {
		c := rune(line[i])
		if unicode.IsDigit(c) {
			return string(c), nil
		}
	}

	return "", fmt.Errorf("no first digit found in line %v", line)
}

// Given a line of input, find and return the last digit of that line as a string
func LastDigit(line string) (string, error) {
	for i := len(line) - 1; i >= 0; i -= 1 {
		c := rune(line[i])
		if unicode.IsDigit(c) {
			return string(c), nil
		}
	}

	return "", fmt.Errorf("no last digit found in line %v", line)
}

// Given a scanner over some input stream, loop over each line,
// find the first and last digits, concatenate them, and sum the result.
//
// The sum of each of these lineNumbers is returned.
func ProcessInput(fileScanner *bufio.Scanner) (int, error) {
	result := 0
	lineNumber := 0

	for fileScanner.Scan() {
		line := fileScanner.Text()
		lineNumber += 1

		firstDigit, err := FirstDigit(line)
		if err != nil {
			return -1, fmt.Errorf("no first digit found on line %v: %v", lineNumber, line)
		}
		lastDigit, err := LastDigit(line)
		if err != nil {
			return -1, fmt.Errorf("no last digit found on line %v: %v", lineNumber, line)
		}

		lineDigitsString := fmt.Sprintf("%v%v", firstDigit, lastDigit)
		lineDigitNumber, err := strconv.Atoi(lineDigitsString)
		if err != nil {
			return -1, fmt.Errorf("line digits '%v' cannot be converted to int", lineDigitsString)
		}

		result += lineDigitNumber
	}

	return result, nil
}

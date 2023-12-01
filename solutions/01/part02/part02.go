package part02

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
)

var (
	NO_DIGIT_ERR       = errors.New("no digit at index")
	INDEX_OUT_OF_BOUND = errors.New("index out of bound for slice")

	stringToDigitHashmap = map[string]int{
		"0":     0,
		"1":     1,
		"2":     2,
		"3":     3,
		"4":     4,
		"5":     5,
		"6":     6,
		"7":     7,
		"8":     8,
		"9":     9,
		"zero":  0,
		"one":   1,
		"two":   2,
		"three": 3,
		"four":  4,
		"five":  5,
		"six":   6,
		"seven": 7,
		"eight": 8,
		"nine":  9,
	}
)

// Given a string and an index into that string, return the digit at that index
// (if it is valid) or an error
func testStringForDigit(line string, index int) (int, error) {
	if index < 0 || index > len(line) {
		return -1, INDEX_OUT_OF_BOUND
	}

	// Check if index is exactly a single digit
	currentCharacter := string(line[index])
	digit, ok := stringToDigitHashmap[currentCharacter]
	if ok {
		return digit, nil
	}

	// Possible string-digits can only be length 3, 4, 5
	for targetLen := 3; targetLen <= 5; targetLen += 1 {
		// Ensure we can actually take this slice
		if index+targetLen > len(line) {
			continue
		}

		currentSlice := line[index : index+targetLen]
		digit, ok := stringToDigitHashmap[currentSlice]
		if ok {
			return digit, nil
		}
	}

	// We have tried every combination, there is no valid digit
	return -1, NO_DIGIT_ERR
}

// Given a line of input, find and return the first digit of that line as a string
func getFirstDigit(line string) (int, error) {
	for i := 0; i < len(line); i += 1 {
		digit, err := testStringForDigit(line, i)
		if err == nil {
			return digit, nil
		}
	}

	return -1, NO_DIGIT_ERR
}

// Given a line of input, find and return the last digit of that line as a string
func getLastDigit(line string) (int, error) {
	for i := len(line) - 1; i >= 0; i -= 1 {
		digit, err := testStringForDigit(line, i)
		if err == nil {
			return digit, nil
		}
	}

	return -1, NO_DIGIT_ERR
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

		lineFirstDigit, err := getFirstDigit(line)
		if err != nil {
			return -1, fmt.Errorf("no first digit found on line %v: %v", lineNumber, line)
		}
		lineLastDigit, err := getLastDigit(line)
		if err != nil {
			return -1, fmt.Errorf("no last digit found on line %v: %v", lineNumber, line)
		}

		lineDigitsString := fmt.Sprintf("%v%v", lineFirstDigit, lineLastDigit)
		lineDigitNumber, err := strconv.Atoi(lineDigitsString)
		if err != nil {
			return -1, fmt.Errorf("line digits '%v' cannot be converted to int", lineDigitsString)
		}

		result += lineDigitNumber
	}

	return result, nil
}

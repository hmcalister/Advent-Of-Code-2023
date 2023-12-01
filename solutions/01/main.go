package main

import (
	"bufio"
	"hmcalister/aoc01/part01"
	"log"
	"os"
)

const INPUT_FILE_PATH = "puzzleInput"

func main() {
	file, err := os.Open(INPUT_FILE_PATH)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()

	fileScanner := bufio.NewScanner(file)
	result, err := part01.ProcessInput(fileScanner)
	if err != nil {
		log.Fatalf("error processing file input: %v", err)
	}

	log.Printf("result: %v", result)
}

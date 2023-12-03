package main

import (
	"bufio"
	"hmcalister/aoc03/part02"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const INPUT_FILE_PATH = "puzzleInput"

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	logFile, err := os.Create("log")
	if err != nil {
		log.Fatal().Msgf("Count not open log file: %v", err)
	}
	log.Logger = zerolog.New(logFile).With().Timestamp().Logger()
}

func main() {
	file, err := os.Open(INPUT_FILE_PATH)
	if err != nil {
		log.Fatal().Msgf("error opening file: %v", err)
	}
	defer file.Close()

	fileScanner := bufio.NewScanner(file)
	result, err := part02.ProcessInput(fileScanner)
	if err != nil {
		log.Fatal().Msgf("error processing file input: %v", err)
	}

	log.Info().
		Int("Result", result).
		Send()
}

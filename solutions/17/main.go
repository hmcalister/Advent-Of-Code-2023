package main

import (
	"bufio"
	"flag"
	"hmcalister/aoc17/part02"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const INPUT_FILE_PATH = "puzzleInput"

func init() {
	logToFileFlag := flag.Bool("logToFile", false, "Flag to log to file, rather than to console output")
	flag.Parse()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	if *logToFileFlag {
		logFile, err := os.Create("log")
		if err != nil {
			log.Fatal().Msgf("Count not open log file: %v", err)
		}
		log.Logger = zerolog.New(logFile).With().Timestamp().Logger()
	}
}

func main() {
	file, err := os.Open(INPUT_FILE_PATH)
	if err != nil {
		log.Panic().Msgf("error opening file: %v", err)
	}
	defer file.Close()

	fileScanner := bufio.NewScanner(file)
	result, err := part02.ProcessInput(fileScanner)
	if err != nil {
		log.Panic().Msgf("error processing file input: %v", err)
	}

	log.Info().
		Int("Result", result).
		Send()
}

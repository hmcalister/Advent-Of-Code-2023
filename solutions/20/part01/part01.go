package part01

import (
	"bufio"
	"hmcalister/aoc20/lib"

	"github.com/rs/zerolog/log"
)

func ProcessInput(fileScanner *bufio.Scanner) (int, error) {
	moduleConfig := lib.ParseFileToModuleConfiguration(fileScanner)

	for i := 0; i < 1000; i += 1 {
		moduleConfig.PushButton()
	}

	lowPulses := moduleConfig.TotalPulses[lib.LOW_PULSE]
	highPulses := moduleConfig.TotalPulses[lib.HIGH_PULSE]
	result := lowPulses * highPulses
	log.Debug().
		Int("LowPulseCount", lowPulses).
		Int("HighPulseCount", highPulses).
		Int("TotalPulses", lowPulses+highPulses).
		Int("Result", result).
		Send()

	return result, nil
}

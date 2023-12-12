package part01

import (
	"bufio"
	"strings"

	"github.com/rs/zerolog/log"
)

type Data struct {
}

func parseLineToData(line string) {
	fields := strings.Fields(line)
	log.Debug().
		Str("ParsedLine", line).
		Interface("Fields", fields).
		Msg("Parsing Line To Data")
}

func ProcessInput(fileScanner *bufio.Scanner) (int, error) {
	return 0, nil
}

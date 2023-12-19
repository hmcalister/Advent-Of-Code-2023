package part01

import (
	"bufio"
	"hmcalister/aoc19/lib"

	"github.com/rs/zerolog/log"
)

func ProcessInput(fileScanner *bufio.Scanner) (int, error) {
	var line string

	workflowMap := make(map[string]lib.Workflow)

	// Parse the workflows
	for fileScanner.Scan() {
		line = fileScanner.Text()
		if len(line) == 0 {
			break
		}
		log.Debug().
			Str("RawLine", line).
			Send()

		newWorkflow := lib.ParseLineToWorkflow(line)
		workflowMap[newWorkflow.WorkflowName] = newWorkflow

	}

	totalAcceptedRatings := 0
	// Parse the parts
	for fileScanner.Scan() {
		line = fileScanner.Text()
		part := lib.ParseLineToPartData(line)

		log.Debug().
			Str("RawLine", line).
			Interface("ParsedPart", part).
			Send()

		if lib.ProcessPart(part, workflowMap) {
			totalAcceptedRatings += part.SumRatings()
		}
	}

	return totalAcceptedRatings, nil
}

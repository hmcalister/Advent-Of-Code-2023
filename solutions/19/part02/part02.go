package part02

import (
	"bufio"
	"hmcalister/aoc19/part02/lib"

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

	allPartSpaceRanges := make([]lib.PartPropertySpaceRange, 0)
	totalAcceptedSpace := 0

	var currentPartSpaceRange lib.PartPropertySpaceRange
	var currentWorkflow lib.Workflow
	var ok bool
	allPartSpaceRanges = append(allPartSpaceRanges, lib.InitialPartPropertySpaceRange())

	for len(allPartSpaceRanges) > 0 {
		currentPartSpaceRange, allPartSpaceRanges = allPartSpaceRanges[0], allPartSpaceRanges[1:]
		log.Debug().
			Str("CurrentPartSpaceRange", currentPartSpaceRange.String()).
			Int("RemainingPartSpaceRanges", len(allPartSpaceRanges)).
			Send()

		currentWorkflow, ok = workflowMap[currentPartSpaceRange.NextWorkflow]
		if !ok {
			log.Fatal().Msgf("failed to find workflow with name %v for PartSpaceRange %v", currentPartSpaceRange.NextWorkflow, currentPartSpaceRange)
		}

		nextPartSpaceRanges := currentWorkflow.FindNextPartSpaceRanges(currentPartSpaceRange)
		for _, nextPartSpaceRange := range nextPartSpaceRanges {
			switch nextPartSpaceRange.NextWorkflow {
			case lib.ACCEPT_PART:
				log.Debug().
					Str("AcceptedRange", nextPartSpaceRange.String()).
					Int("RangeSize", nextPartSpaceRange.Size()).
					Send()
				totalAcceptedSpace += nextPartSpaceRange.Size()
			case lib.REJECT_PART:
				continue
			default:
				allPartSpaceRanges = append(allPartSpaceRanges, nextPartSpaceRange)
			}
		}
	}

	return totalAcceptedSpace, nil
}

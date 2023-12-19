package lib

import (
	"strings"

	"github.com/rs/zerolog/log"
)

const (
	REJECT_PART string = "R"
	ACCEPT_PART string = "A"
)

type Workflow struct {
	// The tag for *this* workflow, the name to address it by
	WorkflowName string

	// A list of functions for this workflow, ordered.
	// These functions are checked in order and the first match is used
	WorkflowFuncs []workflowFunction

	// The target for the workflow, matching with the WorkflowFuncs.
	// If WorkflowFuncs[i] is true, the part is sent to WorkflowTargets[i]
	//
	// Two special workflowNames exist, REJECT_PART and ACCEPT_PART.
	// This should be handled by the controller logic to count part accept/rejection
	WorkflowTargets []string
}

func ParseLineToWorkflow(line string) Workflow {
	nameEndIndex := strings.IndexRune(line, '{')
	workflowName := line[:nameEndIndex]

	workflowFunctionsSection := line[nameEndIndex+1 : len(line)-1]
	workflowFunctionsStrings := strings.FieldsFunc(workflowFunctionsSection, func(r rune) bool {
		return r == ','
	})

	workflowFuncs := make([]workflowFunction, len(workflowFunctionsStrings))
	workflowTargets := make([]string, len(workflowFunctionsStrings))

	for index, workflowFunctionString := range workflowFunctionsStrings {
		workflowFuncs[index], workflowTargets[index] = parseFieldToWorkflowFunction(workflowFunctionString)
	}

	log.Debug().
		Str("WorkflowName", workflowName).
		Int("NumberOfWorkflowFunctions", len(workflowFuncs)).
		Interface("WorkflowTargets", workflowTargets).
		Send()

	return Workflow{
		WorkflowName:    workflowName,
		WorkflowFuncs:   workflowFuncs,
		WorkflowTargets: workflowTargets,
	}
}

func (flow Workflow) FindNextPartSpaceRanges(currentPartSpaceRange PartPropertySpaceRange) []PartPropertySpaceRange {
	nextPartSpaceRanges := make([]PartPropertySpaceRange, 0)

	currentModifiedPartSpaceRange := currentPartSpaceRange

	for functionIndex, function := range flow.WorkflowFuncs {
		log.Debug().
			Str("WorkflowName", flow.WorkflowName).
			Int("FunctionIndex", functionIndex).
			Str("CurrentModifiedPartSpaceRange", currentModifiedPartSpaceRange.String()).
			Send()
		nextPartSpaceRange := currentModifiedPartSpaceRange

		nextPartSpaceRange.NextWorkflow = flow.WorkflowTargets[functionIndex]

		switch function.targetPartProperty {
		case ExtremelyCoolProperty:
			nextPartSpaceRange.ExtremelyCoolRating.RangeStart = function.passingRangeUpdateFunc(nextPartSpaceRange.ExtremelyCoolRating.RangeStart, function.comparisonValue)
			nextPartSpaceRange.ExtremelyCoolRating.RangeEnd = function.passingRangeUpdateFunc(nextPartSpaceRange.ExtremelyCoolRating.RangeEnd, function.comparisonValue)

			currentModifiedPartSpaceRange.ExtremelyCoolRating.RangeStart = function.failingRangeUpdateFunc(currentModifiedPartSpaceRange.ExtremelyCoolRating.RangeStart, function.comparisonValue)
			currentModifiedPartSpaceRange.ExtremelyCoolRating.RangeEnd = function.failingRangeUpdateFunc(currentModifiedPartSpaceRange.ExtremelyCoolRating.RangeEnd, function.comparisonValue)
		case MusicalProperty:
			nextPartSpaceRange.MusicalRating.RangeStart = function.passingRangeUpdateFunc(nextPartSpaceRange.MusicalRating.RangeStart, function.comparisonValue)
			nextPartSpaceRange.MusicalRating.RangeEnd = function.passingRangeUpdateFunc(nextPartSpaceRange.MusicalRating.RangeEnd, function.comparisonValue)

			currentModifiedPartSpaceRange.MusicalRating.RangeStart = function.failingRangeUpdateFunc(currentModifiedPartSpaceRange.MusicalRating.RangeStart, function.comparisonValue)
			currentModifiedPartSpaceRange.MusicalRating.RangeEnd = function.failingRangeUpdateFunc(currentModifiedPartSpaceRange.MusicalRating.RangeEnd, function.comparisonValue)
		case AerodynamicProperty:
			nextPartSpaceRange.AerodynamicRating.RangeStart = function.passingRangeUpdateFunc(nextPartSpaceRange.AerodynamicRating.RangeStart, function.comparisonValue)
			nextPartSpaceRange.AerodynamicRating.RangeEnd = function.passingRangeUpdateFunc(nextPartSpaceRange.AerodynamicRating.RangeEnd, function.comparisonValue)

			currentModifiedPartSpaceRange.AerodynamicRating.RangeStart = function.failingRangeUpdateFunc(currentModifiedPartSpaceRange.AerodynamicRating.RangeStart, function.comparisonValue)
			currentModifiedPartSpaceRange.AerodynamicRating.RangeEnd = function.failingRangeUpdateFunc(currentModifiedPartSpaceRange.AerodynamicRating.RangeEnd, function.comparisonValue)
		case ShinyProperty:
			nextPartSpaceRange.ShinyRating.RangeStart = function.passingRangeUpdateFunc(nextPartSpaceRange.ShinyRating.RangeStart, function.comparisonValue)
			nextPartSpaceRange.ShinyRating.RangeEnd = function.passingRangeUpdateFunc(nextPartSpaceRange.ShinyRating.RangeEnd, function.comparisonValue)

			currentModifiedPartSpaceRange.ShinyRating.RangeStart = function.failingRangeUpdateFunc(currentModifiedPartSpaceRange.ShinyRating.RangeStart, function.comparisonValue)
			currentModifiedPartSpaceRange.ShinyRating.RangeEnd = function.failingRangeUpdateFunc(currentModifiedPartSpaceRange.ShinyRating.RangeEnd, function.comparisonValue)
		}

		log.Trace().
			Str("ResultingRange", nextPartSpaceRange.String()).
			Send()

		if nextPartSpaceRange.Size() == 0 {
			continue
		}
		nextPartSpaceRanges = append(nextPartSpaceRanges, nextPartSpaceRange)
	}
	return nextPartSpaceRanges
}

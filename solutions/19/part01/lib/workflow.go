package lib

import (
	"strconv"
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

type workflowFunction func(PartData) bool

func constructWorkflowFunction(workflowProperty string, comparisonFunc comparisonFunctionType, comparisonValue int) workflowFunction {
	return func(part PartData) bool {
		workflowProperty := workflowProperty

		var partProperty int
		switch workflowProperty {
		case ExtremelyCoolString:
			partProperty = part.ExtremelyCoolRating
		case MusicalString:
			partProperty = part.MusicalRating
		case AerodynamicString:
			partProperty = part.AerodynamicRating
		case ShinyString:
			partProperty = part.ShinyRating
		default:
			log.Fatal().Msgf("failed to parse workflow property %v", workflowProperty)
		}
		return comparisonFunc(partProperty, comparisonValue)
	}
}

// Given a string defining a workflow function and target (something like a<1006:qkq)
// return a workflow function implementing this logic and the target of the function if true
func parseFieldToWorkflowFunction(field string) (workflowFunction, string) {
	colonIndex := strings.IndexRune(field, ':')

	// If there is no condition, we have a "true" match and the field is the target
	if colonIndex == -1 {
		log.Debug().
			Str("WorkflowTarget", field).
			Msg("ParsedWorkflow")

		return func(pd PartData) bool { return true }, field
	}

	// The target of the workflow if resolves to true
	workflowTarget := field[colonIndex+1:]

	// The part property of interest
	workflowProperty := field[:1]

	// The comparison operator to use in the workflow
	workflowComparison := field[1:2]

	// The value to compare to in the workflow
	workflowComparisonValueString := field[2:colonIndex]
	workflowComparisonValue, err := strconv.Atoi(workflowComparisonValueString)
	if err != nil {
		log.Fatal().Msgf("failed to parse workflow comparison value %v to int in workflow string %v", workflowComparisonValueString, field)
	}

	var comparisonFunc comparisonFunctionType
	switch workflowComparison {
	case "<":
		comparisonFunc = lessThanFunc
	case ">":
		comparisonFunc = greaterThanFunc
	default:
		log.Fatal().Msgf("failed to parse comparison rune %v in workflow string %v", workflowComparison, field)
	}

	log.Debug().
		Str("workflowProperty", workflowProperty).
		Str("WorkflowComparison", workflowComparison).
		Int("WorkflowComparisonValue", workflowComparisonValue).
		Str("WorkflowTarget", workflowTarget).
		Msg("ParsedWorkflow")

	return constructWorkflowFunction(workflowProperty, comparisonFunc, workflowComparisonValue), workflowTarget
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

// Pass a part through the workflow functions, returning the target string of the first match
func (flow Workflow) passPartThroughWorkflowFunctions(part PartData) string {
	for funcIndex, f := range flow.WorkflowFuncs {
		if f(part) {
			return flow.WorkflowTargets[funcIndex]
		}
	}

	// Default to returning last value, which should always be true anyway
	return flow.WorkflowTargets[len(flow.WorkflowTargets)-1]
}

// Given a part and a collection of workflows organized into a map (with start point having label "in"),
// process the part through each workflow until the target label is either ACCEPT (true) or REJECT (false)
func ProcessPart(part PartData, workflowMap map[string]Workflow) bool {
	currentWorkflow, ok := workflowMap["in"]
	if !ok {
		log.Fatal().Msg("no entry point workflow with label \"in\"")
	}

	for {
		nextWorkflowName := currentWorkflow.passPartThroughWorkflowFunctions(part)
		if nextWorkflowName == ACCEPT_PART {
			return true
		}
		if nextWorkflowName == REJECT_PART {
			return false
		}

		currentWorkflow, ok = workflowMap[nextWorkflowName]
		if !ok {
			log.Fatal().Msgf("no next workflow with label \"%v\"", nextWorkflowName)
		}
	}
}

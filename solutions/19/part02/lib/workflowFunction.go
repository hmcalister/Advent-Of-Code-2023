package lib

import (
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

type workflowFunction struct {
	targetPartProperty     partPropertyEnum
	comparisonType         comparisonTypeEnum
	passingRangeUpdateFunc comparisonFunctionType
	failingRangeUpdateFunc comparisonFunctionType
	comparisonValue        int
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

		passingRangeUpdateFunc, failingRangeUpdateFunc := getUpdateFuncs(GREATER_THAN)
		return workflowFunction{
			targetPartProperty:     ExtremelyCoolProperty,
			comparisonType:         GREATER_THAN,
			passingRangeUpdateFunc: passingRangeUpdateFunc,
			failingRangeUpdateFunc: failingRangeUpdateFunc,
			comparisonValue:        0,
		}, field
	}

	// The target of the workflow if resolves to true
	workflowTarget := field[colonIndex+1:]

	// The part property of interest
	workflowProperty := partPropertyEnum(field[:1][0])

	// The comparison operator to use in the workflow
	workflowComparison := comparisonTypeEnum(field[1:2][0])

	// The value to compare to in the workflow
	workflowComparisonValueString := field[2:colonIndex]
	workflowComparisonValue, err := strconv.Atoi(workflowComparisonValueString)
	if err != nil {
		log.Fatal().Msgf("failed to parse workflow comparison value %v to int in workflow string %v", workflowComparisonValueString, field)
	}

	passingRangeUpdateFunc, failingRangeUpdateFunc := getUpdateFuncs(workflowComparison)

	log.Debug().
		Str("workflowProperty", workflowProperty.String()).
		Str("WorkflowComparison", workflowComparison.String()).
		Int("WorkflowComparisonValue", workflowComparisonValue).
		Str("WorkflowTarget", workflowTarget).
		Msg("ParsedWorkflow")

	return workflowFunction{
		targetPartProperty:     workflowProperty,
		comparisonType:         workflowComparison,
		passingRangeUpdateFunc: passingRangeUpdateFunc,
		failingRangeUpdateFunc: failingRangeUpdateFunc,
		comparisonValue:        workflowComparisonValue,
	}, workflowTarget
}

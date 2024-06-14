//nolint:mnd
package calculator

import (
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type MathDivisionOperation struct {
	options *sdk.OperationInfo
}

func NewDivision() *MathDivisionOperation {
	return &MathDivisionOperation{
		options: &sdk.OperationInfo{
			Name:        "Division",
			Description: "Divide first number with the second number",
			RequireAuth: false,
			Input:       inputFields,
			SampleOutput: map[string]interface{}{
				"name":          "calculator-division",
				"usage_mode":    "operation",
				"first_number":  50,
				"second_number": 2,
				"result":        25,
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *MathDivisionOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	input := sdk.InputToType[CalculatorOperationProps](ctx)

	fNumber := input.FirstNumber
	sNumber := input.SecondNumber

	if sNumber == 0 {
		return map[string]interface{}{
			"name":          "calculator-division",
			"usage_mode":    "operation",
			"first_number":  fNumber,
			"second_number": sNumber,
			"result":        "Unknown",
		}, nil
	}

	return map[string]interface{}{
		"name":          "calculator-division",
		"usage_mode":    "operation",
		"first_number":  fNumber,
		"second_number": sNumber,
		"result":        fNumber / sNumber,
	}, nil
}

func (c *MathDivisionOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *MathDivisionOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}

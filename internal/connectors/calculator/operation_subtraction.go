//nolint:mnd
package calculator

import (
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type MathSubtractionOperation struct {
	options *sdk.OperationInfo
}

func NewSubtraction() *MathSubtractionOperation {
	return &MathSubtractionOperation{
		options: &sdk.OperationInfo{
			Name:        "Subtraction",
			Description: "Subtracts second number from first number",
			RequireAuth: false,
			Input:       inputFields,
			SampleOutput: map[string]interface{}{
				"name":          "calculator-subtraction",
				"usage_mode":    "operation",
				"first_number":  76.23,
				"second_number": 16.03,
				"result":        60.2,
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *MathSubtractionOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	input := sdk.InputToType[CalculatorOperationProps](ctx)

	fNumber := input.FirstNumber
	sNumber := input.SecondNumber

	return map[string]interface{}{
		"name":          "calculator-subtraction",
		"usage_mode":    "operation",
		"first_number":  fNumber,
		"second_number": sNumber,
		"result":        fNumber - sNumber,
	}, nil
}

func (c *MathSubtractionOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *MathSubtractionOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}

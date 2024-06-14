//nolint:mnd
package calculator

import (
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type MathMultiplicationOperation struct {
	options *sdk.OperationInfo
}

func NewMultiplication() *MathMultiplicationOperation {
	return &MathMultiplicationOperation{
		options: &sdk.OperationInfo{
			Name:        "Multiplication",
			Description: "Multiplies two numbers",
			RequireAuth: false,
			Input:       inputFields,
			SampleOutput: map[string]interface{}{
				"name":          "calculator-multiplication",
				"usage_mode":    "operation",
				"first_number":  45.63,
				"second_number": 67.23,
				"result":        3067.7049,
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *MathMultiplicationOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	input := sdk.InputToType[CalculatorOperationProps](ctx)

	fNumber := input.FirstNumber
	sNumber := input.SecondNumber

	return map[string]interface{}{
		"name":          "calculator-multiplication",
		"usage_mode":    "operation",
		"first_number":  fNumber,
		"second_number": sNumber,
		"result":        fNumber * sNumber,
	}, nil
}

func (c *MathMultiplicationOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *MathMultiplicationOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}

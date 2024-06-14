package calculator

import (
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type MathAdditionOperation struct {
	options *sdk.OperationInfo
}

func NewAddition() *MathAdditionOperation {
	return &MathAdditionOperation{
		options: &sdk.OperationInfo{
			Name:        "Addition",
			Description: "Adds up two numbers",
			RequireAuth: false,
			Input:       inputFields,
			SampleOutput: map[string]interface{}{
				"name":          "calculator-addition",
				"usage_mode":    "operation",
				"first_number":  23,
				"second_number": 37,
				"result":        60,
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *MathAdditionOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	input := sdk.InputToType[CalculatorOperationProps](ctx)

	fNumber := input.FirstNumber
	sNumber := input.SecondNumber

	return map[string]interface{}{
		"name":          "calculator-addition",
		"usage_mode":    "operation",
		"first_number":  fNumber,
		"second_number": sNumber,
		"result":        fNumber + sNumber,
		"teste":         ctx.ResolvedInput,
	}, nil
}

func (c *MathAdditionOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *MathAdditionOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}

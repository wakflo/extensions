//nolint:mnd
package calculator

import (
	"fmt"
	"math"

	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type MathModuloOperation struct {
	options *sdk.OperationInfo
}

func NewModulo() *MathModuloOperation {
	return &MathModuloOperation{
		options: &sdk.OperationInfo{
			Name:        "Modulo",
			Description: "Remainder of a division between first number with the second number",
			RequireAuth: false,
			Input:       inputFields,
			SampleOutput: map[string]interface{}{
				"name":          "calculator-modulo",
				"usage_mode":    "operation",
				"first_number":  5,
				"second_number": 3,
				"result":        2,
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *MathModuloOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	input := sdk.InputToType[CalculatorOperationProps](ctx)

	fmt.Println(input)

	fNumber := input.FirstNumber
	sNumber := input.SecondNumber

	output := math.Mod(fNumber, sNumber)
	if math.IsNaN(output) {
		output = 0
	}

	return map[string]interface{}{
		"name":          "calculator-modulo",
		"usage_mode":    "operation",
		"first_number":  fNumber,
		"second_number": sNumber,
		"result":        output,
	}, nil
}

func (c *MathModuloOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *MathModuloOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}

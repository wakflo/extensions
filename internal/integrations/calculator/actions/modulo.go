package actions

import (
	"math"

	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type moduloActionProps struct {
	sharedProps
}

type ModuloAction struct{}

func (a *ModuloAction) Name() string {
	return "Modulo"
}

func (a *ModuloAction) Description() string {
	return "The Modulo integration action allows you to perform arithmetic operations on numeric values within your workflow. You can add, subtract, multiply, or divide two numbers and store the result in a variable. This action is useful when you need to manipulate data or calculate totals within your automated process."
}

func (a *ModuloAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *ModuloAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &moduloDocs,
	}
}

func (a *ModuloAction) Icon() *string {
	return nil
}

func (a *ModuloAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return inputFields
}

func (a *ModuloAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[moduloActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	fNumber := input.FirstNumber
	sNumber := input.SecondNumber

	output := math.Mod(fNumber, sNumber)
	if math.IsNaN(output) {
		output = 0
	}

	return map[string]interface{}{
		"firstNumber":  fNumber,
		"secondNumber": sNumber,
		"result":       output,
	}, nil
}

func (a *ModuloAction) Auth() *sdk.Auth {
	return nil
}

func (a *ModuloAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"firstNumber":  5,
		"secondNumber": 3,
		"result":       2,
	}
}

func (a *ModuloAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewModuloAction() sdk.Action {
	return &ModuloAction{}
}

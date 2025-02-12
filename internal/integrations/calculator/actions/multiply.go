package actions

import (
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type multiplyActionProps struct {
	sharedProps
}

type MultiplyAction struct{}

func (a *MultiplyAction) Name() string {
	return "Multiply"
}

func (a *MultiplyAction) Description() string {
	return "Multiplies two input values together and returns the result."
}

func (a *MultiplyAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *MultiplyAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &multiplyDocs,
	}
}

func (a *MultiplyAction) Icon() *string {
	return nil
}

func (a *MultiplyAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return inputFields
}

func (a *MultiplyAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[multiplyActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	fNumber := input.FirstNumber
	sNumber := input.SecondNumber

	return map[string]interface{}{
		"firstNumber":  fNumber,
		"secondNumber": sNumber,
		"result":       fNumber * sNumber,
	}, nil
}

func (a *MultiplyAction) Auth() *sdk.Auth {
	return nil
}

func (a *MultiplyAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"firstNumber":  45.63,
		"secondNumber": 67.23,
		"result":       3067.7049,
	}
}

func (a *MultiplyAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewMultiplyAction() sdk.Action {
	return &MultiplyAction{}
}

package actions

import (
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type subtractActionProps struct {
	sharedProps
}

type SubtractAction struct{}

func (a *SubtractAction) Name() string {
	return "Subtract"
}

func (a *SubtractAction) Description() string {
	return "Subtracts one or more values from another value in a specified field, allowing you to perform arithmetic operations and manipulate data within your workflows."
}

func (a *SubtractAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *SubtractAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &subtractDocs,
	}
}

func (a *SubtractAction) Icon() *string {
	return nil
}

func (a *SubtractAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return inputFields
}

func (a *SubtractAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[subtractActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	fNumber := input.FirstNumber
	sNumber := input.SecondNumber

	return map[string]interface{}{
		"firstNumber":  fNumber,
		"secondNumber": sNumber,
		"result":       fNumber - sNumber,
	}, nil
}

func (a *SubtractAction) Auth() *sdk.Auth {
	return nil
}

func (a *SubtractAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"firstNumber":  76.23,
		"secondNumber": 16.03,
		"result":       60.2,
	}
}

func (a *SubtractAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewSubtractAction() sdk.Action {
	return &SubtractAction{}
}

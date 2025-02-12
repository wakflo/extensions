package actions

import (
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type addActionProps struct {
	sharedProps
}

type AddAction struct{}

func (a *AddAction) Name() string {
	return "Add"
}

func (a *AddAction) Description() string {
	return "Adds one or more items to a specified collection or list within your workflow. This action is useful when you need to append new data to an existing dataset or update a list of values in your workflow."
}

func (a *AddAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *AddAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &addDocs,
	}
}

func (a *AddAction) Icon() *string {
	return nil
}

func (a *AddAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return inputFields
}

func (a *AddAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[addActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	fNumber := input.FirstNumber
	sNumber := input.SecondNumber

	return map[string]interface{}{
		"firstNumber":  fNumber,
		"secondNumber": sNumber,
		"result":       fNumber + sNumber,
	}, nil
}

func (a *AddAction) Auth() *sdk.Auth {
	return nil
}

func (a *AddAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"first_number":  23,
		"second_number": 37,
		"result":        60,
	}
}

func (a *AddAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewAddAction() sdk.Action {
	return &AddAction{}
}

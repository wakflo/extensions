package actions

import (
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type divideActionProps struct {
	sharedProps
}

type DivideAction struct{}

func (a *DivideAction) Name() string {
	return "Divide"
}

func (a *DivideAction) Description() string {
	return "Splits an input array or string into multiple equal parts based on a specified number of segments."
}

func (a *DivideAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *DivideAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &divideDocs,
	}
}

func (a *DivideAction) Icon() *string {
	return nil
}

func (a *DivideAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return inputFields
}

func (a *DivideAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[divideActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	fNumber := input.FirstNumber
	sNumber := input.SecondNumber

	if sNumber == 0 {
		return map[string]interface{}{
			"firstNumber":  fNumber,
			"secondNumber": sNumber,
			"result":       "Unknown",
		}, nil
	}

	return map[string]interface{}{
		"firstNumber":  fNumber,
		"secondNumber": sNumber,
		"result":       fNumber / sNumber,
	}, nil
}

func (a *DivideAction) Auth() *sdk.Auth {
	return nil
}

func (a *DivideAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"firstNumber":  50,
		"secondNumber": 2,
		"result":       25,
	}
}

func (a *DivideAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewDivideAction() sdk.Action {
	return &DivideAction{}
}

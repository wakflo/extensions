package actions

import (
	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type divideActionProps struct {
	sharedProps
}

type DivideAction struct{}

// Metadata returns metadata about the action
func (a *DivideAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "divide",
		DisplayName:   "Divide",
		Description:   "Splits an input array or string into multiple equal parts based on a specified number of segments.",
		Type:          core.ActionTypeAction,
		Documentation: divideDocs,
		SampleOutput: map[string]any{
			"firstNumber":  50,
			"secondNumber": 2,
			"result":       25,
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *DivideAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("divide", "Divide")

	form.NumberField("firstNumber", "First number").
		Placeholder("Enter first number").
		Required(true).
		HelpText("The number to be divided (dividend)")

	form.NumberField("secondNumber", "Second number").
		Placeholder("Enter second number").
		Required(true).
		HelpText("The number to divide by (divisor)")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *DivideAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *DivideAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[divideActionProps](ctx)
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

func NewDivideAction() sdk.Action {
	return &DivideAction{}
}

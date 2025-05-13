package actions

import (
	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type multiplyActionProps struct {
	sharedProps
}

type MultiplyAction struct{}

// Metadata returns metadata about the action
func (a *MultiplyAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "multiply",
		DisplayName:   "Multiply",
		Description:   "Multiplies two input values together and returns the result.",
		Type:          core.ActionTypeAction,
		Documentation: multiplyDocs,
		SampleOutput: map[string]any{
			"firstNumber":  45.63,
			"secondNumber": 67.23,
			"result":       3067.7049,
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *MultiplyAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("multiply", "Multiply")

	form.NumberField("firstNumber", "First number").
		Placeholder("Enter first number").
		Required(true).
		HelpText("The first number to multiply")

	form.NumberField("secondNumber", "Second number").
		Placeholder("Enter second number").
		Required(true).
		HelpText("The second number to multiply")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *MultiplyAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *MultiplyAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[multiplyActionProps](ctx)
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

func NewMultiplyAction() sdk.Action {
	return &MultiplyAction{}
}

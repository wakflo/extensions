package actions

import (
	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type subtractActionProps struct {
	sharedProps
}

type SubtractAction struct{}

// Metadata returns metadata about the action
func (a *SubtractAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "subtract",
		DisplayName:   "Subtract",
		Description:   "Subtracts one or more values from another value in a specified field, allowing you to perform arithmetic operations and manipulate data within your workflows.",
		Type:          core.ActionTypeAction,
		Documentation: subtractDocs,
		SampleOutput: map[string]any{
			"firstNumber":  76.23,
			"secondNumber": 16.03,
			"result":       60.2,
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *SubtractAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("subtract", "Subtract")

	form.NumberField("firstNumber", "First number").
		Placeholder("Enter first number").
		Required(true).
		HelpText("The number to subtract from")

	form.NumberField("secondNumber", "Second number").
		Placeholder("Enter second number").
		Required(true).
		HelpText("The number to subtract")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *SubtractAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *SubtractAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[subtractActionProps](ctx)
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

func NewSubtractAction() sdk.Action {
	return &SubtractAction{}
}

package actions

import (
	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type addActionProps struct {
	sharedProps
}

type AddAction struct{}

// Metadata returns metadata about the action
func (a *AddAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "add",
		DisplayName:   "Add",
		Description:   "Adds one or more items to a specified collection or list within your workflow. This action is useful when you need to append new data to an existing dataset or update a list of values in your workflow.",
		Type:          core.ActionTypeAction,
		Documentation: addDocs,
		SampleOutput: map[string]any{
			"first_number":  23,
			"second_number": 37,
			"result":        60,
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *AddAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("add", "Add")

	form.NumberField("firstNumber", "First number").
		Placeholder("Enter first number").
		Required(true).
		HelpText("The first number to add")

	form.NumberField("secondNumber", "Second number").
		Placeholder("Enter second number").
		Required(true).
		HelpText("The second number to add")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *AddAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *AddAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[addActionProps](ctx)
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

func NewAddAction() sdk.Action {
	return &AddAction{}
}

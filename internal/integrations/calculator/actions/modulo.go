package actions

import (
	"math"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type moduloActionProps struct {
	sharedProps
}

type ModuloAction struct{}

// Metadata returns metadata about the action
func (a *ModuloAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "modulo",
		DisplayName:   "Modulo",
		Description:   "The Modulo integration action allows you to perform arithmetic operations on numeric values within your workflow. You can add, subtract, multiply, or divide two numbers and store the result in a variable. This action is useful when you need to manipulate data or calculate totals within your automated process.",
		Type:          core.ActionTypeAction,
		Documentation: moduloDocs,
		SampleOutput: map[string]any{
			"firstNumber":  5,
			"secondNumber": 3,
			"result":       2,
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *ModuloAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("modulo", "Modulo")

	form.NumberField("firstNumber", "First number").
		Placeholder("Enter first number").
		Required(true).
		HelpText("The dividend")

	form.NumberField("secondNumber", "Second number").
		Placeholder("Enter second number").
		Required(true).
		HelpText("The divisor")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *ModuloAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *ModuloAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[moduloActionProps](ctx)
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

func NewModuloAction() sdk.Action {
	return &ModuloAction{}
}

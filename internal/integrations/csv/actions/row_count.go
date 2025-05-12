package actions

import (
	"github.com/gocarina/gocsv"
	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type rowCountActionProps struct {
	Content string `json:"content"`
}

type RowCountAction struct{}

// Metadata returns metadata about the action
func (a *RowCountAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "row_count",
		DisplayName:   "Row Count",
		Description:   "The Row Count integration action counts the number of rows in a specified table or dataset and returns the result as an output variable. This action is useful for tracking changes to data sets, monitoring data growth, or verifying data integrity.",
		Type:          core.ActionTypeAction,
		Documentation: rowCountDocs,
		SampleOutput: map[string]any{
			"rowCount": 0,
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *RowCountAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("row_count", "Row Count")

	// Add content field
	form.TextField("content", "CSV Content").
		Placeholder("csv content").
		Required(true).
		HelpText("The CSV content to count rows in")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *RowCountAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *RowCountAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[rowCountActionProps](ctx)
	if err != nil {
		return nil, err
	}

	var rows [][]string
	err = gocsv.UnmarshalString(input.Content, &rows)
	if err != nil {
		return nil, err
	}

	return map[string]any{"rowCount": len(rows)}, nil
}

func NewRowCountAction() sdk.Action {
	return &RowCountAction{}
}

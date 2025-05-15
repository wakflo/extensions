package actions

import (
	"context"
	"errors"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type createSpreadsheetActionProps struct {
	Name string `json:"name"`
}

type CreateSpreadsheetAction struct{}

// Metadata returns metadata about the action
func (a *CreateSpreadsheetAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_spreadsheet",
		DisplayName:   "Create Spreadsheet",
		Description:   "Create a new spreadsheet in Google Sheets or Microsoft Excel with customizable settings such as sheet name, row and column count, and formatting options.",
		Type:          core.ActionTypeAction,
		Documentation: createSpreadsheetDocs,
		Icon:          "",
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *CreateSpreadsheetAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_spreadsheet", "Create Spreadsheet")

	form.TextField("name", "name").
		Placeholder("Sheet Name").
		HelpText("The name of the sheet.").
		Required(true)

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *CreateSpreadsheetAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *CreateSpreadsheetAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[createSpreadsheetActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	sheetService, err := sheets.NewService(context.Background(), option.WithTokenSource(*authCtx.TokenSource))
	if err != nil {
		return nil, err
	}

	if input.Name == "" {
		return nil, errors.New("name is required")
	}

	document, err := sheetService.Spreadsheets.Create(&sheets.Spreadsheet{
		Properties: &sheets.SpreadsheetProperties{
			Title: input.Name,
		},
	}).
		Do()
	return document, err
}

func NewCreateSpreadsheetAction() sdk.Action {
	return &CreateSpreadsheetAction{}
}

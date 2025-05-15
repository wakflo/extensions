package actions

import (
	"context"
	"errors"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/googlesheets/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type copyWorksheetActionProps struct {
	SpreadSheetID            string `json:"spreadsheetId"`
	DestinationSpreadSheetID string `json:"destinationSpreadSheetId"`
	SheetID                  string `json:"sheetId"`
}

type CopyWorksheetAction struct{}

// Metadata returns metadata about the action
func (a *CopyWorksheetAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "copy_worksheet",
		DisplayName:   "Copy Worksheet",
		Description:   "Copies an existing worksheet to a new location within the same or different workbook, allowing you to duplicate and reuse worksheets with ease.",
		Type:          core.ActionTypeAction,
		Documentation: copyWorksheetDocs,
		Icon:          "",
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *CopyWorksheetAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("copy_worksheet", "Copy Worksheet")

	shared.RegisterSpreadsheetsProps(form, "destinationSpreadSheetId", "SpreadSheet", "Select Spreadsheet", true)

	shared.RegisterSheetIDProps(form, true)

	shared.RegisterSpreadsheetsProps(form, "destinationSpreadSheetId", "Destination SpreadSheet", "Destination spreadsheet to copy to", true)

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *CopyWorksheetAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *CopyWorksheetAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[copyWorksheetActionProps](ctx)
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

	sheetID := shared.ConvertToInt64(input.SheetID)

	spreadsheetCopy, err := sheetService.Spreadsheets.Sheets.CopyTo(input.SpreadSheetID, sheetID, &sheets.CopySheetToAnotherSpreadsheetRequest{
		DestinationSpreadsheetId: input.DestinationSpreadSheetID,
	}).Do()
	if err != nil {
		return nil, err
	}

	if spreadsheetCopy == nil {
		return nil, errors.New("received nil response from the CopyTo operation")
	}

	return spreadsheetCopy, err
}

func NewCopyWorksheetAction() sdk.Action {
	return &CopyWorksheetAction{}
}

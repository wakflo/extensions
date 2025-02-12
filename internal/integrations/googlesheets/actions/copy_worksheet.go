package actions

import (
	"context"
	"errors"

	"github.com/wakflo/extensions/internal/integrations/googlesheets/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type copyWorksheetActionProps struct {
	SpreadSheetID            string `json:"spreadsheetId"`
	DestinationSpreadSheetID string `json:"destinationSpreadSheetId"`
	SheetID                  string `json:"sheetId"`
}

type CopyWorksheetAction struct{}

func (a *CopyWorksheetAction) Name() string {
	return "Copy Worksheet"
}

func (a *CopyWorksheetAction) Description() string {
	return "Copies an existing worksheet to a new location within the same or different workbook, allowing you to duplicate and reuse worksheets with ease."
}

func (a *CopyWorksheetAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *CopyWorksheetAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &copyWorksheetDocs,
	}
}

func (a *CopyWorksheetAction) Icon() *string {
	return nil
}

func (a *CopyWorksheetAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"spreadsheetId":            shared.GetSpreadsheetsInput("Spreadsheet", "spreadsheet ID", true),
		"sheetId":                  shared.GetSheetIDInput("Sheet", "select sheet", true),
		"destinationSpreadSheetId": shared.GetSpreadsheetsInput("Destination Spreadsheet", "Destination spreadsheet to copy to", true),
	}
}

func (a *CopyWorksheetAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[copyWorksheetActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	sheetService, err := sheets.NewService(context.Background(), option.WithTokenSource(*ctx.Auth.TokenSource))
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

func (a *CopyWorksheetAction) Auth() *sdk.Auth {
	return nil
}

func (a *CopyWorksheetAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *CopyWorksheetAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewCopyWorksheetAction() sdk.Action {
	return &CopyWorksheetAction{}
}

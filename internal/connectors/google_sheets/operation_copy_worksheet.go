package googlesheets

import (
	"context"
	"errors"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"

	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type copyWorkSheetOperationProps struct {
	SpreadSheetID            string `json:"spreadsheetId"`
	DestinationSpreadSheetID string `json:"destinationSpreadSheetId"`
	SheetID                  string `json:"sheetId"`
}

type CopyWorkSheetOperation struct {
	options *sdk.OperationInfo
}

func NewCopyWorkSheetOperation() *CopyWorkSheetOperation {
	return &CopyWorkSheetOperation{
		options: &sdk.OperationInfo{
			Name:        "Copy Worksheet",
			Description: "Copy a worksheet into another or the same spreadsheet",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"spreadsheetId":            getSpreadsheetsInput("Spreadsheet", "spreadsheet ID", true),
				"sheetId":                  getSheetIDInput("Sheet", "select sheet", true),
				"destinationSpreadSheetId": getSpreadsheetsInput("Destination Spreadsheet", "Destination spreadsheet to copy to", true),
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *CopyWorkSheetOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Token == nil {
		return nil, errors.New("missing google auth token")
	}

	if ctx.Auth.TokenSource == nil {
		return nil, errors.New("missing google token source")
	}

	input, err := sdk.InputToTypeSafely[copyWorkSheetOperationProps](ctx)
	if err != nil {
		return nil, err
	}
	sheetService, err := sheets.NewService(context.Background(), option.WithTokenSource(*ctx.Auth.TokenSource))
	if err != nil {
		return nil, err
	}

	sheetID := convertToInt64(input.SheetID)

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

func (c *CopyWorkSheetOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *CopyWorkSheetOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}

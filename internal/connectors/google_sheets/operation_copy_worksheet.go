package googlesheets

import (
	"context"
	"errors"
	"strconv"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type copyWorkSheetOperationProps struct {
	SpreadSheetID            string `json:"spreadSheetId"`
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
				"spreadSheetId": autoform.NewShortTextField().
					SetDisplayName("Spreadsheet ID").
					SetDescription("The ID of the spreadsheet.").
					SetRequired(true).
					Build(),
				"sheetId": autoform.NewShortTextField().
					SetDisplayName("Sheet ID").
					SetDescription("The ID of the sheet.").
					SetRequired(true).
					Build(),
				"destinationSpreadSheetId": autoform.NewShortTextField().
					SetDisplayName("Destination Spreadsheet ID").
					SetDescription("The ID of the destination spreadsheet.").
					SetRequired(true).
					Build(),
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

	input := sdk.InputToType[copyWorkSheetOperationProps](ctx)
	sheetService, err := sheets.NewService(context.Background(), option.WithTokenSource(*ctx.Auth.TokenSource))
	if err != nil {
		return nil, err
	}

	if input.SpreadSheetID == "" {
		return nil, errors.New("spreadsheet ID is required")
	}

	if input.SheetID == "" {
		return nil, errors.New("sheet ID is required")
	}

	if input.SheetID == "" {
		return nil, errors.New("sheet ID is required")
	}

	// this converts the sheet ID to int as required by the Sheets API
	sheetIDInt, err := strconv.ParseInt(input.SheetID, 10, 64)
	if err != nil {
		return nil, err
	}

	spreadsheetCopy, err := sheetService.Spreadsheets.Sheets.CopyTo(input.SpreadSheetID, sheetIDInt, &sheets.CopySheetToAnotherSpreadsheetRequest{
		DestinationSpreadsheetId: input.DestinationSpreadSheetID,
	}).Do()

	return spreadsheetCopy, err
}

func (c *CopyWorkSheetOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *CopyWorkSheetOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}

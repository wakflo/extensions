package googlesheets

import (
	"context"
	"errors"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"

	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type findWorkSheetByIDOperationProps struct {
	SpreadSheetID     string `json:"spreadsheetId,omitempty"`
	SheetID           string `json:"sheetId"`
	IncludeTeamDrives bool   `json:"includeTeamDrives"`
}

type FindWorkSheetByIDOperation struct {
	options *sdk.OperationInfo
}

func NewFindWorkSheetByIDOperation() *FindWorkSheetByIDOperation {
	return &FindWorkSheetByIDOperation{
		options: &sdk.OperationInfo{
			Name:        "Find Worksheet By ID",
			Description: "Find a worksheet by ID",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"spreadsheetId":     getSpreadsheetsInput("Spreadsheet", "spreadsheet ID", true),
				"sheetId":           getSheetIDInput("Sheet", "select sheet", true),
				"includeTeamDrives": includeTeamFieldInput,
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *FindWorkSheetByIDOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Token == nil {
		return nil, errors.New("missing google auth token")
	}

	input := sdk.InputToType[findWorkSheetByIDOperationProps](ctx)
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

	spreadsheet, err := sheetService.Spreadsheets.Get(input.SpreadSheetID).
		Do()

	for _, sheet := range spreadsheet.Sheets {
		if sheet.Properties.SheetId == convertToInt64(input.SheetID) {
			return sheet, nil
		}
	}
	return spreadsheet, err
}

func (c *FindWorkSheetByIDOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *FindWorkSheetByIDOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}

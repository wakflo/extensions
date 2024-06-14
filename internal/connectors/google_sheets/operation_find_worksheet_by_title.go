package googlesheets

import (
	"context"
	"errors"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type findWorkSheetByTitleOperationProps struct {
	SpreadSheetID string `json:"spreadSheetId"`
	SheetTitle    string `json:"sheetTitle"`
}

type FindWorkSheetByTitleOperation struct {
	options *sdk.OperationInfo
}

func NewFindWorkSheetByTitleOperation() *FindWorkSheetByTitleOperation {
	return &FindWorkSheetByTitleOperation{
		options: &sdk.OperationInfo{
			Name:        "Find Worksheet By Title",
			Description: "Find a worksheet by title",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"spreadSheetId": autoform.NewShortTextField().
					SetDisplayName("Spreadsheet ID").
					SetDescription("The ID of the spreadsheet.").
					SetRequired(true).
					Build(),
				"sheetTitle": autoform.NewShortTextField().
					SetDisplayName("Sheet Title").
					SetDescription("The title of the sheet.").
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

func (c *FindWorkSheetByTitleOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Token == nil {
		return nil, errors.New("missing google auth token")
	}

	input := sdk.InputToType[findWorkSheetByTitleOperationProps](ctx)
	sheetService, err := sheets.NewService(context.Background(), option.WithTokenSource(*ctx.Auth.TokenSource))
	if err != nil {
		return nil, err
	}

	if input.SpreadSheetID == "" {
		return nil, errors.New("spreadsheet ID is required")
	}

	if input.SheetTitle == "" {
		return nil, errors.New("sheet title is required")
	}

	spreadsheet, err := sheetService.Spreadsheets.Get(input.SpreadSheetID).
		// Fields("id, name, mimeType, webViewLink, kind, createdTime").
		Do()

	for _, sheet := range spreadsheet.Sheets {
		if sheet.Properties.Title == input.SheetTitle {
			return sheet, nil
		}

		if sheet.Properties.Title != input.SheetTitle {
			const notFound string = "not sheet found"
			return notFound, nil
		}
	}
	return spreadsheet, err
}

func (c *FindWorkSheetByTitleOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *FindWorkSheetByTitleOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}

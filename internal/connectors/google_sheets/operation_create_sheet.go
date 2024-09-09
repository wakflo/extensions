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

type createSheetOperationProps struct {
	Name string `json:"name"`
}

type CreateSheetOperation struct {
	options *sdk.OperationInfo
}

func NewCreateSheetOperation() *CreateSheetOperation {
	return &CreateSheetOperation{
		options: &sdk.OperationInfo{
			Name:        "Create Sheet",
			Description: "Create a new sheet",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"name": autoform.NewShortTextField().
					SetDisplayName("Sheet Name").
					SetDescription("The name of the sheet.").
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

func (c *CreateSheetOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Token == nil {
		return nil, errors.New("missing google auth token")
	}

	input := sdk.InputToType[createSheetOperationProps](ctx)
	sheetService, err := sheets.NewService(context.Background(), option.WithTokenSource(*ctx.Auth.TokenSource))
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

func (c *CreateSheetOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *CreateSheetOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}

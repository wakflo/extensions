package actions

import (
	"context"
	"errors"

	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type createSpreadsheetActionProps struct {
	Name string `json:"name"`
}

type CreateSpreadsheetAction struct{}

func (a *CreateSpreadsheetAction) Name() string {
	return "Create Spreadsheet"
}

func (a *CreateSpreadsheetAction) Description() string {
	return "Create a new spreadsheet in Google Sheets or Microsoft Excel with customizable settings such as sheet name, row and column count, and formatting options."
}

func (a *CreateSpreadsheetAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *CreateSpreadsheetAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &createSpreadsheetDocs,
	}
}

func (a *CreateSpreadsheetAction) Icon() *string {
	return nil
}

func (a *CreateSpreadsheetAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"name": autoform.NewShortTextField().
			SetDisplayName("Sheet Name").
			SetDescription("The name of the sheet.").
			SetRequired(true).
			Build(),
	}
}

func (a *CreateSpreadsheetAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createSpreadsheetActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

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

func (a *CreateSpreadsheetAction) Auth() *sdk.Auth {
	return nil
}

func (a *CreateSpreadsheetAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *CreateSpreadsheetAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewCreateSpreadsheetAction() sdk.Action {
	return &CreateSpreadsheetAction{}
}

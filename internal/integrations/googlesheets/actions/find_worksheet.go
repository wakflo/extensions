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

type findWorksheetActionProps struct {
	SpreadSheetID string `json:"spreadsheetId,omitempty"`
	SheetTitle    string `json:"sheetTitle"`
}

type FindWorksheetAction struct{}

func (a *FindWorksheetAction) Name() string {
	return "Find Worksheet"
}

func (a *FindWorksheetAction) Description() string {
	return "Locates and retrieves a specific worksheet from a spreadsheet application, allowing you to automate tasks that rely on this worksheet's data."
}

func (a *FindWorksheetAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *FindWorksheetAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &findWorksheetDocs,
	}
}

func (a *FindWorksheetAction) Icon() *string {
	return nil
}

func (a *FindWorksheetAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"spreadsheetId": shared.GetSpreadsheetsInput("Spreadsheet", "spreadsheet ID", true),
		"sheetTitle":    shared.GetSheetTitleInput("Sheet", "select sheet", true),
	}
}

func (a *FindWorksheetAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[findWorksheetActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

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
		Do()
	if err != nil {
		return nil, err
	}

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

func (a *FindWorksheetAction) Auth() *sdk.Auth {
	return nil
}

func (a *FindWorksheetAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *FindWorksheetAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewFindWorksheetAction() sdk.Action {
	return &FindWorksheetAction{}
}

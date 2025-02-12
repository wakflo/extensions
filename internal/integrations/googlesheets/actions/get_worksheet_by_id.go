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

type getWorksheetByIdActionProps struct {
	SpreadSheetID     string `json:"spreadsheetId,omitempty"`
	SheetID           string `json:"sheetId"`
	IncludeTeamDrives bool   `json:"includeTeamDrives"`
}

type GetWorksheetByIdAction struct{}

func (a *GetWorksheetByIdAction) Name() string {
	return "Get Worksheet By ID"
}

func (a *GetWorksheetByIdAction) Description() string {
	return "Retrieves a specific worksheet by its unique identifier (ID), allowing you to access and manipulate its contents within your workflow."
}

func (a *GetWorksheetByIdAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *GetWorksheetByIdAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &getWorksheetByIdDocs,
	}
}

func (a *GetWorksheetByIdAction) Icon() *string {
	return nil
}

func (a *GetWorksheetByIdAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"spreadsheetId":     shared.GetSpreadsheetsInput("Spreadsheet", "spreadsheet ID", true),
		"sheetId":           shared.GetSheetIDInput("Sheet", "select sheet", true),
		"includeTeamDrives": shared.IncludeTeamFieldInput,
	}
}

func (a *GetWorksheetByIdAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[getWorksheetByIdActionProps](ctx.BaseContext)
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

	if input.SheetID == "" {
		return nil, errors.New("sheet ID is required")
	}

	spreadsheet, err := sheetService.Spreadsheets.Get(input.SpreadSheetID).
		Do()
	if err != nil {
		return nil, err
	}

	for _, sheet := range spreadsheet.Sheets {
		if sheet.Properties.SheetId == shared.ConvertToInt64(input.SheetID) {
			return sheet, nil
		}
	}
	return spreadsheet, err
}

func (a *GetWorksheetByIdAction) Auth() *sdk.Auth {
	return nil
}

func (a *GetWorksheetByIdAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *GetWorksheetByIdAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewGetWorksheetByIdAction() sdk.Action {
	return &GetWorksheetByIdAction{}
}

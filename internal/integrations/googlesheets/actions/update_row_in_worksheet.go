package actions

import (
	"context"
	"errors"
	"fmt"

	"github.com/wakflo/extensions/internal/integrations/googlesheets/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type updateRowInWorksheetActionProps struct {
	SpreadSheetID string          `json:"spreadSheetId"`
	SheetTitle    string          `json:"sheetTitle"`
	SheetRow      string          `json:"sheetRow"`
	Values        [][]interface{} `json:"values"`
}

type UpdateRowInWorksheetAction struct{}

func (a *UpdateRowInWorksheetAction) Name() string {
	return "Update Row in Worksheet"
}

func (a *UpdateRowInWorksheetAction) Description() string {
	return "Updates a specific row in a worksheet by modifying its values based on the provided data. This action allows you to dynamically update existing rows in your worksheet, making it easy to maintain and refresh your data."
}

func (a *UpdateRowInWorksheetAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *UpdateRowInWorksheetAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &updateRowInWorksheetDocs,
	}
}

func (a *UpdateRowInWorksheetAction) Icon() *string {
	return nil
}

func (a *UpdateRowInWorksheetAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"spreadsheetId": shared.GetSpreadsheetsInput("Spreadsheet", "spreadsheet ID", true),
		"sheetTitle":    shared.GetSheetTitleInput("Sheet", "select sheet", true),
		"sheetRow": autoform.NewShortTextField().
			SetDisplayName("Sheet Row").
			SetDescription("The row range of the sheet in the format of A1 notation.").
			SetRequired(true).
			Build(),
		"values": autoform.NewArrayField().
			SetDisplayName("Values").
			SetItems(
				autoform.NewArrayField().
					SetDisplayName("Values").
					SetItems(
						autoform.NewShortTextField().
							SetDisplayName("Value").
							SetDescription("value").
							SetRequired(true).
							Build(),
					).
					SetDescription("The values to be updated in the row.").
					SetRequired(false).
					Build(),
			).
			SetDescription("The values to be updated in the row.").
			SetRequired(false).
			Build(),
	}
}

func (a *UpdateRowInWorksheetAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[updateRowInWorksheetActionProps](ctx.BaseContext)
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
		return nil, errors.New("sheet Title is required")
	}

	if input.SheetRow == "" {
		return nil, errors.New("sheet row range is required")
	}

	range_ := fmt.Sprintf("%s!%s", input.SheetTitle, input.SheetRow)

	// sampleData := [][]interface{}{{"apple", "banana", "orange", "pineapple", "mango"}}

	spreadsheet, err := sheetService.Spreadsheets.Values.Update(input.SpreadSheetID, range_, &sheets.ValueRange{
		Range:          range_,
		MajorDimension: "ROWS",
		Values:         input.Values,
	}).
		ValueInputOption("RAW").
		Do()
	return spreadsheet, err
}

func (a *UpdateRowInWorksheetAction) Auth() *sdk.Auth {
	return nil
}

func (a *UpdateRowInWorksheetAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *UpdateRowInWorksheetAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewUpdateRowInWorksheetAction() sdk.Action {
	return &UpdateRowInWorksheetAction{}
}

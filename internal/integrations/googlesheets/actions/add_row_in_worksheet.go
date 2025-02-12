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

type addRowInWorksheetActionProps struct {
	SpreadSheetID string          `json:"spreadSheetId"`
	SheetTitle    string          `json:"sheetTitle"`
	SheetRow      string          `json:"sheetRow"`
	Values        [][]interface{} `json:"values"`
}

type AddRowInWorksheetAction struct{}

func (a *AddRowInWorksheetAction) Name() string {
	return "Add Row in Worksheet"
}

func (a *AddRowInWorksheetAction) Description() string {
	return "Adds a new row to an existing worksheet, allowing you to dynamically update your data and workflows with fresh information. This integration action enables seamless data manipulation, making it easy to append new records, track changes, or perform calculations based on updated data."
}

func (a *AddRowInWorksheetAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *AddRowInWorksheetAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &addRowInWorksheetDocs,
	}
}

func (a *AddRowInWorksheetAction) Icon() *string {
	return nil
}

func (a *AddRowInWorksheetAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"spreadsheetId": shared.GetSpreadsheetsInput("Spreadsheet", "spreadsheet ID", true),
		"sheetTitle":    shared.GetSheetTitleInput("Sheet", "select sheet", true),
		"sheetRow": autoform.NewShortTextField().
			SetDisplayName("Sheet Row").
			SetDescription("The row range of the sheet in the format of A1 notation where you want to append the new row.").
			SetRequired(true).
			Build(),
		"values": autoform.NewArrayField().
			SetDisplayName("Values").
			SetItems(
				autoform.NewShortTextField().
					SetDisplayName("Label").
					SetDescription("Label").
					SetRequired(true).
					Build(),
			).
			SetDescription("The values to be added in the row you're about to create.").
			SetRequired(false).
			Build(),
	}
}

func (a *AddRowInWorksheetAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[addRowInWorksheetActionProps](ctx.BaseContext)
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

	// if len(input.Values) <= 0 {
	// 	return nil, errors.New("values are required")
	// }

	// this variable concatenates the sheet ID and the sheet row range since we want to append a row in a particular range of the worksheet
	// range_ := input.SheetID + input.SheetRow
	range_ := fmt.Sprintf("%s!%s", input.SheetTitle, input.SheetRow)

	// sample array data to add as new row
	sampleData := [][]interface{}{{"apple", "banana", "orange", "pineapple", "mango"}}

	spreadsheet, err := sheetService.Spreadsheets.Values.Append(input.SpreadSheetID, range_, &sheets.ValueRange{
		Range:          range_,
		MajorDimension: "ROWS",
		Values:         sampleData,
	}).
		ValueInputOption("RAW").
		InsertDataOption("INSERT_ROWS").
		Do()
	return spreadsheet, err
}

func (a *AddRowInWorksheetAction) Auth() *sdk.Auth {
	return nil
}

func (a *AddRowInWorksheetAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *AddRowInWorksheetAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewAddRowInWorksheetAction() sdk.Action {
	return &AddRowInWorksheetAction{}
}

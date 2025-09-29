package actions

import (
	"context"
	"errors"
	"fmt"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/googlesheets/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type ValueItem struct {
	Value string `json:"value"`
}

type ValuesWrapper struct {
	Values []ValueItem `json:"values"`
}

type updateRowInWorksheetActionProps struct {
	SpreadSheetID string        `json:"spreadsheetId"`
	SheetTitle    string        `json:"sheetTitle"`
	SheetRow      string        `json:"sheetRow"`
	Values        ValuesWrapper `json:"values"`
}

type UpdateRowInWorksheetAction struct{}

// Metadata returns metadata about the action
func (a *UpdateRowInWorksheetAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "update_row_in_worksheet",
		DisplayName:   "Update Row in Worksheet",
		Description:   "Updates a specific row in a worksheet by modifying its values based on the provided data. This action allows you to dynamically update existing rows in your worksheet, making it easy to maintain and refresh your data.",
		Type:          core.ActionTypeAction,
		Documentation: updateRowInWorksheetDocs,
		Icon:          "",
		SampleOutput: map[string]any{
			"updatedRange":   "Sheet1!A2:E2",
			"updatedRows":    1,
			"updatedColumns": 5,
			"updatedCells":   5,
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
// Keep it simple like AddRowInWorksheetAction
func (a *UpdateRowInWorksheetAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("update_row_in_worksheet", "Update Row in Worksheet")

	shared.RegisterSpreadsheetsProps(form, "spreadsheetId", "Spreadsheet", "Select Spreadsheet", true)

	shared.RegisterSheetTitleProps(form, true)

	form.TextField("sheetRow", "sheetRow").
		Placeholder("Sheet Row").
		HelpText("The row range of the sheet in the format of A1 notation (e.g., A2:E2 to update row 2).").
		Required(true)

	// Use the same simple structure as AddRowInWorksheetAction
	labelsArray := form.ArrayField("values", "Values")
	labelGroup := labelsArray.ObjectTemplate("values", "")
	labelGroup.TextField("value", "Value").
		Placeholder("Value").
		Required(false).
		HelpText("The value of the cell.")

	schema := form.Build()

	return schema
}

func (a *UpdateRowInWorksheetAction) Auth() *core.AuthMetadata {
	return nil
}

func (a *UpdateRowInWorksheetAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[updateRowInWorksheetActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	sheetService, err := sheets.NewService(context.Background(), option.WithTokenSource(*authCtx.TokenSource))
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

	// Convert the Values structure to [][]interface{} for Google Sheets API
	var rowData []interface{}
	for _, item := range input.Values.Values {
		rowData = append(rowData, item.Value)
	}

	// If no values provided, return error
	if len(rowData) == 0 {
		return nil, errors.New("at least one value is required")
	}

	// Wrap in another array since we're updating one row
	sheetData := [][]interface{}{rowData}

	range_ := fmt.Sprintf("%s!%s", input.SheetTitle, input.SheetRow)

	spreadsheet, err := sheetService.Spreadsheets.Values.Update(input.SpreadSheetID, range_, &sheets.ValueRange{
		Range:          range_,
		MajorDimension: "ROWS",
		Values:         sheetData,
	}).
		ValueInputOption("RAW").
		Do()
	if err != nil {
		return nil, fmt.Errorf("failed to update row: %w", err)
	}

	// Return meaningful response
	return map[string]interface{}{
		"success":        true,
		"updatedRange":   spreadsheet.UpdatedRange,
		"updatedRows":    spreadsheet.UpdatedRows,
		"updatedColumns": spreadsheet.UpdatedColumns,
		"updatedCells":   spreadsheet.UpdatedCells,
		"spreadsheetId":  spreadsheet.SpreadsheetId,
	}, nil
}

func NewUpdateRowInWorksheetAction() sdk.Action {
	return &UpdateRowInWorksheetAction{}
}

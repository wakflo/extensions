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

type updateRowInWorksheetActionProps struct {
	SpreadSheetID string          `json:"spreadSheetId"`
	SheetTitle    string          `json:"sheetTitle"`
	SheetRow      string          `json:"sheetRow"`
	Values        [][]interface{} `json:"values"`
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
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *UpdateRowInWorksheetAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("update_row_in_worksheet", "Update Row in Worksheet")

	shared.RegisterSpreadsheetsProps(form, "spreadSheetId", "Spreadsheet", "spreadsheet ID", true)

	shared.RegisterSheetTitleProps(form, true)

	form.TextField("sheetRow", "sheetRow").
		Placeholder("Sheet Row").
		HelpText("The row range of the sheet in the format of A1 notation.").
		Required(true)

	valuesArray := form.ArrayField("values", "Values")
	valuesArray.HelpText("The values to be updated in the row.")

	rowGroup := valuesArray.ObjectTemplate("row", "Row")

	innerValuesArray := rowGroup.ArrayField("innerValues", "Values")
	innerValuesArray.HelpText("Column values for this row")

	valueGroup := innerValuesArray.ObjectTemplate("value", "Value")
	valueGroup.TextField("value", "Value").
		Placeholder("Enter cell value").
		Required(false).
		HelpText("Individual cell value")

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

	range_ := fmt.Sprintf("%s!%s", input.SheetTitle, input.SheetRow)

	spreadsheet, err := sheetService.Spreadsheets.Values.Update(input.SpreadSheetID, range_, &sheets.ValueRange{
		Range:          range_,
		MajorDimension: "ROWS",
		Values:         input.Values,
	}).
		ValueInputOption("RAW").
		Do()
	return spreadsheet, err
}

func NewUpdateRowInWorksheetAction() sdk.Action {
	return &UpdateRowInWorksheetAction{}
}

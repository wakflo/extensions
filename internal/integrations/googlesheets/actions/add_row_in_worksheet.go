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

type addRowInWorksheetActionProps struct {
	SpreadSheetID string          `json:"spreadSheetId"`
	SheetTitle    string          `json:"sheetTitle"`
	SheetRow      string          `json:"sheetRow"`
	Values        [][]interface{} `json:"values"`
}

type AddRowInWorksheetAction struct{}

// Metadata returns metadata about the action
func (a *AddRowInWorksheetAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "add_row_in_worksheet",
		DisplayName:   "Add Row in Worksheet",
		Description:   "Adds a new row to an existing worksheet, allowing you to dynamically update your data and workflows with fresh information. This integration action enables seamless data manipulation, making it easy to append new records, track changes, or perform calculations based on updated data.",
		Type:          core.ActionTypeAction,
		Documentation: addRowInWorksheetDocs,
		Icon:          "",
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *AddRowInWorksheetAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("add_row_in_worksheet", "Add Row in Worksheet")

	shared.RegisterSpreadsheetsProps(form, "spreadsheetId", "SpreadSheet", "Select Spreadsheet", true)

	shared.RegisterSheetTitleProps(form, true)

	form.TextField("sheetRow", "sheetRow").
		Placeholder("Sheet Row").
		HelpText("The row range of the sheet in the format of A1 notation where you want to append the new row.").
		Required(true)

	labelsArray := form.ArrayField("values", "Values")
	labelGroup := labelsArray.ObjectTemplate("values", "")
	labelGroup.TextField("value", "Value").
		Placeholder("Value").
		Required(false).
		HelpText("The Value of the row.")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *AddRowInWorksheetAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *AddRowInWorksheetAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[addRowInWorksheetActionProps](ctx)
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

	// this variable concatenates the sheet ID and the sheet row range since we want to append a row in a particular range of the worksheet
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

func NewAddRowInWorksheetAction() sdk.Action {
	return &AddRowInWorksheetAction{}
}

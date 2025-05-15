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

type readRowInWorksheetActionProps struct {
	SpreadSheetID     string `json:"spreadsheetId,omitempty"`
	SheetID           string `json:"sheetId"`
	SheetTitle        string `json:"sheetTitle"`
	Range             string `json:"range"`
	StartRow          string `json:"start-row"`
	EndRow            string `json:"end-row"`
	IncludeTeamDrives bool   `json:"includeTeamDrives"`
}

type ReadRowInWorksheetAction struct{}

// Metadata returns metadata about the action
func (a *ReadRowInWorksheetAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "read_row_in_worksheet",
		DisplayName:   "Read Row in Worksheet",
		Description:   "Reads a single row from a worksheet and returns its values as an object. This action is useful when you need to retrieve specific data from a worksheet or perform actions based on the contents of a particular row.",
		Type:          core.ActionTypeAction,
		Documentation: readRowInWorksheetDocs,
		Icon:          "",
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *ReadRowInWorksheetAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("read_row_in_worksheet", "Read Row in Worksheet")

	shared.RegisterSpreadsheetsProps(form, "spreadsheetId", "Spreadsheet", "spreadsheet ID", true)

	shared.RegisterSheetIDProps(form, true)

	shared.RegisterSheetTitleProps(form, true)

	form.TextField("start-row", "start-row").
		Placeholder("Start Row").
		HelpText("The row range of the sheet in the format of A1 notation where to start.").
		Required(true)

	form.TextField("end-row", "end-row").
		Placeholder("End Row").
		HelpText("The row range of the sheet in the format of A1 notation where to end.").
		Required(true)

	form.CheckboxField("includeTeamDrives", "includeTeamDrives").
		Placeholder("Include Team Drives").
		HelpText("Include files from Team Drives in results")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *ReadRowInWorksheetAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *ReadRowInWorksheetAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[readRowInWorksheetActionProps](ctx)
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

	if input.SheetID == "" {
		return nil, errors.New("sheet ID is required")
	}
	if input.SheetTitle == "" {
		return nil, errors.New("sheet title is required")
	}

	readRange := fmt.Sprintf("%s!%s:%s", input.SheetTitle, input.StartRow, input.EndRow)

	spreadsheetCall := sheetService.Spreadsheets.Get(input.SpreadSheetID).
		Ranges(readRange).
		IncludeGridData(true)

	spreadsheet, err := spreadsheetCall.Do()
	if err != nil {
		return nil, err
	}

	var rows [][]interface{}
	for _, sheet := range spreadsheet.Sheets {
		if sheet.Properties.SheetId == shared.ConvertToInt64(input.SheetID) {
			if sheet.Data != nil && len(sheet.Data) > 0 {
				gridData := sheet.Data[0]
				for _, rowData := range gridData.RowData {
					var row []interface{}
					for _, cell := range rowData.Values {
						var value interface{}
						if cell.UserEnteredValue != nil {
							if cell.UserEnteredValue.NumberValue != nil {
								value = *cell.UserEnteredValue.NumberValue
							} else if cell.UserEnteredValue.StringValue != nil {
								value = *cell.UserEnteredValue.StringValue
							} else if cell.UserEnteredValue.BoolValue != nil {
								value = *cell.UserEnteredValue.BoolValue
							} else if cell.UserEnteredValue.FormulaValue != nil {
								value = *cell.UserEnteredValue.FormulaValue
							} else {
								value = nil
							}
						}
						row = append(row, value)
					}
					rows = append(rows, row)
				}
				break
			}
		}
	}

	if rows == nil {
		return nil, errors.New("no rows found in the specified range")
	}

	return core.JSON(rows), nil
}

func NewReadRowInWorksheetAction() sdk.Action {
	return &ReadRowInWorksheetAction{}
}

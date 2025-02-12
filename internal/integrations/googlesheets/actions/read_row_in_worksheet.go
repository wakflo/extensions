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

func (a *ReadRowInWorksheetAction) Name() string {
	return "Read Row in Worksheet"
}

func (a *ReadRowInWorksheetAction) Description() string {
	return "Reads a single row from a worksheet and returns its values as an object. This action is useful when you need to retrieve specific data from a worksheet or perform actions based on the contents of a particular row."
}

func (a *ReadRowInWorksheetAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *ReadRowInWorksheetAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &readRowInWorksheetDocs,
	}
}

func (a *ReadRowInWorksheetAction) Icon() *string {
	return nil
}

func (a *ReadRowInWorksheetAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"spreadsheetId": shared.GetSpreadsheetsInput("Spreadsheet", "spreadsheet ID", true),
		"sheetId":       shared.GetSheetIDInput("Select Sheet to get the ID", "select sheet", true),
		"sheetTitle":    shared.GetSheetTitleInput("Sheet ", "select sheet to read from", true),
		"start-row": autoform.NewShortTextField().
			SetDisplayName("Start Row").
			SetDescription("The row range of the sheet in the format of A1 notation where to start.").
			SetRequired(true).
			Build(),
		"end-row": autoform.NewShortTextField().
			SetDisplayName("End Row").
			SetDescription("The row range of the sheet in the format of A1 notation where to end.").
			SetRequired(true).
			Build(),
		"includeTeamDrives": shared.IncludeTeamFieldInput,
	}
}

func (a *ReadRowInWorksheetAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[readRowInWorksheetActionProps](ctx.BaseContext)
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

	return sdk.JSON(rows), nil
}

func (a *ReadRowInWorksheetAction) Auth() *sdk.Auth {
	return nil
}

func (a *ReadRowInWorksheetAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *ReadRowInWorksheetAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewReadRowInWorksheetAction() sdk.Action {
	return &ReadRowInWorksheetAction{}
}

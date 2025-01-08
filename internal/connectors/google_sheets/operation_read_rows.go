// Copyright 2022-present Wakflo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package googlesheets

import (
	"context"
	"errors"
	"fmt"

	"github.com/wakflo/go-sdk/autoform"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"

	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type readRowsOperationProps struct {
	SpreadSheetID     string `json:"spreadsheetId,omitempty"`
	SheetID           string `json:"sheetId"`
	SheetTitle        string `json:"sheetTitle"`
	Range             string `json:"range"`
	StartRow          string `json:"start-row"`
	EndRow            string `json:"end-row"`
	IncludeTeamDrives bool   `json:"includeTeamDrives"`
}

type ReadRowsOperation struct {
	options *sdk.OperationInfo
}

func NewReadRowsOperation() *ReadRowsOperation {
	return &ReadRowsOperation{
		options: &sdk.OperationInfo{
			Name:        "Read rows",
			Description: "Read rows from a Worksheet",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"spreadsheetId": getSpreadsheetsInput("Spreadsheet", "spreadsheet ID", true),
				"sheetId":       getSheetIDInput("Select Sheet to get the ID", "select sheet", true),
				"sheetTitle":    getSheetTitleInput("Sheet ", "select sheet to read from", true),
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
				"includeTeamDrives": includeTeamFieldInput,
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *ReadRowsOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Token == nil {
		return nil, errors.New("missing google auth token")
	}

	input := sdk.InputToType[readRowsOperationProps](ctx)
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
		if sheet.Properties.SheetId == convertToInt64(input.SheetID) {
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

func (c *ReadRowsOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *ReadRowsOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}

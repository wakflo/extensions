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

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type addColumnInWorkSheetOperationProps struct {
	SpreadSheetID    string `json:"spreadsheetId"`
	SheetID          string `json:"sheetId"`
	SheetColumnIndex string `json:"sheetColumnIndex"`
}

type AddColumnInWorkSheetOperation struct {
	options *sdk.OperationInfo
}

func NewAddColumnInWorkSheetOperation() *AddColumnInWorkSheetOperation {
	return &AddColumnInWorkSheetOperation{
		options: &sdk.OperationInfo{
			Name:        "Add Column in Sheet",
			Description: "Add a column in a sheet",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"spreadsheetId": getSpreadsheetsInput("Spreadsheet", "spreadsheet ID", true),
				"sheetId":       getSheetIDInput("Sheet", "select sheet", true),
				"sheetColumnIndex": autoform.NewShortTextField().
					SetDisplayName("Sheet Column Index").
					SetDescription("The index of the column where the new column should be added. Index is zero based.").
					SetRequired(true).
					Build(),
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *AddColumnInWorkSheetOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Token == nil {
		return nil, errors.New("missing google auth token")
	}

	input := sdk.InputToType[addColumnInWorkSheetOperationProps](ctx)
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

	if input.SheetColumnIndex == "" {
		return nil, errors.New("sheet column index is required")
	}

	// Create the InsertDimensionRequest
	insertDimensionRequest := &sheets.InsertDimensionRequest{
		Range: &sheets.DimensionRange{
			SheetId:    convertToInt64(input.SheetID),
			Dimension:  "COLUMNS",
			StartIndex: convertToInt64(input.SheetColumnIndex),
			EndIndex:   convertToInt64(input.SheetColumnIndex) + 1,
		},
		InheritFromBefore: false,
	}

	// Create the BatchUpdateSpreadsheetRequest
	batchUpdateRequest := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				InsertDimension: insertDimensionRequest,
			},
		},
	}

	spreadsheet, err := sheetService.Spreadsheets.BatchUpdate(input.SpreadSheetID, batchUpdateRequest).
		Do()
	return spreadsheet, err
}

func (c *AddColumnInWorkSheetOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *AddColumnInWorkSheetOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}

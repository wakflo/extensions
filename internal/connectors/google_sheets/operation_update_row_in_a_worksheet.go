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

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type updateRowInWorkSheetOperationProps struct {
	SpreadSheetID string          `json:"spreadSheetId"`
	SheetTitle    string          `json:"sheetTitle"`
	SheetRow      string          `json:"sheetRow"`
	Values        [][]interface{} `json:"values"`
}

type UpdateRowInWorkSheetOperation struct {
	options *sdk.OperationInfo
}

func NewUpdateRowInWorkSheetOperation() *UpdateRowInWorkSheetOperation {
	return &UpdateRowInWorkSheetOperation{
		options: &sdk.OperationInfo{
			Name:        "Update Row In Sheet",
			Description: "Update a row in a sheet",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"spreadsheetId": getSpreadsheetsInput("Spreadsheet", "spreadsheet ID", true),
				"sheetTitle":    getSheetTitleInput("Sheet", "select sheet", true),
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
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *UpdateRowInWorkSheetOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Token == nil {
		return nil, errors.New("missing google auth token")
	}

	input, errs := sdk.InputToTypeSafely[updateRowInWorkSheetOperationProps](ctx)
	if errs != nil {
		return nil, errs
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

func (c *UpdateRowInWorkSheetOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *UpdateRowInWorkSheetOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}

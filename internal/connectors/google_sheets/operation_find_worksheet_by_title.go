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

	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type findWorkSheetByTitleOperationProps struct {
	SpreadSheetID string `json:"spreadsheetId,omitempty"`
	SheetTitle    string `json:"sheetTitle"`
}

type FindWorkSheetByTitleOperation struct {
	options *sdk.OperationInfo
}

func NewFindWorkSheetByTitleOperation() *FindWorkSheetByTitleOperation {
	return &FindWorkSheetByTitleOperation{
		options: &sdk.OperationInfo{
			Name:        "Find Worksheet By Title",
			Description: "Find a worksheet by title",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"spreadsheetId": getSpreadsheetsInput("Spreadsheet", "spreadsheet ID", true),
				"sheetTitle":    getSheetTitleInput("Sheet", "select sheet", true),
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *FindWorkSheetByTitleOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Token == nil {
		return nil, errors.New("missing google auth token")
	}

	input := sdk.InputToType[findWorkSheetByTitleOperationProps](ctx)
	sheetService, err := sheets.NewService(context.Background(), option.WithTokenSource(*ctx.Auth.TokenSource))
	if err != nil {
		return nil, err
	}

	if input.SpreadSheetID == "" {
		return nil, errors.New("spreadsheet ID is required")
	}

	if input.SheetTitle == "" {
		return nil, errors.New("sheet title is required")
	}

	spreadsheet, err := sheetService.Spreadsheets.Get(input.SpreadSheetID).
		Do()

	for _, sheet := range spreadsheet.Sheets {
		if sheet.Properties.Title == input.SheetTitle {
			return sheet, nil
		}

		if sheet.Properties.Title != input.SheetTitle {
			const notFound string = "not sheet found"
			return notFound, nil
		}
	}
	return spreadsheet, err
}

func (c *FindWorkSheetByTitleOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *FindWorkSheetByTitleOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}

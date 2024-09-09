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
	"encoding/json"
	"errors"
	"io"
	"strconv"

	"github.com/gookit/goutil/arrutil"
	fastshot "github.com/opus-domini/fast-shot"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

var (
	// #nosec
	tokenURL   = "https://oauth2.googleapis.com/token"
	sharedAuth = autoform.NewOAuthField("https://accounts.google.com/o/oauth2/auth", &tokenURL, []string{
		"https://www.googleapis.com/auth/spreadsheets https://www.googleapis.com/auth/drive.readonly",
	}).Build()
)

func getSpreadsheetsInput(title string, desc string, required bool) *sdkcore.AutoFormSchema {
	getSpreadsheetFiles := func(ctx *sdkcore.DynamicFieldContext) (interface{}, error) {
		input := sdk.DynamicInputToType[struct {
			IncludeTeamDrives bool `json:"includeTeamDrives"`
		}](ctx)

		client := fastshot.NewClient("https://www.googleapis.com/drive/v3").
			Auth().BearerToken(ctx.Auth.AccessToken).
			Header().
			AddAccept("application/json").
			Build()

		queryParams := map[string]string{
			"q":                 "mimeType='application/vnd.google-apps.spreadsheet' and trashed = false",
			"supportsAllDrives": "true",
		}

		if input.IncludeTeamDrives {
			queryParams["includeItemsFromAllDrives"] = "true"
		}

		rsp, err := client.GET("/files").Query().
			AddParams(queryParams).Send()
		if err != nil {
			return nil, err
		}

		defer rsp.Body().Close()

		if rsp.Status().IsError() {
			return nil, errors.New(rsp.Status().Text())
		}

		defer rsp.Body().Close()
		byts, err := io.ReadAll(rsp.Body().Raw())
		if err != nil {
			return nil, err
		}

		var body FileList
		err = json.Unmarshal(byts, &body)
		if err != nil {
			return nil, err
		}

		return body.Files, nil
	}

	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName(title).
		SetDescription(desc).
		SetDynamicOptions(&getSpreadsheetFiles).
		SetRequired(required).Build()
}

func getSheetIDInput(title string, desc string, required bool) *sdkcore.AutoFormSchema {
	getSheetID := func(ctx *sdkcore.DynamicFieldContext) (interface{}, error) {
		input := sdk.DynamicInputToType[struct {
			SpreadSheetID string `json:"spreadsheetId,omitempty"`
		}](ctx)

		client := fastshot.NewClient("https://sheets.googleapis.com/v4/spreadsheets").
			Auth().BearerToken(ctx.Auth.AccessToken).
			Header().
			AddAccept("application/json").
			Build()

		rsp, err := client.GET("/" + input.SpreadSheetID).Send()
		if err != nil {
			return nil, err
		}

		defer rsp.Body().Close()

		if rsp.Status().IsError() {
			return nil, errors.New(rsp.Status().Text())
		}

		defer rsp.Body().Close()
		byts, err := io.ReadAll(rsp.Body().Raw())
		if err != nil {
			return nil, err
		}

		var body Spreadsheet

		err = json.Unmarshal(byts, &body)
		if err != nil {
			return nil, err
		}

		sheet := body.Sheets
		return arrutil.Map[Sheet, map[string]any](sheet, func(input Sheet) (target map[string]any, find bool) {
			return map[string]any{
				"id":   input.Properties.SheetID,
				"name": input.Properties.Title,
			}, true
		}), nil
	}

	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName(title).
		SetDescription(desc).
		SetDynamicOptions(&getSheetID).
		SetRequired(required).Build()
}

func getSheetTitleInput(title string, desc string, required bool) *sdkcore.AutoFormSchema {
	getSheetTitle := func(ctx *sdkcore.DynamicFieldContext) (interface{}, error) {
		input := sdk.DynamicInputToType[struct {
			SpreadSheetID string `json:"spreadsheetId,omitempty"`
		}](ctx)

		client := fastshot.NewClient("https://sheets.googleapis.com/v4/spreadsheets").
			Auth().BearerToken(ctx.Auth.AccessToken).
			Header().
			AddAccept("application/json").
			Build()

		rsp, err := client.GET("/" + input.SpreadSheetID).Send()
		if err != nil {
			return nil, err
		}

		defer rsp.Body().Close()

		if rsp.Status().IsError() {
			return nil, errors.New(rsp.Status().Text())
		}

		defer rsp.Body().Close()
		byts, err := io.ReadAll(rsp.Body().Raw())
		if err != nil {
			return nil, err
		}

		var body Spreadsheet

		err = json.Unmarshal(byts, &body)
		if err != nil {
			return nil, err
		}

		sheet := body.Sheets
		return arrutil.Map[Sheet, map[string]any](sheet, func(input Sheet) (target map[string]any, find bool) {
			return map[string]any{
				"id":   input.Properties.Title,
				"name": input.Properties.Title,
			}, true
		}), nil
	}

	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName(title).
		SetDescription(desc).
		SetDynamicOptions(&getSheetTitle).
		SetRequired(required).Build()
}

var includeTeamFieldInput = autoform.NewBooleanField().
	SetDisplayName("Include Team Drives Sheets").
	SetDescription("Determines if sheets from Team Drives sheets should be included in the results.").
	SetDefaultValue(false).
	Build()

func convertToInt64(s string) int64 {
	convertedString, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}

	return convertedString
}

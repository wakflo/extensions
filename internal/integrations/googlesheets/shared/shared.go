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

package shared

import (
	"encoding/json"
	"errors"
	"io"
	"strconv"

	"github.com/gookit/goutil/arrutil"
	"github.com/juicycleff/smartform/v1"
	fastshot "github.com/opus-domini/fast-shot"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"

	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

var (
	googleSheetsForm = smartform.NewAuthForm("google-sheets-auth", "Google Sheets OAuth", smartform.AuthStrategyOAuth2)
	_                = googleSheetsForm.
				OAuthField("oauth", "Google Sheets OAuth").
				AuthorizationURL("https://accounts.google.com/o/oauth2/auth").
				TokenURL("https://oauth2.googleapis.com/token").
				Scopes([]string{
			"https://www.googleapis.com/auth/spreadsheets https://www.googleapis.com/auth/drive.readonly",
		}).
		Build()
)

var SharedGoogleSheetsAuth = googleSheetsForm.Build()

func RegisterSpreadsheetsProps(form *smartform.FormBuilder, value string, title string, description string, required bool) *smartform.FieldBuilder {
	getSpreadsheetFiles := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}

		input := sdk.DynamicInputToType[struct {
			IncludeTeamDrives bool `json:"includeTeamDrives"`
		}](ctx)

		client := fastshot.NewClient("https://www.googleapis.com/drive/v3").
			Auth().BearerToken(authCtx.AccessToken).
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

		byts, err := io.ReadAll(rsp.Body().Raw())
		if err != nil {
			return nil, err
		}

		var body FileList
		err = json.Unmarshal(byts, &body)
		if err != nil {
			return nil, err
		}

		return ctx.Respond(body.Files, len(body.Files))
	}

	return form.SelectField(value, title).
		Placeholder("Select a spreadsheet").
		Required(required).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getSpreadsheetFiles)).
				WithFieldReference("includeTeamDrives", "includeTeamDrives").
				WithSearchSupport().
				End().
				RefreshOn("includeTeamDrives").
				GetDynamicSource(),
		).
		HelpText(description)
}

func RegisterSheetIDProps(form *smartform.FormBuilder, required bool) *smartform.FieldBuilder {
	getSheetID := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}

		input := sdk.DynamicInputToType[struct {
			SpreadSheetID string `json:"spreadsheetId,omitempty"`
		}](ctx)

		client := fastshot.NewClient("https://sheets.googleapis.com/v4/spreadsheets").
			Auth().BearerToken(authCtx.AccessToken).
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
		items := arrutil.Map[Sheet, map[string]any](sheet, func(input Sheet) (target map[string]any, find bool) {
			return map[string]any{
				"id":   input.Properties.SheetID,
				"name": input.Properties.Title,
			}, true
		})
		return ctx.Respond(items, len(items))
	}

	return form.SelectField("sheetId", "Sheet ID").
		Placeholder("Select a sheet").
		Required(required).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getSheetID)).
				WithFieldReference("spreadsheetId", "spreadsheetId").
				WithSearchSupport().
				End().
				RefreshOn("spreadsheetId").
				GetDynamicSource(),
		).
		HelpText("Select a sheet ID")
}

func RegisterSheetTitleProps(form *smartform.FormBuilder, required bool) *smartform.FieldBuilder {
	getSheetTitle := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}

		input := sdk.DynamicInputToType[struct {
			SpreadSheetID string `json:"spreadsheetId,omitempty"`
		}](ctx)

		client := fastshot.NewClient("https://sheets.googleapis.com/v4/spreadsheets").
			Auth().BearerToken(authCtx.AccessToken).
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
		items := arrutil.Map[Sheet, map[string]any](sheet, func(input Sheet) (target map[string]any, find bool) {
			return map[string]any{
				"id":   input.Properties.Title,
				"name": input.Properties.Title,
			}, true
		})
		return ctx.Respond(items, len(items))
	}

	return form.SelectField("sheetTitle", "Sheet Title").
		Placeholder("Select a sheet").
		Required(required).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getSheetTitle)).
				WithFieldReference("spreadsheetId", "spreadsheetId").
				WithSearchSupport().
				End().
				RefreshOn("spreadsheetId").
				GetDynamicSource(),
		).
		HelpText("Select a sheet title")
}

func ConvertToInt64(s string) int64 {
	convertedString, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}

	return convertedString
}

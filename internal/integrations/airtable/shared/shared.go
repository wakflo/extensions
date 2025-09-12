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
	"fmt"
	"io"
	"net/http"

	"github.com/juicycleff/smartform/v1"
	fastshot "github.com/opus-domini/fast-shot"
	"github.com/wakflo/go-sdk/v2"

	sdkcontext "github.com/wakflo/go-sdk/v2/context"

	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

const desc = `
    To obtain your personal token, follow these steps:

    1. Log in to your Airtable account.
    2. Visit https://airtable.com/create/tokens/ to create one
    3. Click on "+ Add a base" and select the base you want to use or all bases.
    4. Click on "+ Add a scope" and select "data.records.read", "data.records.write" and "schema.bases.read".
    5. Click on "Create token" and copy the token.
    `

var (
	form = smartform.NewAuthForm("airtable-auth", "Airtable API Authentication", smartform.AuthStrategyCustom)

	_ = form.TextField("api-key", "Personal Access Token (Required)").
		Required(true).
		HelpText(desc)

	AirtableSharedAuth = form.Build()
)

var BaseAPI = "https://api.airtable.com"

func RegisterBasesProps(form *smartform.FormBuilder) *smartform.FieldBuilder {
	getBases := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}

		client := fastshot.NewClient(BaseAPI).
			Auth().BearerToken(authCtx.Extra["api-key"]).
			Header().
			AddAccept("application/json").
			Build()

		rsp, err := client.GET("/v0/meta/bases").Send()
		if err != nil {
			return nil, err
		}

		if rsp.Status().IsError() {
			return nil, errors.New(rsp.Status().Text())
		}

		bytes, err := io.ReadAll(rsp.Raw().Body) //nolint:bodyclose
		if err != nil {
			return nil, err
		}

		var baseRsp Response
		err = json.Unmarshal(bytes, &baseRsp)
		if err != nil {
			return nil, err
		}

		return ctx.Respond(baseRsp.Bases, len(baseRsp.Bases))
	}

	return form.SelectField("bases", "Bases").
		Placeholder("Select a base").
		Required(true).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getBases)).
				WithSearchSupport().
				End().
				GetDynamicSource(),
		).
		HelpText("Bases")
}

func RegisterTablesProps(form *smartform.FormBuilder) *smartform.FieldBuilder {
	getTables := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}

		input := sdk.DynamicInputToType[struct {
			Bases string `json:"bases"`
		}](ctx)

		client := fastshot.NewClient(BaseAPI).
			Auth().BearerToken(authCtx.Extra["api-key"]).
			Header().
			AddAccept("application/json").
			Build()

		fullURL := fmt.Sprintf("/v0/meta/bases/%s/tables", input.Bases)

		rsp, err := client.GET(fullURL).Send()
		if err != nil {
			return nil, err
		}

		if rsp.Status().IsError() {
			return nil, errors.New(rsp.Status().Text())
		}

		bytes, err := io.ReadAll(rsp.Raw().Body) //nolint:bodyclose
		if err != nil {
			return nil, err
		}

		var bases TableResponse
		err = json.Unmarshal(bytes, &bases)
		if err != nil {
			return nil, err
		}

		return ctx.Respond(bases.Tables, len(bases.Tables))
	}

	return form.SelectField("table", "Tables").
		Placeholder("Select a table").
		Required(true).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getTables)).
				WithFieldReference("bases", "bases").
				WithSearchSupport().
				End().
				RefreshOn("bases").
				GetDynamicSource(),
		).
		HelpText("Tables")
}

func RegisterFieldsProps(form *smartform.FormBuilder) *smartform.FieldBuilder {
	getFields := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}

		input := sdk.DynamicInputToType[struct {
			Bases string `json:"bases"`
			Table string `json:"table"`
		}](ctx)

		client := fastshot.NewClient(BaseAPI).
			Auth().BearerToken(authCtx.Extra["api-key"]).
			Header().
			AddAccept("application/json").
			Build()

		fullURL := fmt.Sprintf("/v0/meta/bases/%s/tables", input.Bases)

		rsp, err := client.GET(fullURL).Send()
		if err != nil {
			return nil, err
		}

		if rsp.Status().IsError() {
			return nil, errors.New(rsp.Status().Text())
		}

		bytes, err := io.ReadAll(rsp.Raw().Body) //nolint:bodyclose
		if err != nil {
			return nil, err
		}

		var tables TableResponse
		err = json.Unmarshal(bytes, &tables)
		if err != nil {
			return nil, err
		}

		var selectedTable Table
		for _, table := range tables.Tables {
			if table.ID == input.Table {
				selectedTable = table
				break
			}
		}

		return ctx.Respond(selectedTable.Fields, len(selectedTable.Fields))
	}

	return form.SelectField("fields", "Search Fields").
		Placeholder("Select fields").
		Required(true).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getFields)).
				WithFieldReference("bases", "bases").
				WithFieldReference("table", "table").
				WithSearchSupport().
				End().
				RefreshOn("table").
				GetDynamicSource(),
		).
		HelpText("Fields for selected Table")
}

func RegisterViewsProps(form *smartform.FormBuilder) *smartform.FieldBuilder {
	getViews := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}

		input := sdk.DynamicInputToType[struct {
			Bases string `json:"bases"`
			Table string `json:"table"`
		}](ctx)

		client := fastshot.NewClient(BaseAPI).
			Auth().BearerToken(authCtx.Extra["api-key"]).
			Header().
			AddAccept("application/json").
			Build()

		fullURL := fmt.Sprintf("/v0/meta/bases/%s/tables", input.Bases)

		rsp, err := client.GET(fullURL).Send()
		if err != nil {
			return nil, err
		}

		if rsp.Status().IsError() {
			return nil, errors.New(rsp.Status().Text())
		}

		bytes, err := io.ReadAll(rsp.Raw().Body) //nolint:bodyclose
		if err != nil {
			return nil, err
		}

		var tables TableResponse
		err = json.Unmarshal(bytes, &tables)
		if err != nil {
			return nil, err
		}

		var selectedTable Table
		for _, table := range tables.Tables {
			if table.ID == input.Table {
				selectedTable = table
				break
			}
		}

		return ctx.Respond(selectedTable.Views, len(selectedTable.Views))
	}

	return form.SelectField("view", "View").
		Placeholder("Select a view").
		Required(true).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getViews)).
				WithFieldReference("bases", "bases").
				WithFieldReference("table", "table").
				WithSearchSupport().
				End().
				RefreshOn("table").
				GetDynamicSource(),
		).
		HelpText("View for selected Table")
}

func RegisterRecordsProps(form *smartform.FormBuilder) *smartform.FieldBuilder {
	getRecords := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}

		input := sdk.DynamicInputToType[struct {
			BasesID string `json:"bases"`
			TableID string `json:"table"`
		}](ctx)

		// Build the full URL
		fullURL := fmt.Sprintf("%s/v0/%s/%s", BaseAPI, input.BasesID, input.TableID)

		// Use the existing AirtableRequest helper function
		response, err := AirtableRequest(authCtx.Extra["api-key"], fullURL, "GET")
		if err != nil {
			return nil, err
		}

		// Convert response to our expected structure
		responseMap, ok := response.(map[string]interface{})
		if !ok {
			return nil, errors.New("unexpected response format")
		}

		// Extract records array
		recordsInterface, exists := responseMap["records"]
		if !exists {
			return nil, errors.New("no records field in response")
		}

		recordsArray, ok := recordsInterface.([]interface{})
		if !ok {
			return nil, errors.New("records field is not an array")
		}

		// Transform records into selectable options
		options := make([]interface{}, 0, len(recordsArray))

		for _, recordInterface := range recordsArray {
			record, ok := recordInterface.(map[string]interface{})
			if !ok {
				continue
			}

			// Extract record ID
			recordID, ok := record["id"].(string)
			if !ok {
				continue
			}

			// Extract fields
			fields, ok := record["fields"].(map[string]interface{})
			if !ok {
				fields = make(map[string]interface{})
			}

			// Get the task name for display
			displayName := fmt.Sprintf("Record %s", recordID) // Default fallback

			if taskName, ok := fields["Task Name"].(string); ok && taskName != "" {
				displayName = taskName
			}

			option := map[string]interface{}{
				"id":   recordID,
				"name": displayName,
			}

			options = append(options, option)
		}

		return ctx.Respond(options, len(options))
	}

	return form.SelectField("record-id", "Record").
		Placeholder("Select a record").
		Required(true).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getRecords)).
				WithFieldReference("bases", "bases").
				WithFieldReference("table", "table").
				WithSearchSupport().
				End().
				RefreshOn("bases", "table").
				GetDynamicSource(),
		).
		HelpText("Select a record from the table")
}

func AirtableRequest(accessToken, reqURL, requestType string) (interface{}, error) {
	req, err := http.NewRequest(requestType, reqURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+accessToken)
	client := &http.Client{}
	res, errs := client.Do(req)
	if errs != nil {
		return nil, errs
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var response interface{}
	if newErrs := json.Unmarshal(body, &response); newErrs != nil {
		return nil, errors.New("error parsing response")
	}

	return response, nil
}

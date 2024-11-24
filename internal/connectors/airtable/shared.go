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

package airtable

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	fastshot "github.com/opus-domini/fast-shot"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

var sharedAuth = autoform.NewCustomAuthField().
	SetFields(map[string]*sdkcore.AutoFormSchema{
		"api-key": autoform.NewLongTextField().
			SetDisplayName("Personal Access Token").
			SetDescription("Personal Access Token").
			SetRequired(true).
			Build(),
	}).
	Build()

var baseAPI = "https://api.airtable.com"

func getBasesInput() *sdkcore.AutoFormSchema {
	getBases := func(ctx *sdkcore.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		client := fastshot.NewClient(baseAPI).
			Auth().BearerToken(ctx.Auth.Extra["api-key"]).
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

	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName("Bases").
		SetDescription("Bases").
		SetDependsOn([]string{"connection"}).
		SetDynamicOptions(&getBases).
		SetRequired(true).Build()
}

func getTablesInput() *sdkcore.AutoFormSchema {
	getTables := func(ctx *sdkcore.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		input := sdk.DynamicInputToType[struct {
			Bases string `json:"bases"`
		}](ctx)
		client := fastshot.NewClient(baseAPI).
			Auth().BearerToken(ctx.Auth.Extra["api-key"]).
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

	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName("Tables").
		SetDescription("Tables").
		SetDependsOn([]string{"connection"}).
		SetDynamicOptions(&getTables).
		SetRequired(true).Build()
}

func getFieldsInput() *sdkcore.AutoFormSchema {
	getFields := func(ctx *sdkcore.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		input := sdk.DynamicInputToType[struct {
			Bases string `json:"bases"`
			Table string `json:"table"`
		}](ctx)
		client := fastshot.NewClient(baseAPI).
			Auth().BearerToken(ctx.Auth.Extra["api-key"]).
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

	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName("Search Fields").
		SetDescription("Fields for selected Table").
		SetDynamicOptions(&getFields).
		SetRequired(true).Build()
}

func getViewsInput() *sdkcore.AutoFormSchema {
	getViews := func(ctx *sdkcore.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		input := sdk.DynamicInputToType[struct {
			Bases string `json:"bases"`
			Table string `json:"table"`
		}](ctx)
		client := fastshot.NewClient(baseAPI).
			Auth().BearerToken(ctx.Auth.Extra["api-key"]).
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

	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName("View").
		SetDescription("view for selected Table").
		SetDynamicOptions(&getViews).
		SetRequired(true).Build()
}

func airtableRequest(accessToken, reqURL, requestType string) (interface{}, error) {
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

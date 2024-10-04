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

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type updateRecordOperationProps struct {
	Name  string `json:"name"`
	Bases string `json:"bases"`
	Table string `json:"table"`
}

type UpdateRecordOperation struct {
	options *sdk.OperationInfo
}

func NewUpdateRecordOperation() *UpdateRecordOperation {
	return &UpdateRecordOperation{
		options: &sdk.OperationInfo{
			Name:        "Update Record",
			Description: "Update a record",
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"record-id": autoform.NewShortTextField().
					SetDisplayName("Record ID").
					SetDescription("The record's ID").
					SetRequired(true).Build(),
				"bases":  getBasesInput(),
				"table":  getTablesInput(),
				"fields": getFieldsInput(),
				"views":  getViewsInput(),
			},
			SampleOutput: map[string]interface{}{},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
			RequireAuth: true,
		},
	}
}

func (c *UpdateRecordOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Extra["api-key"] == "" {
		return nil, errors.New("missing airtable access token")
	}
	accessToken := ctx.Auth.Extra["api-key"]
	input := sdk.InputToType[updateRecordOperationProps](ctx)

	// data := map[string]interface{}{
	//	"records": map[string]interface{}{
	//		"fields": map[string]interface{}{
	//			"name":      input.Name,
	//			"workspace": input.Bases,
	//		},
	//	},
	// }
	//
	// taskJSON, err := json.Marshal(data)
	// if err != nil {
	//	return nil, err
	// }

	reqURL := fmt.Sprintf("https://api.airtable.com/v0/meta/bases/%s/tables", input.Bases)

	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+accessToken)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var response interface{}
	if errs := json.Unmarshal(body, &response); errs != nil {
		return nil, errors.New("error parsing response")
	}

	return response, nil
}

func (c *UpdateRecordOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *UpdateRecordOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}

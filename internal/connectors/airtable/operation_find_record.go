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
	"errors"
	"fmt"
	"net/http"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type findRecordOperationProps struct {
	Bases        string `json:"bases"`
	Table        string `json:"table"`
	RecordID     string `json:"record-id"`
	SearchFields string `json:"search-fields,omitempty"`
	Views        string `json:"views"`
}

type FindRecordOperation struct {
	options *sdk.OperationInfo
}

func NewFindRecordOperation() *FindRecordOperation {
	return &FindRecordOperation{
		options: &sdk.OperationInfo{
			Name:        "Find Airtable Record",
			Description: "Find a record in airtable.",
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"record-id": autoform.NewShortTextField().
					SetDisplayName("Record ID").
					SetDescription("The ID of the record").
					SetRequired(true).Build(),
				"bases": getBasesInput(),
				"table": getTablesInput(),
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

func (c *FindRecordOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Extra["api-key"] == "" {
		return nil, errors.New("missing airtable access token")
	}
	apiKey := ctx.Auth.Extra["api-key"]
	input := sdk.InputToType[findRecordOperationProps](ctx)

	reqURL := fmt.Sprintf("https://api.airtable.com/v0/%s/%s/%s", input.Bases, input.Table, input.RecordID)

	response, err := airtableRequest(apiKey, reqURL, http.MethodGet)
	if err != nil {
		return nil, errors.New("error fetching data")
	}

	return response, nil
}

func (c *FindRecordOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *FindRecordOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}

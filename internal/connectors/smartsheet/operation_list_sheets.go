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

package smartsheet

import (
	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type listInvoicesOperationProps struct {
	Name string `json:"name"`
}

type ListSheetsOperation struct {
	options *sdk.OperationInfo
}

func NewListSheetsOperation() *ListSheetsOperation {
	return &ListSheetsOperation{
		options: &sdk.OperationInfo{
			Name:        "List Sheets",
			Description: "list all sheets",
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"name": autoform.NewShortTextField().
					SetDisplayName("").
					SetDescription("").
					SetRequired(false).Build(),
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

func (c *ListSheetsOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	_ = sdk.InputToType[listInvoicesOperationProps](ctx)

	url := "/2.0/sheets"

	sheets, err := getSmartsheetClient(ctx.Auth.AccessToken, url)
	if err != nil {
		return nil, err
	}
	return sheets, nil
}

func (c *ListSheetsOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *ListSheetsOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}

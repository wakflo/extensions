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

package harvest

import (
	"errors"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type listInvoicesOperationProps struct {
	Name string `json:"name"`
}

type ListInvoicesOperation struct {
	options *sdk.OperationInfo
}

func NewListInvoicesOperation() *ListInvoicesOperation {
	return &ListInvoicesOperation{
		options: &sdk.OperationInfo{
			Name:        "List Invoices",
			Description: "list all invoices",
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

func (c *ListInvoicesOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing Harvest auth token")
	}

	_ = sdk.InputToType[listInvoicesOperationProps](ctx)

	url := "/v2/invoices"

	invoices, err := getHarvestClient(ctx.Auth.AccessToken, url)
	if err != nil {
		return nil, err
	}
	invoiceArray, ok := invoices["invoices"].(interface{})
	if !ok {
		return nil, errors.New("failed to extract issues from response")
	}
	return invoiceArray, nil
}

func (c *ListInvoicesOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *ListInvoicesOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}

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

package stripe

import (
	"errors"
	"net/http"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type retrieveCustomerOperationProps struct {
	CustomerID string `json:"customer"`
}

type RetrieveCustomerOperation struct {
	options *sdk.OperationInfo
}

func NewRetrieveCustomerOperation() *RetrieveCustomerOperation {
	return &RetrieveCustomerOperation{
		options: &sdk.OperationInfo{
			Name:        "Retrieve Customer By ID",
			Description: "Retrieve a customer by ID",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"customer": autoform.NewShortTextField().
					SetDisplayName("Customer ID").
					SetDescription("undefined").
					SetRequired(true).
					Build(),
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *RetrieveCustomerOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Extra["api-key"] == "" {
		return nil, errors.New("missing stripe secret api-key")
	}
	apiKey := ctx.Auth.Extra["api-key"]

	input := sdk.InputToType[retrieveCustomerOperationProps](ctx)

	reqURL := "/v1/customers/" + input.CustomerID

	resp, err := stripClient(apiKey, reqURL, http.MethodGet, nil, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *RetrieveCustomerOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *RetrieveCustomerOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}

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
	"net/url"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type searchCustomerOperationProps struct {
	Email string `json:"email"`
}

type SearchCustomerOperation struct {
	options *sdk.OperationInfo
}

func NewSearchCustomerOperation() *SearchCustomerOperation {
	return &SearchCustomerOperation{
		options: &sdk.OperationInfo{
			Name:        "Search Customer",
			Description: "Search customer",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"email": autoform.NewShortTextField().
					SetDisplayName("Email").
					SetDescription("customer email").
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

func (c *SearchCustomerOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Extra["api-key"] == "" {
		return nil, errors.New("missing stripe secret api-key")
	}
	apiKey := ctx.Auth.Extra["api-key"]

	input := sdk.InputToType[searchCustomerOperationProps](ctx)

	params := url.Values{}
	params.Add("query", "email:'"+input.Email+"'")

	reqURL := "/v1/customers/search"

	resp, err := stripClient(apiKey, reqURL, http.MethodGet, nil, params)
	if err != nil {
		return nil, err
	}

	nodes, ok := resp["data"].([]interface{})
	if !ok {
		return nil, errors.New("failed to extract data from response")
	}

	return nodes, nil
}

func (c *SearchCustomerOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *SearchCustomerOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}

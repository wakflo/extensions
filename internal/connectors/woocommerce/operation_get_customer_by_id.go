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

package woocommerce

import (
	"errors"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type GetCustomerOperation struct {
	options *sdk.OperationInfo
}

type getCustomerOperationProps struct {
	CustomerID int `json:"customer-id"`
}

func NewGetCustomerOperation() *GetCustomerOperation {
	return &GetCustomerOperation{
		options: &sdk.OperationInfo{
			Name:        "Get customer by ID",
			Description: "Get customer by ID",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"customer-id": autoform.NewNumberField().
					SetDisplayName("Customer ID").
					SetDescription("Enter customer ID").
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

func (c *GetCustomerOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	baseURL := ctx.Auth.Extra["shop-url"]
	consumerKey := ctx.Auth.Extra["consumer-key"]
	consumerSecret := ctx.Auth.Extra["consumer-secret"]

	if baseURL == "" || consumerKey == "" || consumerSecret == "" {
		return nil, errors.New("missing WooCommerce authentication credentials")
	}

	input := sdk.InputToType[getCustomerOperationProps](ctx)

	wooClient := initializeWooCommerceClient(baseURL, consumerKey, consumerSecret)

	customer, err := wooClient.Services.Customer.One(input.CustomerID)
	if err != nil {
		return nil, err
	}

	return customer, nil
}

func (c *GetCustomerOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *GetCustomerOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}

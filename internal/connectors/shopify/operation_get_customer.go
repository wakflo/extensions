// Copyright 2022-present Wakflo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package shopify

import (
	"context"
	"errors"
	"fmt"
	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type getCustomerOperationProps struct {
	CustomerID uint64 `json:"customerId"`
}
type GetCustomerOperation struct {
	options *sdk.OperationInfo
}

func NewGetCustomerOperation() *GetCustomerOperation {
	return &GetCustomerOperation{
		options: &sdk.OperationInfo{
			Name:        "Get Customer",
			Description: "Get an existing customer's information.",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"customerId": autoform.NewNumberField().
					SetDisplayName("Customer ID").
					SetDescription("The ID of the customer.").
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
	if ctx.Auth.Extra["token"] == "" {
		return nil, errors.New("missing shopify auth token")
	}
	input := sdk.InputToType[getCustomerOperationProps](ctx)

	domain := ctx.Auth.Extra["domain"]
	shopName := domain + ".myshopify.com"
	client := getShopifyClient(shopName, ctx.Auth.Extra["token"])

	customer, err := client.Customer.Get(context.Background(), input.CustomerID, nil)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, fmt.Errorf("no customer found with ID '%d'", input.CustomerID)
	}

	return sdk.JSON(map[string]interface{}{

		"raw_customer": customer,
	}), nil
}
func (c *GetCustomerOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}
func (c *GetCustomerOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}

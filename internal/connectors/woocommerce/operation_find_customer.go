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
	"fmt"

	"github.com/hiscaler/woocommerce-go"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type FindCustomerOperation struct {
	options *sdk.OperationInfo
}

type findCustomerOperationProps struct {
	Email string `json:"email"`
}

func NewFindCustomerOperation() *FindCustomerOperation {
	return &FindCustomerOperation{
		options: &sdk.OperationInfo{
			Name:        "Find a customer",
			Description: "Find a customer",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"productId": autoform.NewShortTextField().
					SetDisplayName("Customer Email").
					SetDescription("Enter the email address of the customer.").
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

func (c *FindCustomerOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	baseURL := ctx.Auth.Extra["shop-url"]
	consumerKey := ctx.Auth.Extra["consumer-key"]
	consumerSecret := ctx.Auth.Extra["consumer-secret"]

	if baseURL == "" || consumerKey == "" || consumerSecret == "" {
		return nil, errors.New("missing WooCommerce authentication credentials")
	}

	input := sdk.InputToType[findCustomerOperationProps](ctx)

	wooClient := initializeWooCommerceClient(baseURL, consumerKey, consumerSecret)

	params := woocommerce.CustomersQueryParams{
		Email: input.Email,
	}

	customer, total, totalPages, isLastPage, err := wooClient.Services.Customer.All(params)
	if err != nil {
		return nil, err
	}
	fmt.Println(totalPages, total, isLastPage)

	return customer, nil
}

func (c *FindCustomerOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *FindCustomerOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}

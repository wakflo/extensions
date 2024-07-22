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

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type getTransactionsOperationProps struct {
	OrderID uint64 `json:"orderId"`
}
type GetTransactionsOperation struct {
	options *sdk.OperationInfo
}

func NewGetTransactionsOperation() *GetTransactionsOperation {
	return &GetTransactionsOperation{
		options: &sdk.OperationInfo{
			Name:        "Get Order Transactions",
			Description: "Get an order's transactions.",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"orderId": autoform.NewNumberField().
					SetDisplayName("Order ID").
					SetDescription("The ID of the order.").
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

func (c *GetTransactionsOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Extra["token"] == "" {
		return nil, errors.New("missing shopify auth token")
	}

	input := sdk.InputToType[getTransactionsOperationProps](ctx)
	domain := ctx.Auth.Extra["domain"]
	shopName := domain + ".myshopify.com"
	client := getShopifyClient(shopName, ctx.Auth.Extra["token"])

	transactions, err := client.Transaction.List(context.Background(), input.OrderID, nil)
	if err != nil {
		return nil, err
	}

	if transactions == nil {
		return nil, errors.New("no transactions found with ID ")
	}
	return sdk.JSON(map[string]interface{}{
		"transactions_details": transactions,
	}), nil
}

func (c *GetTransactionsOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *GetTransactionsOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}

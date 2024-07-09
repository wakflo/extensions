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

package shopify

import (
	"context"
	"errors"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type getTransactionOperationProps struct {
	OrderID       uint64 `json:"orderId"`
	TransactionID uint64 `json:"transactionId"`
}

type GetTransactionOperation struct {
	options *sdk.OperationInfo
}

func NewGetTransactionOperation() *GetTransactionOperation {
	return &GetTransactionOperation{
		options: &sdk.OperationInfo{
			Name:        "Get  Transaction",
			Description: "Get an existing transaction's information.",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"orderId": autoform.NewNumberField().
					SetDisplayName("Order ID").
					SetDescription("The ID of the order.").
					SetRequired(true).
					Build(),
				"transactionId": autoform.NewNumberField().
					SetDisplayName("Transaction ID").
					SetDescription("The ID of the transaction.").
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

func (c *GetTransactionOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {

	if ctx.Auth.Extra["token"] == "" {
		return nil, errors.New("missing shopify auth token")
	}

	input := sdk.InputToType[getTransactionOperationProps](ctx)

	domain := ctx.Auth.Extra["domain"]
	shopName := domain + ".myshopify.com"
	client := getShopifyClient(shopName, ctx.Auth.Extra["token"])

	transaction, err := client.Transaction.Get(context.Background(), input.OrderID, input.TransactionID, nil)
	if err != nil {
		return nil, err
	}

	if transaction == nil {
		return nil, errors.New("no transaction found with ID ")
	}

	return sdk.JSON(map[string]interface{}{
		"transaction": transaction,
	}), nil

}

func (c *GetTransactionOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *GetTransactionOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}

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

	goshopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/shopspring/decimal"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type createTransactionOperationProps struct {
	OrderID       uint64           `json:"orderId"`
	ParentID      *int64           `json:"parentId"`
	Kind          string           `json:"kind"`
	Currency      string           `json:"currency"`
	Amount        *decimal.Decimal `json:"amount"`
	Authorization string           `json:"authorization"`
	Source        string           `json:"source"`
}
type CreateTransactionOperation struct {
	options *sdk.OperationInfo
}

func NewCreateTransactionOperation() *CreateTransactionOperation {
	return &CreateTransactionOperation{
		options: &sdk.OperationInfo{
			Name:        "Create Transaction",
			Description: "Create a new transaction.",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"orderId": autoform.NewNumberField().
					SetDisplayName("Order ID").
					SetDescription("The ID of the order to create a transaction for.").
					SetRequired(true).
					Build(),
				"kind": autoform.NewSelectField().
					SetDisplayName("Type").
					SetOptions(shopifyTransactionKinds).
					SetRequired(true).
					Build(),
				"parentId": autoform.NewNumberField().
					SetDisplayName("Parent ID").
					SetDescription("The ID of an associated transaction.").
					SetRequired(false).
					Build(),
				"currency": autoform.NewShortTextField().
					SetDisplayName("Currency").
					SetRequired(false).
					Build(),
				"amount": autoform.NewNumberField().
					SetDisplayName("Amount").
					SetRequired(false).
					Build(),
				"authorization": autoform.NewShortTextField().
					SetDisplayName("Authorization Key.").
					SetRequired(false).
					Build(),
				"source": autoform.NewShortTextField().
					SetDisplayName("Source").
					SetDescription("An optional origin of the transaction. Set to external to import a cash transaction for the associated order.").
					SetRequired(false).
					Build(),
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *CreateTransactionOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Extra["token"] == "" {
		return nil, errors.New("missing shopify auth token")
	}
	input := sdk.InputToType[createTransactionOperationProps](ctx)
	domain := ctx.Auth.Extra["domain"]
	shopName := domain + ".myshopify.com"
	client := getShopifyClient(shopName, ctx.Auth.Extra["token"])
	newTransaction := goshopify.Transaction{
		ParentId:      input.ParentID,
		Kind:          input.Kind,
		OrderId:       input.OrderID,
		Amount:        input.Amount,
		Authorization: input.Authorization,
		Currency:      input.Currency,
		Source:        input.Source,
	}
	transaction, err := client.Transaction.Create(context.Background(), input.OrderID, newTransaction)
	if err != nil {
		return nil, err
	}

	if transaction == nil {
		return nil, fmt.Errorf("no order found with ID '%d'", input.OrderID)
	}

	return map[string]interface{}{
		"new_transaction": transaction,
	}, nil
}

func (c *CreateTransactionOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *CreateTransactionOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}

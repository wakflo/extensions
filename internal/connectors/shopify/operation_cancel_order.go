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

type cancelOrderOperationProps struct {
	OrderID uint64 `json:"orderId"`
	Reason  string `json:"reason,omitempty"`
}
type CancelOrderOperation struct {
	options *sdk.OperationInfo
}

func NewCancelOrderOperation() *CancelOrderOperation {
	return &CancelOrderOperation{
		options: &sdk.OperationInfo{
			Name:        "Cancel Order",
			Description: "Cancel an existing order by ID.",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"orderId": autoform.NewNumberField().
					SetDisplayName("Order").
					SetDescription("The ID of the order to cancel.").
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

func (c *CancelOrderOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Extra["token"] == "" {
		return nil, errors.New("missing shopify auth token")
	}
	input := sdk.InputToType[cancelOrderOperationProps](ctx)

	domain := ctx.Auth.Extra["domain"]
	shopName := domain + ".myshopify.com"

	client := getShopifyClient(shopName, ctx.Auth.Extra["token"])

	order, err := client.Order.Cancel(context.Background(), input.OrderID, nil)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, fmt.Errorf("no order found with ID '%d'", input.OrderID)
	}
	return sdk.JSON(map[string]interface{}{
		"Result": "Order Successfully cancelled!",
	}), nil
}

func (c *CancelOrderOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *CancelOrderOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}

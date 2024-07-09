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
	"fmt"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type getOrderOperationProps struct {
	OrderID uint64 `json:"orderId"`
}
type GetOrderOperation struct {
	options *sdk.OperationInfo
}

func NewGetOrderOperation() *GetOrderOperation {
	return &GetOrderOperation{
		options: &sdk.OperationInfo{
			Name:        "Get Order",
			Description: "Get an existing order",
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

func (c *GetOrderOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Extra["token"] == "" {
		return nil, errors.New("missing shopify auth token")
	}
	input := sdk.InputToType[getOrderOperationProps](ctx)

	domain := ctx.Auth.Extra["domain"]
	shopName := domain + ".myshopify.com"

	client := getShopifyClient(shopName, ctx.Auth.Extra["token"])

	order, err := client.Order.Get(context.Background(), input.OrderID, nil)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, fmt.Errorf("no order found with ID '%d'", input.OrderID)
	}
	return sdk.JSON(map[string]interface{}{
		"order details": order,
	}), nil
}

func (c *GetOrderOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *GetOrderOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}

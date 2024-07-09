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

type updateOrderOperationProps struct {
	OrderID uint64 `json:"orderId"`
	Tags    string `json:"tags"`
	Note    string `json:"note"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
}
type UpdateOrderOperation struct {
	options *sdk.OperationInfo
}

func NewUpdateOrderOperation() *UpdateOrderOperation {
	return &UpdateOrderOperation{
		options: &sdk.OperationInfo{
			Name:        "Update orders",
			Description: "Update an existing order",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"orderId": autoform.NewNumberField().
					SetDisplayName("Order ID").
					SetDescription("The ID of the order to update").
					SetRequired(true).
					Build(),
				"note": autoform.NewLongTextField().
					SetDisplayName("Note about the order.").
					SetRequired(false).
					Build(),
				"tags": autoform.NewLongTextField().
					SetDisplayName("Tags").
					SetDescription("A string of comma-separated tags for filtering and search.").
					SetRequired(false).
					Build(),
				"phone": autoform.NewShortTextField().
					SetDisplayName("Phone number.").
					SetRequired(false).
					Build(),
				"email": autoform.NewShortTextField().
					SetDisplayName("Email address.").
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
func (c *UpdateOrderOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Extra["token"] == "" {
		return nil, errors.New("missing shopify auth token")
	}
	input := sdk.InputToType[updateOrderOperationProps](ctx)
	domain := ctx.Auth.Extra["domain"]
	shopName := domain + ".myshopify.com"
	client := getShopifyClient(shopName, ctx.Auth.Extra["token"])
	existingOrder, err := client.Order.Get(context.Background(), input.OrderID, nil)
	if err != nil {
		return nil, err
	}
	if input.Note != "" {
		existingOrder.Note = input.Note
	}
	if input.Tags != "" {
		existingOrder.Tags = input.Tags
	}
	if input.Phone != "" {
		existingOrder.Phone = input.Phone
	}
	if input.Email != "" {
		existingOrder.Email = input.Email
	}
	updatedOrder, err := client.Order.Update(context.Background(), *existingOrder)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"updated_order": updatedOrder,
	}, nil
}
func (c *UpdateOrderOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}
func (c *UpdateOrderOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}

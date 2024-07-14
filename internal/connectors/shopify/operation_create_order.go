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

	goshopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/shopspring/decimal"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type createOrderOperationProps struct {
	ProductID  uint64               `json:"productId"`
	LineItems  []goshopify.LineItem `json:"line_items"`
	VariantID  uint64               `json:"variantId"`
	CustomerID uint64               `json:"customerId"`
	Title      string               `json:"title"`
	Tags       string               `json:"tags"`
	Note       string               `json:"note"`
	Quantity   int                  `json:"quantity"`
	Price      *decimal.Decimal     `json:"price"`
}
type CreateOrderOperation struct {
	options *sdk.OperationInfo
}

func NewCreateOrderOperation() *CreateOrderOperation {
	return &CreateOrderOperation{
		options: &sdk.OperationInfo{
			Name:        "Create orders",
			Description: "Create a new orders",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"productId": autoform.NewNumberField().
					SetDisplayName("Product ID").
					SetDescription("The ID of the product to create the order with.").
					SetRequired(false).
					Build(),
				"variantId": autoform.NewNumberField().
					SetDisplayName("Product Variant").
					SetDescription("The ID of the variant to create the order with.").
					SetRequired(false).
					Build(),
				"customerId": autoform.NewNumberField().
					SetDisplayName("Customer ID").
					SetDescription("The ID of the customer to use.").
					SetRequired(false).
					Build(),
				"title": autoform.NewShortTextField().
					SetDisplayName("Title").
					SetRequired(false).
					Build(),
				"note": autoform.NewLongTextField().
					SetDisplayName("Note about the order.").
					SetRequired(false).
					Build(),
				"tags": autoform.NewLongTextField().
					SetDisplayName("Tags").
					SetDescription("A string of comma-separated tags for filtering and search").
					SetRequired(false).
					Build(),
				"quantity": autoform.NewNumberField().
					SetDisplayName("Quantity").
					SetDescription("The ID of the variant to create the order with.").
					SetRequired(false).
					SetDefaultValue(1).
					Build(),
				"price": autoform.NewShortTextField().
					SetDisplayName("Price").
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

func (c *CreateOrderOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Extra["token"] == "" {
		return nil, errors.New("missing shopify auth token")
	}
	input := sdk.InputToType[createOrderOperationProps](ctx)
	domain := ctx.Auth.Extra["domain"]
	shopName := domain + ".myshopify.com"
	client := getShopifyClient(shopName, ctx.Auth.Extra["token"])
	newOrder := goshopify.Order{
		LineItems: []goshopify.LineItem{
			{
				ProductId: input.ProductID,
				VariantId: input.VariantID,
				Quantity:  input.Quantity,
				Price:     input.Price,
				Title:     input.Title,
			},
		},
		Note: input.Note,
		Tags: input.Tags,
		Customer: &goshopify.Customer{
			Id: input.CustomerID,
		},
	}
	order, err := client.Order.Create(context.Background(), newOrder)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"new_order": order,
	}, nil
}

func (c *CreateOrderOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *CreateOrderOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}

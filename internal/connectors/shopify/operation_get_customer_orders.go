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
	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type getCustomerOrdersOperationProps struct {
	CustomerID uint64 `json:"customerId"`
}
type GetCustomerOrdersOperation struct {
	options *sdk.OperationInfo
}

func NewGetCustomerOrdersOperation() *GetCustomerOrdersOperation {
	return &GetCustomerOrdersOperation{
		options: &sdk.OperationInfo{
			Name:        "Get Customer Orders",
			Description: "Get an existing customer's orders.",
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

func (c *GetCustomerOrdersOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Extra["token"] == "" {
		return nil, errors.New("missing shopify auth token")
	}
	input := sdk.InputToType[getCustomerOrdersOperationProps](ctx)
	domain := ctx.Auth.Extra["domain"]
	shopName := domain + ".myshopify.com"
	client := getShopifyClient(shopName, ctx.Auth.Extra["token"])
	orders, err := client.Order.List(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer orders: %w", err)
	}
	customerOrders := filterOrdersByCustomerID(orders, input.CustomerID)
	simplifiedOrders := make([]map[string]interface{}, len(customerOrders))
	for i, order := range customerOrders {
		simplifiedOrders[i] = map[string]interface{}{
			"ID":                order.Id,
			"Name":              order.Name,
			"Email":             order.Email,
			"CreatedAt":         order.CreatedAt,
			"UpdatedAt":         order.UpdatedAt,
			"TotalPrice":        order.TotalPrice,
			"Currency":          order.Currency,
			"FinancialStatus":   order.FinancialStatus,
			"FulfillmentStatus": order.FulfillmentStatus,
		}
	}
	return map[string]interface{}{
		"customer_orders": simplifiedOrders,
	}, nil
}

func filterOrdersByCustomerID(orders []goshopify.Order, customerID uint64) []goshopify.Order {
	var customerOrders []goshopify.Order
	for _, order := range orders {
		if order.Customer != nil && order.Customer.Id == customerID {
			customerOrders = append(customerOrders, order)
		}
	}
	return customerOrders
}

func (c *GetCustomerOrdersOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *GetCustomerOrdersOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}

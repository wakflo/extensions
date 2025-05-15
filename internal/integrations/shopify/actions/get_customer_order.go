package actions

import (
	"context"
	"fmt"

	goshopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/shopify/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type getCustomerOrderActionProps struct {
	CustomerID uint64 `json:"customerId"`
}

type GetCustomerOrderAction struct{}

func (a *GetCustomerOrderAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_customer_order",
		DisplayName:   "Get Customer Order",
		Description:   "Retrieves customer orders from the specified system or database, allowing you to automate tasks that require access to order information.",
		Type:          core.ActionTypeAction,
		Documentation: getCustomerOrderDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *GetCustomerOrderAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_customer_order", "Get Customer Order")

	form.NumberField("customerId", "Customer ID").
		Placeholder("The ID of the customer.").
		Required(true).
		HelpText("The ID of the customer.")

	schema := form.Build()
	return schema
}

func (a *GetCustomerOrderAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[getCustomerOrderActionProps](ctx)
	if err != nil {
		return nil, err
	}

	client, err := shared.CreateClient(ctx)
	if err != nil {
		return nil, err
	}

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

func (a *GetCustomerOrderAction) Auth() *core.AuthMetadata {
	return nil
}

func NewGetCustomerOrderAction() sdk.Action {
	return &GetCustomerOrderAction{}
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

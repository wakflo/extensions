package actions

import (
	"context"
	"fmt"

	goshopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/wakflo/extensions/internal/integrations/shopify/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type getCustomerOrderActionProps struct {
	CustomerID uint64 `json:"customerId"`
}

type GetCustomerOrderAction struct{}

func (a *GetCustomerOrderAction) Name() string {
	return "Get Customer Order"
}

func (a *GetCustomerOrderAction) Description() string {
	return "Retrieves customer orders from the specified system or database, allowing you to automate tasks that require access to order information."
}

func (a *GetCustomerOrderAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *GetCustomerOrderAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &getCustomerOrderDocs,
	}
}

func (a *GetCustomerOrderAction) Icon() *string {
	return nil
}

func (a *GetCustomerOrderAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"customerId": autoform.NewNumberField().
			SetDisplayName("Customer ID").
			SetDescription("The ID of the customer.").
			SetRequired(true).
			Build(),
	}
}

func (a *GetCustomerOrderAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[getCustomerOrderActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	client, err := shared.CreateClient(ctx.BaseContext)
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

func (a *GetCustomerOrderAction) Auth() *sdk.Auth {
	return nil
}

func (a *GetCustomerOrderAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *GetCustomerOrderAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
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

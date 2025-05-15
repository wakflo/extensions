package actions

import (
	"context"
	"errors"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/shopify/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type listOrdersActionProps struct {
	Limit             int    `json:"limit,omitempty"`
	Status            string `json:"status,omitempty"`
	FinancialStatus   string `json:"financial_status,omitempty"`
	FulfillmentStatus string `json:"fulfillment_status,omitempty"`
}

type ListOrdersAction struct{}

// Metadata returns metadata about the action
func (a *ListOrdersAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "list_orders",
		DisplayName:   "List Orders",
		Description:   "Retrieve a list of orders from your Shopify store, allowing you to automate tasks and workflows based on order data.",
		Type:          core.ActionTypeAction,
		Documentation: listOrdersDocs,
		SampleOutput: map[string]any{
			"orders": []map[string]any{
				{
					"id":           "123456789",
					"name":         "#1001",
					"email":        "customer@example.com",
					"created_at":   "2023-01-01T12:00:00Z",
					"updated_at":   "2023-01-01T12:30:00Z",
					"processed_at": "2023-01-01T12:05:00Z",
					"closed_at":    nil,
					"cancelled_at": nil,
				},
			},
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *ListOrdersAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("list_orders", "List Orders")

	form.NumberField("limit", "Limit").
		Required(false).
		HelpText("Maximum number of orders to return (default: 50, max: 250).")

	form.TextField("status", "Status").
		Required(false).
		Placeholder("any").
		HelpText("Filter orders by status (open, closed, cancelled, any).")

	form.TextField("financial_status", "Financial Status").
		Required(false).
		Placeholder("paid").
		HelpText("Filter orders by financial status (authorized, pending, paid, partially_paid, refunded, voided, partially_refunded, any).")

	form.TextField("fulfillment_status", "Fulfillment Status").
		Required(false).
		Placeholder("fulfilled").
		HelpText("Filter orders by fulfillment status (fulfilled, partial, unfulfilled, any).")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *ListOrdersAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *ListOrdersAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[listOrdersActionProps](ctx)
	if err != nil {
		return nil, err
	}

	// Create Shopify client
	client, err := shared.CreateClient(ctx)
	if err != nil {
		return nil, err
	}

	// Prepare query options
	options := map[string]interface{}{}
	if input.Limit > 0 {
		options["limit"] = input.Limit
	}
	if input.Status != "" {
		options["status"] = input.Status
	}
	if input.FinancialStatus != "" {
		options["financial_status"] = input.FinancialStatus
	}
	if input.FulfillmentStatus != "" {
		options["fulfillment_status"] = input.FulfillmentStatus
	}

	// Get orders from Shopify
	orders, err := client.Order.List(context.Background(), options)
	if err != nil {
		return nil, err
	}

	if orders == nil {
		return nil, errors.New("no orders found")
	}

	return map[string]interface{}{
		"orders": orders,
	}, nil
}

func NewListOrdersAction() sdk.Action {
	return &ListOrdersAction{}
}

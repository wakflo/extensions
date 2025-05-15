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

type listDraftOrdersActionProps struct {
	Limit  int    `json:"limit,omitempty"`
	Status string `json:"status,omitempty"`
}

type ListDraftOrdersAction struct{}

// Metadata returns metadata about the action
func (a *ListDraftOrdersAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "list_draft_orders",
		DisplayName:   "List Draft Orders",
		Description:   "Retrieve and list all draft orders in your Shopify store, allowing you to review and manage pending orders with ease.",
		Type:          core.ActionTypeAction,
		Documentation: listDraftOrdersDocs,
		SampleOutput: map[string]any{
			"draft_orders": []map[string]any{
				{
					"id":              987654321,
					"name":            "#D1",
					"email":           "customer@example.com",
					"created_at":      "2023-01-01T12:00:00Z",
					"updated_at":      "2023-01-01T12:00:00Z",
					"status":          "open",
					"currency":        "USD",
					"total_price":     "100.00",
					"subtotal_price":  "90.00",
					"total_tax":       "10.00",
					"invoice_sent_at": nil,
					"note":            "Customer requested rush processing",
					"line_items": []map[string]any{
						{
							"id":         123456789,
							"product_id": 111222333,
							"variant_id": 444555666,
							"title":      "Example Product",
							"price":      "45.00",
							"quantity":   2,
							"sku":        "PROD-123",
							"grams":      500,
						},
					},
					"shipping_address": map[string]any{
						"address1": "123 Main St",
						"city":     "New York",
						"province": "NY",
						"country":  "US",
						"zip":      "10001",
						"phone":    "+1234567890",
					},
				},
			},
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *ListDraftOrdersAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("list_draft_orders", "List Draft Orders")

	form.NumberField("limit", "Limit").
		Required(false).
		HelpText("Maximum number of draft orders to return (default: 50, max: 250).")

	form.TextField("status", "Status").
		Required(false).
		Placeholder("open").
		HelpText("Filter draft orders by status (open, invoice_sent, completed).")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *ListDraftOrdersAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *ListDraftOrdersAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[listDraftOrdersActionProps](ctx)
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

	// Get draft orders from Shopify
	draftOrders, err := client.DraftOrder.List(context.Background(), options)
	if err != nil {
		return nil, err
	}

	if draftOrders == nil {
		return nil, errors.New("no draft orders found")
	}

	return map[string]interface{}{
		"draft_orders": draftOrders,
	}, nil
}

func NewListDraftOrdersAction() sdk.Action {
	return &ListDraftOrdersAction{}
}

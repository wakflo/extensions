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

type listCustomersActionProps struct {
	Limit  int    `json:"limit,omitempty"`
	Search string `json:"search,omitempty"`
}

type ListCustomersAction struct{}

// Metadata returns metadata about the action
func (a *ListCustomersAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "list_customers",
		DisplayName:   "List Customers",
		Description:   "Retrieves a list of customers from your Shopify store, allowing you to automate tasks that require customer information.",
		Type:          core.ActionTypeAction,
		Documentation: listCustomersDocs,
		SampleOutput: map[string]any{
			"customers": []map[string]any{
				{
					"id":           "123456789",
					"email":        "customer@example.com",
					"first_name":   "John",
					"last_name":    "Doe",
					"phone":        "+1234567890",
					"created_at":   "2023-01-01T12:00:00Z",
					"updated_at":   "2023-01-01T12:00:00Z",
					"orders_count": 5,
					"state":        "enabled",
					"total_spent":  "500.00",
					"tags":         "vip,repeat",
					"addresses": []map[string]any{
						{
							"address1": "123 Main St",
							"city":     "New York",
							"province": "NY",
							"country":  "US",
							"zip":      "10001",
							"phone":    "+1234567890",
							"default":  true,
						},
					},
				},
			},
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *ListCustomersAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("list_customers", "List Customers")

	form.NumberField("limit", "Limit").
		Required(false).
		HelpText("Maximum number of customers to return (default: 50, max: 250).")

	form.TextField("search", "Search Query").
		Required(false).
		Placeholder("john@example.com").
		HelpText("Search for customers by email, name, or phone number.")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *ListCustomersAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *ListCustomersAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[listCustomersActionProps](ctx)
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
	if input.Search != "" {
		options["query"] = input.Search
	}

	// Get customers from Shopify
	customers, err := client.Customer.List(context.Background(), options)
	if err != nil {
		return nil, err
	}

	if customers == nil {
		return nil, errors.New("no customers found")
	}

	return map[string]interface{}{
		"customers": customers,
	}, nil
}

func NewListCustomersAction() sdk.Action {
	return &ListCustomersAction{}
}

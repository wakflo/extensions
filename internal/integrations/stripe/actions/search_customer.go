package actions

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/wakflo/extensions/internal/integrations/stripe/shared"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type searchCustomerActionProps struct {
	Email string `json:"email"`
}

type SearchCustomerAction struct{}

// Metadata returns metadata about the action
func (a *SearchCustomerAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "search_customer",
		DisplayName:   "Search Customer",
		Description:   "Search for a Stripe customer by email address",
		Type:          core.ActionTypeAction,
		Documentation: searchCustomerDocs,
		SampleOutput: []map[string]any{
			{
				"id":          "cus_1234567890",
				"object":      "customer",
				"name":        "John Doe",
				"email":       "john.doe@example.com",
				"phone":       "+1234567890",
				"description": "Customer description",
				"created":     1620000000,
				"address": map[string]string{
					"city":        "San Francisco",
					"country":     "US",
					"line1":       "123 Market St",
					"state":       "CA",
					"postal_code": "94102",
				},
			},
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *SearchCustomerAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("search_customer", "Search Customer")

	form.TextField("email", "Email").
		Placeholder("customer@example.com").
		Required(true).
		HelpText("The customer's email address to search for.")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *SearchCustomerAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *SearchCustomerAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[searchCustomerActionProps](ctx)
	if err != nil {
		return nil, err
	}

	// Get the auth context and API key
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	apiKey := authCtx.Extra["api-key"]
	if apiKey == "" {
		return nil, errors.New("missing stripe secret api-key")
	}

	// Build query parameters
	params := url.Values{}
	params.Add("query", "email:'"+input.Email+"'")

	reqURL := "/v1/customers/search"

	// Call the Stripe API
	resp, err := shared.StripeClient(apiKey, reqURL, http.MethodGet, nil, params)
	if err != nil {
		return nil, err
	}

	// Extract the data from the response
	nodes, ok := resp["data"].([]interface{})
	if !ok {
		return nil, errors.New("failed to extract data from response")
	}

	return nodes, nil
}

func NewSearchCustomerAction() sdk.Action {
	return &SearchCustomerAction{}
}

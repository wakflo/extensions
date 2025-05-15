package actions

import (
	"errors"
	"net/http"

	"github.com/wakflo/extensions/internal/integrations/stripe/shared"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type retrieveCustomerActionProps struct {
	CustomerID string `json:"customer"`
}

type RetrieveCustomerAction struct{}

// Metadata returns metadata about the action
func (a *RetrieveCustomerAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "retrieve_customer",
		DisplayName:   "Retrieve Customer By ID",
		Description:   "Retrieve a Stripe customer by their ID",
		Type:          core.ActionTypeAction,
		Documentation: retrieveCustomerDocs,
		SampleOutput: map[string]any{
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
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *RetrieveCustomerAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("retrieve_customer", "Retrieve Customer")

	form.TextField("customer", "Customer ID").
		Placeholder("cus_1234567890").
		Required(true).
		HelpText("The Stripe customer ID to retrieve.")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *RetrieveCustomerAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *RetrieveCustomerAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[retrieveCustomerActionProps](ctx)
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

	// Build the request URL
	reqURL := "/v1/customers/" + input.CustomerID

	// Call the Stripe API
	resp, err := shared.StripeClient(apiKey, reqURL, http.MethodGet, nil, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func NewRetrieveCustomerAction() sdk.Action {
	return &RetrieveCustomerAction{}
}

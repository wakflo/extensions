package actions

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/stripe/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type createInvoiceActionProps struct {
	CustomerID  string `json:"customer"`
	Currency    string `json:"currency"`
	Description string `json:"description"`
}

type CreateInvoiceAction struct{}

// Metadata returns metadata about the action
func (a *CreateInvoiceAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_invoice",
		DisplayName:   "Create Invoice",
		Description:   "Create a new invoice in Stripe for a specified customer",
		Type:          core.ActionTypeAction,
		Documentation: createInvoiceDocs,
		SampleOutput: map[string]any{
			"id":          "in_1234567890",
			"object":      "invoice",
			"customer":    "cus_1234567890",
			"currency":    "usd",
			"description": "Invoice description",
			"created":     1620000000,
			"status":      "draft",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *CreateInvoiceAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_invoice", "Create Invoice")

	form.TextField("customer", "Customer ID").
		Placeholder("cus_1234567890").
		Required(true).
		HelpText("Stripe customer ID.")

	form.TextField("currency", "Currency").
		Placeholder("USD").
		Required(true).
		HelpText("Currency for the invoice (e.g., USD).")

	form.TextField("description", "Description").
		Placeholder("Invoice description").
		Required(false).
		HelpText("A description for the invoice.")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *CreateInvoiceAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *CreateInvoiceAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[createInvoiceActionProps](ctx)
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

	// Build request data
	data := url.Values{}
	data.Set("customer", input.CustomerID)
	data.Set("currency", input.Currency)
	if input.Description != "" {
		data.Set("description", input.Description)
	}

	payload := []byte(data.Encode())
	reqURL := "/v1/invoices"

	// Assuming stripClient is part of the shared package
	resp, err := shared.StripeClient(apiKey, reqURL, http.MethodPost, payload, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func NewCreateInvoiceAction() sdk.Action {
	return &CreateInvoiceAction{}
}

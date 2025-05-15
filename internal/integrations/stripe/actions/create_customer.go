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

type createCustomerActionProps struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	City        string `json:"city,omitempty"`
	Country     string `json:"country,omitempty"`
	Address     string `json:"line1,omitempty"`
	Phone       string `json:"phone,omitempty"`
	Description string `json:"description,omitempty"`
	State       string `json:"state,omitempty"`
	PostalCode  string `json:"postal_code,omitempty"`
}

type CreateCustomerAction struct{}

// Metadata returns metadata about the action
func (a *CreateCustomerAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_customer",
		DisplayName:   "Create Customer",
		Description:   "Create a new customer in Stripe with specified details",
		Type:          core.ActionTypeAction,
		Documentation: createCustomerDocs,
		SampleOutput: map[string]any{
			"id":          "cus_1234567890",
			"object":      "customer",
			"name":        "John Doe",
			"email":       "john.doe@example.com",
			"phone":       "+1234567890",
			"description": "New customer",
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

func (a *CreateCustomerAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_customer", "Create Customer")

	form.TextField("name", "Name").
		Placeholder("John Doe").
		Required(true).
		HelpText("The customer's full name.")

	form.TextField("email", "Email").
		Placeholder("john.doe@example.com").
		Required(true).
		HelpText("The customer's email address.")

	form.TextField("description", "Description").
		Placeholder("Customer description").
		Required(false).
		HelpText("A description of the customer.")

	form.TextField("phone", "Phone").
		Placeholder("+1234567890").
		Required(false).
		HelpText("The customer's phone number.")

	form.TextField("line1", "Address Line 1").
		Placeholder("123 Market St").
		Required(false).
		HelpText("The first line of the customer's address.")

	form.TextField("postal_code", "Postal Code").
		Placeholder("94102").
		Required(false).
		HelpText("The customer's postal code.")

	form.TextField("city", "City").
		Placeholder("San Francisco").
		Required(false).
		HelpText("The customer's city.")

	form.TextField("state", "State").
		Placeholder("CA").
		Required(false).
		HelpText("The customer's state or province.")

	form.TextField("country", "Country").
		Placeholder("US").
		Required(false).
		HelpText("The customer's country.")

	schema := form.Build()

	return schema
}

func (a *CreateCustomerAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *CreateCustomerAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[createCustomerActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	apiKey := authCtx.Extra["api-key"]
	if apiKey == "" {
		return nil, errors.New("missing stripe secret api-key")
	}

	data := url.Values{}
	data.Set("name", input.Name)
	data.Set("email", input.Email)
	data.Set("description", input.Description)
	if input.Phone != "" {
		data.Set("phone", input.Phone)
	}
	if input.City != "" {
		data.Set("address[city]", input.City)
	}
	if input.Country != "" {
		data.Set("address[country]", input.Country)
	}
	if input.Address != "" {
		data.Set("address[line1]", input.Address)
	}
	if input.State != "" {
		data.Set("address[state]", input.State)
	}
	if input.PostalCode != "" {
		data.Set("address[postal_code]", input.PostalCode)
	}

	payload := []byte(data.Encode())
	reqURL := "/v1/customers"

	resp, err := shared.StripeClient(apiKey, reqURL, http.MethodPost, payload, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func NewCreateCustomerAction() sdk.Action {
	return &CreateCustomerAction{}
}

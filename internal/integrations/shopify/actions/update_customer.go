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

type updateCustomerActionProps struct {
	CustomerID uint64 `json:"customerId"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	Tags       string `json:"tags"`
}

type UpdateCustomerAction struct{}

// Metadata returns metadata about the action
func (a *UpdateCustomerAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "update_customer",
		DisplayName:   "Update Customer",
		Description:   "Updates customer information in your CRM or database by mapping and synchronizing data from various sources, ensuring accurate and up-to-date records.",
		Type:          core.ActionTypeAction,
		Documentation: updateCustomerDocs,
		SampleOutput: map[string]any{
			"updated_customer": map[string]any{},
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *UpdateCustomerAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("update_customer", "Update Customer")

	form.NumberField("customerId", "Customer ID").
		Required(true).
		HelpText("The id of the customer to update.")

	form.TextField("firstName", "First name").
		Required(false).
		HelpText("Customer first name.")

	form.TextField("lastName", "Last name").
		Required(false).
		HelpText("Customer last name.")

	form.TextField("phone", "Phone").
		Required(false).
		HelpText("Customer phone number.")

	form.TextField("email", "Email").
		Required(false).
		HelpText("Customer email address.")

	form.TextField("tags", "Tags").
		Required(false).
		HelpText("A string of comma-separated tags for filtering and search")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *UpdateCustomerAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *UpdateCustomerAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[updateCustomerActionProps](ctx)
	if err != nil {
		return nil, err
	}

	client, err := shared.CreateClient(ctx)
	if err != nil {
		return nil, err
	}

	existingCustomer, err := client.Customer.Get(context.Background(), input.CustomerID, nil)
	if err != nil {
		return nil, err
	}
	if input.FirstName != "" {
		existingCustomer.FirstName = input.FirstName
	}
	if input.LastName != "" {
		existingCustomer.LastName = input.LastName
	}
	if input.Email != "" {
		existingCustomer.Email = input.Email
	}
	if input.Phone != "" {
		existingCustomer.Phone = input.Phone
	}
	if input.Tags != "" {
		existingCustomer.Tags = input.Tags
	}
	updatedCustomer, err := client.Customer.Update(context.Background(), *existingCustomer)
	if err != nil {
		return nil, errors.New("failed to update customer")
	}
	return map[string]interface{}{
		"updated_customer": updatedCustomer,
	}, nil
}

func NewUpdateCustomerAction() sdk.Action {
	return &UpdateCustomerAction{}
}

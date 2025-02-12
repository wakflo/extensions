package actions

import (
	"context"
	"errors"

	"github.com/wakflo/extensions/internal/integrations/shopify/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
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

func (a *UpdateCustomerAction) Name() string {
	return "Update Customer"
}

func (a *UpdateCustomerAction) Description() string {
	return "Updates customer information in your CRM or database by mapping and synchronizing data from various sources, ensuring accurate and up-to-date records."
}

func (a *UpdateCustomerAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *UpdateCustomerAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &updateCustomerDocs,
	}
}

func (a *UpdateCustomerAction) Icon() *string {
	return nil
}

func (a *UpdateCustomerAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"customerId": autoform.NewNumberField().
			SetDisplayName("Customer ID").
			SetDescription("The id of the customer to update.").
			SetRequired(true).
			Build(),
		"firstName": autoform.NewShortTextField().
			SetDisplayName("First name").
			SetDescription("Customer first name.").
			SetRequired(false).
			Build(),
		"lastName": autoform.NewShortTextField().
			SetDisplayName("Last name").
			SetDescription("Customer last name.").
			SetRequired(false).
			Build(),
		"phone": autoform.NewShortTextField().
			SetDisplayName("Phone").
			SetDescription("Customer phone number.").
			SetRequired(false).
			Build(),
		"email": autoform.NewShortTextField().
			SetDisplayName("Email").
			SetDescription("Customer email address.").
			SetRequired(false).
			Build(),
		"tags": autoform.NewLongTextField().
			SetDisplayName("Tags").
			SetDescription("A string of comma-separated tags for filtering and search").
			SetRequired(false).
			Build(),
	}
}

func (a *UpdateCustomerAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[updateCustomerActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	client, err := shared.CreateClient(ctx.BaseContext)
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

func (a *UpdateCustomerAction) Auth() *sdk.Auth {
	return nil
}

func (a *UpdateCustomerAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *UpdateCustomerAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewUpdateCustomerAction() sdk.Action {
	return &UpdateCustomerAction{}
}

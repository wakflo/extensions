package actions

import (
	"context"
	"errors"

	shopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/wakflo/extensions/internal/integrations/shopify/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type createCustomerActionProps struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Tags      string `json:"tags"`
}

type CreateCustomerAction struct{}

func (a *CreateCustomerAction) Name() string {
	return "Create Customer"
}

func (a *CreateCustomerAction) Description() string {
	return "Create a new customer in your CRM system by providing required details such as name, email, phone number, and other relevant information. This integration action allows you to automate the process of creating new customers, reducing manual errors and increasing efficiency."
}

func (a *CreateCustomerAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *CreateCustomerAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &createCustomerDocs,
	}
}

func (a *CreateCustomerAction) Icon() *string {
	return nil
}

func (a *CreateCustomerAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
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

func (a *CreateCustomerAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createCustomerActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	client, err := shared.CreateClient(ctx.BaseContext)
	if err != nil {
		return nil, err
	}
	customer := shopify.Customer{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Tags:      input.Tags,
		Phone:     input.Phone,
	}

	newCustomer, err := client.Customer.Create(context.Background(), customer)
	if err != nil {
		return nil, err
	}
	if newCustomer == nil {
		return nil, errors.New("customer not created! ")
	}
	return map[string]interface{}{
		"new_customer": newCustomer,
	}, nil
}

func (a *CreateCustomerAction) Auth() *sdk.Auth {
	return nil
}

func (a *CreateCustomerAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *CreateCustomerAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewCreateCustomerAction() sdk.Action {
	return &CreateCustomerAction{}
}

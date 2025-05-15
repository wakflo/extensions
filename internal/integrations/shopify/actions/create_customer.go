package actions

import (
	"context"
	"errors"

	shopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/shopify/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type createCustomerActionProps struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Tags      string `json:"tags"`
}

type CreateCustomerAction struct{}

func (a *CreateCustomerAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_customer",
		DisplayName:   "Create Customer",
		Description:   "Create a new customer in your CRM system by providing required details such as name, email, phone number, and other relevant information. This integration action allows you to automate the process of creating new customers, reducing manual errors and increasing efficiency.",
		Type:          core.ActionTypeAction,
		Documentation: createCustomerDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *CreateCustomerAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_customer", "Create Customer")

	form.TextField("firstName", "First name").
		Placeholder("Customer first name.").
		HelpText("Customer first name.")

	form.TextField("lastName", "Last name").
		Placeholder("Customer last name.").
		HelpText("Customer last name.")

	form.TextField("phone", "Phone").
		Placeholder("Customer phone number.").
		HelpText("Customer phone number.")

	form.TextField("email", "Email").
		Placeholder("Customer email address.").
		HelpText("Customer email address.")

	form.TextareaField("tags", "Tags").
		Placeholder("A string of comma-separated tags for filtering and search").
		HelpText("A string of comma-separated tags for filtering and search")

	schema := form.Build()
	return schema
}

func (a *CreateCustomerAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[createCustomerActionProps](ctx)
	if err != nil {
		return nil, err
	}

	client, err := shared.CreateClient(ctx)
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

func (a *CreateCustomerAction) Auth() *core.AuthMetadata {
	return nil
}

func NewCreateCustomerAction() sdk.Action {
	return &CreateCustomerAction{}
}

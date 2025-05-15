package actions

import (
	"github.com/hiscaler/woocommerce-go"
	"github.com/hiscaler/woocommerce-go/entity"
	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/woocommerce/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type createCustomerActionProps struct {
	Email         string `json:"email,omitempty"`
	FirstName     string `json:"first_name,omitempty"`
	LastName      string `json:"last_name,omitempty"`
	Username      string `json:"username,omitempty"`
	Password      string `json:"password,omitempty"`
	Phone         string `json:"phone,omitempty"`
	City          string `json:"city,omitempty"`
	Country       string `json:"country,omitempty"`
	State         string `json:"state,omitempty"`
	StreetAddress string `json:"street_address,omitempty"`
	PostalCode    string `json:"post_code,omitempty"`
}

type CreateCustomerAction struct{}

func (a *CreateCustomerAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_customer",
		DisplayName:   "Create Customer",
		Description:   "Create a new customer on Woocommerce by providing required details such as name, email, phone number, and other relevant information.",
		Type:          core.ActionTypeAction,
		Documentation: createCustomerDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *CreateCustomerAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_custpmer", "Create Customer")

	form.TextField("first_name", "First Name").
		Placeholder("Enter first name").
		Required(true).
		HelpText("Enter first name")

	form.TextField("last_name", "Last Name").
		Placeholder("Enter last name").
		Required(true).
		HelpText("Enter last name")

	form.TextField("email", "Email").
		Placeholder("Enter email").
		Required(true).
		HelpText("Enter email")

	form.TextField("username", "Username").
		Placeholder("Enter username").
		Required(true).
		HelpText("Enter username")

	form.TextField("password", "Password").
		Placeholder("Enter password").
		Required(true).
		HelpText("Enter password")

	form.TextField("phone", "Phone").
		Placeholder("Enter phone").
		Required(true).
		HelpText("Enter phone")

	form.TextField("country", "Country").
		Placeholder("Enter country").
		Required(true).
		HelpText("Enter country")

	form.TextField("city", "City").
		Placeholder("Enter city").
		Required(true).
		HelpText("Enter city")

	form.TextField("state", "State").
		Placeholder("Enter state").
		Required(true).
		HelpText("Enter state")

	form.TextField("street_address", "Street Address").
		Placeholder("Enter street address").
		Required(true).
		HelpText("Enter street address")

	form.TextField("post_code", "Post Code").
		Placeholder("Enter post code").
		Required(true).
		HelpText("Enter post code")

	labelsArray := form.ArrayField("labels", "Labels")
	labelGroup := labelsArray.ObjectTemplate("label", "")
	labelGroup.TextField("value", "Label").
		Placeholder("Label").
		Required(true).
		HelpText("The task's labels (a list of names that may represent either personal or shared labels)")

	schema := form.Build()

	return schema
}

func (a *CreateCustomerAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[createCustomerActionProps](ctx)
	if err != nil {
		return nil, err
	}

	wooClient, err := shared.InitClient(ctx)
	if err != nil {
		return nil, err
	}

	billing := entity.Billing{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Address1:  input.StreetAddress,
		City:      input.City,
		State:     input.State,
		Postcode:  input.PostalCode,
		Country:   input.Country,
		Email:     input.Email,
		Phone:     input.Phone,
	}

	// Create a query parameters struct
	params := woocommerce.CreateCustomerRequest{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Username:  input.Username,
		Password:  input.Password,
		Billing:   &billing,
	}

	newCustomer, err := wooClient.Services.Customer.Create(params)
	if err != nil {
		return nil, err
	}

	return newCustomer, nil
}

func (a *CreateCustomerAction) Auth() *core.AuthMetadata {
	return nil
}

func NewCreateCustomerAction() sdk.Action {
	return &CreateCustomerAction{}
}

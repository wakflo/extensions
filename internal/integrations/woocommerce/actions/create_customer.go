package actions

import (
	"github.com/hiscaler/woocommerce-go"
	"github.com/hiscaler/woocommerce-go/entity"
	"github.com/wakflo/extensions/internal/integrations/woocommerce/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
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
		"first_name": autoform.NewShortTextField().
			SetDisplayName("First Name").
			SetDescription("Enter first name").
			SetRequired(true).
			Build(),
		"last_name": autoform.NewShortTextField().
			SetDisplayName("Last Name").
			SetRequired(true).
			Build(),
		"email": autoform.NewShortTextField().
			SetDisplayName(" Email").
			SetDescription("Enter email address").
			SetRequired(true).
			Build(),
		"username": autoform.NewShortTextField().
			SetDisplayName("Username").
			SetDescription("Enter username").
			SetRequired(true).
			Build(),
		"password": autoform.NewShortTextField().
			SetDisplayName("Password").
			SetDescription("Enter Password").
			SetRequired(true).
			Build(),
		"phone": autoform.NewShortTextField().
			SetDisplayName("Phone").
			SetDescription("Enter Phone number").
			Build(),
		"country": autoform.NewShortTextField().
			SetDisplayName("Country").
			SetDescription("Enter Country").
			SetRequired(true).
			Build(),
		"city": autoform.NewShortTextField().
			SetDisplayName("City").
			SetRequired(true).
			SetDescription("Enter City").
			Build(),
		"state": autoform.NewShortTextField().
			SetDisplayName("State").
			SetRequired(true).
			SetDescription("Enter State").
			Build(),
		"street_address": autoform.NewLongTextField().
			SetDisplayName("Address").
			SetDescription("Enter the street address").
			SetRequired(true).
			Build(),
		"post_code": autoform.NewShortTextField().
			SetDisplayName("Postal Code").
			SetDescription("Enter State").
			Build(),
	}
}

func (a *CreateCustomerAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createCustomerActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	wooClient, err := shared.InitClient(ctx.BaseContext)
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

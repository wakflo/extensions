package actions

import (
	"fmt"

	"github.com/hiscaler/woocommerce-go"
	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/woocommerce/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type findCustomerActionProps struct {
	Email string `json:"email"`
}

type FindCustomerAction struct{}

func (a *FindCustomerAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "find_customer",
		DisplayName:   "Find Customer By Email",
		Description:   "Searches for a customer by their email address .",
		Type:          core.ActionTypeAction,
		Documentation: findCustomerDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *FindCustomerAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("find_customer", "Find Customer")

	form.TextField("email", "Email").
		Placeholder("Enter the email address of the customer").
		Required(true).
		HelpText("Enter the email address of the customer")

	schema := form.Build()

	return schema
}

func (a *FindCustomerAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[findCustomerActionProps](ctx)
	if err != nil {
		return nil, err
	}

	wooClient, err := shared.InitClient(ctx)
	if err != nil {
		return nil, err
	}

	params := woocommerce.CustomersQueryParams{
		Email: input.Email,
	}

	customers, _, _, _, err := wooClient.Services.Customer.All(params)
	if err != nil {
		return nil, err
	}

	// Filter customers by exact email match
	for _, customer := range customers {
		if customer.Email == input.Email {
			return customer, nil
		}
	}

	return nil, fmt.Errorf("no customer found with email: %s", input.Email)
}

func (a *FindCustomerAction) Auth() *core.AuthMetadata {
	return nil
}

func NewFindCustomerAction() sdk.Action {
	return &FindCustomerAction{}
}

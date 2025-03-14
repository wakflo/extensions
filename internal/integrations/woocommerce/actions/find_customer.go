package actions

import (
	"fmt"

	"github.com/hiscaler/woocommerce-go"
	"github.com/wakflo/extensions/internal/integrations/woocommerce/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type findCustomerActionProps struct {
	Email string `json:"email"`
}

type FindCustomerAction struct{}

func (a *FindCustomerAction) Name() string {
	return "Find Customer"
}

func (a *FindCustomerAction) Description() string {
	return "Searches for a customer by their unique identifier (e.g., email address or customer ID) and retrieves relevant information, such as name, contact details, and account history."
}

func (a *FindCustomerAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *FindCustomerAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &findCustomerDocs,
	}
}

func (a *FindCustomerAction) Icon() *string {
	return nil
}

func (a *FindCustomerAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"email": autoform.NewShortTextField().
			SetDisplayName("Email").
			SetDescription("Enter the email address of the customer").
			SetRequired(true).
			Build(),
	}
}

func (a *FindCustomerAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[findCustomerActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	wooClient, err := shared.InitClient(ctx.BaseContext)
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

func (a *FindCustomerAction) Auth() *sdk.Auth {
	return nil
}

func (a *FindCustomerAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *FindCustomerAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewFindCustomerAction() sdk.Action {
	return &FindCustomerAction{}
}

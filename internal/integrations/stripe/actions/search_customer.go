package actions

import (
	"fmt"

	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type searchCustomerActionProps struct {
	Name string `json:"name"`
}

type SearchCustomerAction struct{}

func (a *SearchCustomerAction) Name() string {
	return "Search Customer"
}

func (a *SearchCustomerAction) Description() string {
	return "Searches for a customer by their name, email, or phone number in your CRM system and retrieves relevant information such as contact details, order history, and account status."
}

func (a *SearchCustomerAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *SearchCustomerAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &searchCustomerDocs,
	}
}

func (a *SearchCustomerAction) Icon() *string {
	return nil
}

func (a *SearchCustomerAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"name": autoform.NewShortTextField().
			SetLabel("Name").
			SetRequired(true).
			SetPlaceholder("Your name").
			Build(),
	}
}

func (a *SearchCustomerAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[searchCustomerActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	// implement action logic
	out := map[string]any{
		"message": fmt.Sprintf("Hello %s!", input.Name),
	}

	return out, nil
}

func (a *SearchCustomerAction) Auth() *sdk.Auth {
	return nil
}

func (a *SearchCustomerAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *SearchCustomerAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewSearchCustomerAction() sdk.Action {
	return &SearchCustomerAction{}
}

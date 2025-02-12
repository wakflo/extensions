package actions

import (
	"fmt"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type createCustomerActionProps struct {
	Name string `json:"name"`
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
		"name": autoform.NewShortTextField().
			SetLabel("Name").
			SetRequired(true).
			SetPlaceholder("Your name").
			Build(),
	}
}

func (a *CreateCustomerAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createCustomerActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	// implement action logic
	out := map[string]any{
		"message": fmt.Sprintf("Hello %s!", input.Name),
	}

	return out, nil
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

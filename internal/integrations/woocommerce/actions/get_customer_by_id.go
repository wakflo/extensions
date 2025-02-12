package actions

import (
	"github.com/wakflo/extensions/internal/integrations/woocommerce/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type getCustomerByIdActionProps struct {
	CustomerID int `json:"customer-id"`
}

type GetCustomerByIdAction struct{}

func (a *GetCustomerByIdAction) Name() string {
	return "Get Customer By ID"
}

func (a *GetCustomerByIdAction) Description() string {
	return "Retrieves a customer record by their unique identifier (ID) from your CRM or database, allowing you to access and utilize customer information in subsequent workflow steps."
}

func (a *GetCustomerByIdAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *GetCustomerByIdAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &getCustomerByIdDocs,
	}
}

func (a *GetCustomerByIdAction) Icon() *string {
	return nil
}

func (a *GetCustomerByIdAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"customer-id": autoform.NewNumberField().
			SetDisplayName("Customer ID").
			SetDescription("Enter customer ID").
			SetRequired(true).
			Build(),
	}
}

func (a *GetCustomerByIdAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[getCustomerByIdActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	wooClient, err := shared.InitClient(ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	customer, err := wooClient.Services.Customer.One(input.CustomerID)
	if err != nil {
		return nil, err
	}

	return customer, nil
}

func (a *GetCustomerByIdAction) Auth() *sdk.Auth {
	return nil
}

func (a *GetCustomerByIdAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *GetCustomerByIdAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewGetCustomerByIdAction() sdk.Action {
	return &GetCustomerByIdAction{}
}

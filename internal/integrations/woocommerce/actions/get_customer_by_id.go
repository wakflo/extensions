package actions

import (
	"github.com/wakflo/extensions/internal/integrations/woocommerce/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type getCustomerByIDActionProps struct {
	CustomerID int `json:"customer-id"`
}

type GetCustomerByIDAction struct{}

func (a *GetCustomerByIDAction) Name() string {
	return "Get Customer By ID"
}

func (a *GetCustomerByIDAction) Description() string {
	return "Retrieves a customer record by their unique identifier (ID) from your CRM or database, allowing you to access and utilize customer information in subsequent workflow steps."
}

func (a *GetCustomerByIDAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *GetCustomerByIDAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &getCustomerByIDDocs,
	}
}

func (a *GetCustomerByIDAction) Icon() *string {
	return nil
}

func (a *GetCustomerByIDAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"customer-id": autoform.NewNumberField().
			SetDisplayName("Customer ID").
			SetDescription("Enter customer ID").
			SetRequired(true).
			Build(),
	}
}

func (a *GetCustomerByIDAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[getCustomerByIDActionProps](ctx.BaseContext)
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

func (a *GetCustomerByIDAction) Auth() *sdk.Auth {
	return nil
}

func (a *GetCustomerByIDAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *GetCustomerByIDAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewGetCustomerByIDAction() sdk.Action {
	return &GetCustomerByIDAction{}
}

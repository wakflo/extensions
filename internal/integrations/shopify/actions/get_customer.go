package actions

import (
	"context"
	"fmt"

	"github.com/wakflo/extensions/internal/integrations/shopify/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type getCustomerActionProps struct {
	CustomerID uint64 `json:"customerId"`
}

type GetCustomerAction struct{}

func (a *GetCustomerAction) Name() string {
	return "Get Customer"
}

func (a *GetCustomerAction) Description() string {
	return "Retrieves customer information from a specified data source or system, allowing you to access and utilize existing customer data within your workflow."
}

func (a *GetCustomerAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *GetCustomerAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &getCustomerDocs,
	}
}

func (a *GetCustomerAction) Icon() *string {
	return nil
}

func (a *GetCustomerAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"customerId": autoform.NewNumberField().
			SetDisplayName("Customer ID").
			SetDescription("The ID of the customer.").
			SetRequired(true).
			Build(),
	}
}

func (a *GetCustomerAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[getCustomerActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	client, err := shared.CreateClient(ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	customer, err := client.Customer.Get(context.Background(), input.CustomerID, nil)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, fmt.Errorf("no customer found with ID '%d'", input.CustomerID)
	}

	return sdk.JSON(map[string]interface{}{
		"raw_customer": customer,
	}), nil
}

func (a *GetCustomerAction) Auth() *sdk.Auth {
	return nil
}

func (a *GetCustomerAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *GetCustomerAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewGetCustomerAction() sdk.Action {
	return &GetCustomerAction{}
}

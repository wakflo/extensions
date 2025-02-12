package actions

import (
	"context"
	"errors"

	"github.com/wakflo/extensions/internal/integrations/shopify/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type listCustomersActionProps struct {
	Name string `json:"name"`
}

type ListCustomersAction struct{}

func (a *ListCustomersAction) Name() string {
	return "List Customers"
}

func (a *ListCustomersAction) Description() string {
	return "Retrieves a list of customers from your CRM or database, allowing you to automate tasks that require customer information."
}

func (a *ListCustomersAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *ListCustomersAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &listCustomersDocs,
	}
}

func (a *ListCustomersAction) Icon() *string {
	return nil
}

func (a *ListCustomersAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{}
}

func (a *ListCustomersAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	client, err := shared.CreateClient(ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	customers, err := client.Customer.List(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	if customers == nil {
		return nil, errors.New("no customer found")
	}

	return sdk.JSON(map[string]interface{}{
		"customers details": customers,
	}), err
}

func (a *ListCustomersAction) Auth() *sdk.Auth {
	return nil
}

func (a *ListCustomersAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *ListCustomersAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewListCustomersAction() sdk.Action {
	return &ListCustomersAction{}
}

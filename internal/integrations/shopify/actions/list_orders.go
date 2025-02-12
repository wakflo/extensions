package actions

import (
	"context"
	"errors"

	"github.com/wakflo/extensions/internal/integrations/shopify/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type listOrdersActionProps struct {
	Name string `json:"name"`
}

type ListOrdersAction struct{}

func (a *ListOrdersAction) Name() string {
	return "List Orders"
}

func (a *ListOrdersAction) Description() string {
	return "Retrieve a list of orders from your e-commerce platform or order management system, allowing you to automate tasks and workflows based on order data."
}

func (a *ListOrdersAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *ListOrdersAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &listOrdersDocs,
	}
}

func (a *ListOrdersAction) Icon() *string {
	return nil
}

func (a *ListOrdersAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{}
}

func (a *ListOrdersAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	client, err := shared.CreateClient(ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	orders, err := client.Order.List(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	if orders == nil {
		return nil, errors.New("no orders found")
	}

	return sdk.JSON(map[string]interface{}{
		"Orders": orders,
	}), err
}

func (a *ListOrdersAction) Auth() *sdk.Auth {
	return nil
}

func (a *ListOrdersAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *ListOrdersAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewListOrdersAction() sdk.Action {
	return &ListOrdersAction{}
}

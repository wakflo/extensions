package actions

import (
	"github.com/hiscaler/woocommerce-go"
	"github.com/wakflo/extensions/internal/integrations/woocommerce/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type listOrdersActionProps struct {
	ProjectID int `json:"project-id"`
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
	return map[string]*sdkcore.AutoFormSchema{
		"projectId": autoform.NewLongTextField().
			SetDisplayName("").
			SetDescription("").
			Build(),
	}
}

func (a *ListOrdersAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	wooClient, err := shared.InitClient(ctx.BaseContext)
	if err != nil {
		return nil, err
	}
	params := woocommerce.OrdersQueryParams{}

	orders, _, _, _, err := wooClient.Services.Order.All(params)
	if err != nil {
		return nil, err
	}

	return orders, nil
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

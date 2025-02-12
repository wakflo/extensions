package actions

import (
	"context"
	"errors"

	"github.com/wakflo/extensions/internal/integrations/shopify/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type listDraftOrdersActionProps struct {
	Name string `json:"name"`
}

type ListDraftOrdersAction struct{}

func (a *ListDraftOrdersAction) Name() string {
	return "List Draft Orders"
}

func (a *ListDraftOrdersAction) Description() string {
	return "Retrieve and list all draft orders in your e-commerce platform, allowing you to review and manage pending orders with ease."
}

func (a *ListDraftOrdersAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *ListDraftOrdersAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &listDraftOrdersDocs,
	}
}

func (a *ListDraftOrdersAction) Icon() *string {
	return nil
}

func (a *ListDraftOrdersAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{}
}

func (a *ListDraftOrdersAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	client, err := shared.CreateClient(ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	darftOrders, err := client.DraftOrder.List(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	if darftOrders == nil {
		return nil, errors.New("no draft orders found")
	}

	return sdk.JSON(map[string]interface{}{
		"Draft orders": darftOrders,
	}), err
}

func (a *ListDraftOrdersAction) Auth() *sdk.Auth {
	return nil
}

func (a *ListDraftOrdersAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *ListDraftOrdersAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewListDraftOrdersAction() sdk.Action {
	return &ListDraftOrdersAction{}
}

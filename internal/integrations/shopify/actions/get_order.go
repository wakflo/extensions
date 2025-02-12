package actions

import (
	"context"
	"fmt"

	"github.com/wakflo/extensions/internal/integrations/shopify/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type getOrderActionProps struct {
	OrderID uint64 `json:"orderId"`
}

type GetOrderAction struct{}

func (a *GetOrderAction) Name() string {
	return "Get Order"
}

func (a *GetOrderAction) Description() string {
	return "Retrieves an order from the system, allowing you to access and manipulate order details within your workflow."
}

func (a *GetOrderAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *GetOrderAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &getOrderDocs,
	}
}

func (a *GetOrderAction) Icon() *string {
	return nil
}

func (a *GetOrderAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"orderId": autoform.NewNumberField().
			SetDisplayName("Order ID").
			SetDescription("The ID of the order.").
			SetRequired(true).
			Build(),
	}
}

func (a *GetOrderAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[getOrderActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	client, err := shared.CreateClient(ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	order, err := client.Order.Get(context.Background(), input.OrderID, nil)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, fmt.Errorf("no order found with ID '%d'", input.OrderID)
	}
	return sdk.JSON(map[string]interface{}{
		"order details": order,
	}), nil
}

func (a *GetOrderAction) Auth() *sdk.Auth {
	return nil
}

func (a *GetOrderAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *GetOrderAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewGetOrderAction() sdk.Action {
	return &GetOrderAction{}
}

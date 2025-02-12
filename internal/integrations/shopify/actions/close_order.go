package actions

import (
	"context"
	"fmt"

	"github.com/wakflo/extensions/internal/integrations/shopify/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type closeOrderActionProps struct {
	OrderID uint64 `json:"orderId"`
}

type CloseOrderAction struct{}

func (a *CloseOrderAction) Name() string {
	return "Close Order"
}

func (a *CloseOrderAction) Description() string {
	return "Automatically closes an order in your system, marking it as fulfilled and updating relevant fields to reflect the order's status."
}

func (a *CloseOrderAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *CloseOrderAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &closeOrderDocs,
	}
}

func (a *CloseOrderAction) Icon() *string {
	return nil
}

func (a *CloseOrderAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"orderId": autoform.NewNumberField().
			SetDisplayName("Order").
			SetDescription("The ID of the order to close.").
			SetRequired(true).
			Build(),
	}
}

func (a *CloseOrderAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[closeOrderActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	client, err := shared.CreateClient(ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	order, err := client.Order.Close(context.Background(), input.OrderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, fmt.Errorf("no order found with ID '%d'", input.OrderID)
	}
	orderMap := map[string]interface{}{
		"ID":    order.Id,
		"Email": order.Email,
		"Note":  order.Note,
	}
	return sdk.JSON(map[string]interface{}{
		"Closed order details": orderMap,
	}), nil
}

func (a *CloseOrderAction) Auth() *sdk.Auth {
	return nil
}

func (a *CloseOrderAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *CloseOrderAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewCloseOrderAction() sdk.Action {
	return &CloseOrderAction{}
}

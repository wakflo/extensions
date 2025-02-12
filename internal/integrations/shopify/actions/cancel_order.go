package actions

import (
	"context"
	"fmt"

	"github.com/wakflo/extensions/internal/integrations/shopify/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type cancelOrderActionProps struct {
	OrderID uint64 `json:"orderId"`
	Reason  string `json:"reason,omitempty"`
}

type CancelOrderAction struct{}

func (a *CancelOrderAction) Name() string {
	return "Cancel Order"
}

func (a *CancelOrderAction) Description() string {
	return "Cancels an existing order, revoking any associated payment processing and updating the order status to reflect cancellation."
}

func (a *CancelOrderAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *CancelOrderAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &cancelOrderDocs,
	}
}

func (a *CancelOrderAction) Icon() *string {
	return nil
}

func (a *CancelOrderAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"orderId": autoform.NewNumberField().
			SetDisplayName("Order").
			SetDescription("The ID of the order to cancel.").
			SetRequired(true).
			Build(),
	}
}

func (a *CancelOrderAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[cancelOrderActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	client, err := shared.CreateClient(ctx.BaseContext)
	if err != nil {
		return nil, err
	}
	order, err := client.Order.Cancel(context.Background(), input.OrderID, nil)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, fmt.Errorf("no order found with ID '%d'", input.OrderID)
	}
	return sdk.JSON(map[string]interface{}{
		"Result": "Order Successfully cancelled!",
	}), nil
}

func (a *CancelOrderAction) Auth() *sdk.Auth {
	return nil
}

func (a *CancelOrderAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *CancelOrderAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewCancelOrderAction() sdk.Action {
	return &CancelOrderAction{}
}

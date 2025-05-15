package actions

import (
	"context"
	"fmt"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/shopify/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type cancelOrderActionProps struct {
	OrderID uint64 `json:"orderId"`
	Reason  string `json:"reason,omitempty"`
}

type CancelOrderAction struct{}

func (a *CancelOrderAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "cancel_order",
		DisplayName:   "Cancel Order",
		Description:   "Cancels an existing order, revoking any associated payment processing and updating the order status to reflect cancellation.",
		Type:          core.ActionTypeAction,
		Documentation: cancelOrderDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *CancelOrderAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("cancel_order", "Cancel Order")

	form.NumberField("orderId", "Order").
		Placeholder("The ID of the order to cancel.").
		Required(true).
		HelpText("The ID of the order to cancel.")

	schema := form.Build()
	return schema
}

func (a *CancelOrderAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[cancelOrderActionProps](ctx)
	if err != nil {
		return nil, err
	}

	client, err := shared.CreateClient(ctx)
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
	return core.JSON(map[string]interface{}{
		"Result": "Order Successfully cancelled!",
	}), nil
}

func (a *CancelOrderAction) Auth() *core.AuthMetadata {
	return nil
}

func NewCancelOrderAction() sdk.Action {
	return &CancelOrderAction{}
}

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

type closeOrderActionProps struct {
	OrderID uint64 `json:"orderId"`
}

type CloseOrderAction struct{}

func (a *CloseOrderAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "close_order",
		DisplayName:   "Close Order",
		Description:   "Automatically closes an order in your system, marking it as fulfilled and updating relevant fields to reflect the order's status.",
		Type:          core.ActionTypeAction,
		Documentation: closeOrderDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *CloseOrderAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("close_order", "Close Order")

	form.NumberField("orderId", "Order").
		Placeholder("The ID of the order to close.").
		Required(true).
		HelpText("The ID of the order to close.")

	schema := form.Build()
	return schema
}

func (a *CloseOrderAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[closeOrderActionProps](ctx)
	if err != nil {
		return nil, err
	}

	client, err := shared.CreateClient(ctx)
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
	return core.JSON(map[string]interface{}{
		"Closed order details": orderMap,
	}), nil
}

func (a *CloseOrderAction) Auth() *core.AuthMetadata {
	return nil
}

func NewCloseOrderAction() sdk.Action {
	return &CloseOrderAction{}
}

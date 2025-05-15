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

type getOrderActionProps struct {
	OrderID uint64 `json:"orderId"`
}

type GetOrderAction struct{}

func (a *GetOrderAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_order",
		DisplayName:   "Get Order",
		Description:   "Retrieves an order from the system, allowing you to access and manipulate order details within your workflow.",
		Type:          core.ActionTypeAction,
		Documentation: getOrderDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *GetOrderAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_order", "Get Order")

	form.NumberField("orderId", "Order ID").
		Placeholder("The ID of the order.").
		Required(true).
		HelpText("The ID of the order.")

	schema := form.Build()
	return schema
}

func (a *GetOrderAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[getOrderActionProps](ctx)
	if err != nil {
		return nil, err
	}

	client, err := shared.CreateClient(ctx)
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
	return core.JSON(map[string]interface{}{
		"order details": order,
	}), nil
}

func (a *GetOrderAction) Auth() *core.AuthMetadata {
	return nil
}

func NewGetOrderAction() sdk.Action {
	return &GetOrderAction{}
}

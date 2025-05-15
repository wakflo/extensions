package actions

import (
	"context"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/shopify/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type updateOrderActionProps struct {
	OrderID uint64 `json:"orderId"`
	Tags    string `json:"tags"`
	Note    string `json:"note"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
}

type UpdateOrderAction struct{}

// Metadata returns metadata about the action
func (a *UpdateOrderAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "update_order",
		DisplayName:   "Update Order",
		Description:   "Updates an existing order with new information, such as customer details, shipping address, or payment method.",
		Type:          core.ActionTypeAction,
		Documentation: updateOrderDocs,
		SampleOutput: map[string]any{
			"updated_order": map[string]any{},
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *UpdateOrderAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("update_order", "Update Order")

	form.NumberField("orderId", "Order ID").
		Required(true).
		HelpText("The ID of the order to update")

	form.TextField("note", "Note about the order.").
		Required(false)

	form.TextField("tags", "Tags").
		Required(false).
		HelpText("A string of comma-separated tags for filtering and search.")

	form.TextField("phone", "Phone number.").
		Required(false)

	form.TextField("email", "Email address.").
		Required(false)

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *UpdateOrderAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *UpdateOrderAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[updateOrderActionProps](ctx)
	if err != nil {
		return nil, err
	}

	client, err := shared.CreateClient(ctx)
	if err != nil {
		return nil, err
	}
	existingOrder, err := client.Order.Get(context.Background(), input.OrderID, nil)
	if err != nil {
		return nil, err
	}
	if input.Note != "" {
		existingOrder.Note = input.Note
	}
	if input.Tags != "" {
		existingOrder.Tags = input.Tags
	}
	if input.Phone != "" {
		existingOrder.Phone = input.Phone
	}
	if input.Email != "" {
		existingOrder.Email = input.Email
	}
	updatedOrder, err := client.Order.Update(context.Background(), *existingOrder)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"updated_order": updatedOrder,
	}, nil
}

func NewUpdateOrderAction() sdk.Action {
	return &UpdateOrderAction{}
}

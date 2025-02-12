package actions

import (
	"context"

	"github.com/wakflo/extensions/internal/integrations/shopify/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type updateOrderActionProps struct {
	OrderID uint64 `json:"orderId"`
	Tags    string `json:"tags"`
	Note    string `json:"note"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
}

type UpdateOrderAction struct{}

func (a *UpdateOrderAction) Name() string {
	return "Update Order"
}

func (a *UpdateOrderAction) Description() string {
	return "Updates an existing order with new information, such as customer details, shipping address, or payment method."
}

func (a *UpdateOrderAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *UpdateOrderAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &updateOrderDocs,
	}
}

func (a *UpdateOrderAction) Icon() *string {
	return nil
}

func (a *UpdateOrderAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"orderId": autoform.NewNumberField().
			SetDisplayName("Order ID").
			SetDescription("The ID of the order to update").
			SetRequired(true).
			Build(),
		"note": autoform.NewLongTextField().
			SetDisplayName("Note about the order.").
			SetRequired(false).
			Build(),
		"tags": autoform.NewLongTextField().
			SetDisplayName("Tags").
			SetDescription("A string of comma-separated tags for filtering and search.").
			SetRequired(false).
			Build(),
		"phone": autoform.NewShortTextField().
			SetDisplayName("Phone number.").
			SetRequired(false).
			Build(),
		"email": autoform.NewShortTextField().
			SetDisplayName("Email address.").
			SetRequired(false).
			Build(),
	}
}

func (a *UpdateOrderAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[updateOrderActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	client, err := shared.CreateClient(ctx.BaseContext)
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

func (a *UpdateOrderAction) Auth() *sdk.Auth {
	return nil
}

func (a *UpdateOrderAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *UpdateOrderAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewUpdateOrderAction() sdk.Action {
	return &UpdateOrderAction{}
}

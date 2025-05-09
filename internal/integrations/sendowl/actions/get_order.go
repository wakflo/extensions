package actions

import (
	"github.com/wakflo/extensions/internal/integrations/sendowl/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type getOrderActionProps struct {
	OrderID string `json:"order_id"`
}

type GetOrderAction struct{}

func (a *GetOrderAction) Name() string {
	return "Get Order"
}

func (a *GetOrderAction) Description() string {
	return "Retrieve detailed information about a specific order from your SendOwl account using the order ID."
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
	icon := "mdi:receipt-text"
	return &icon
}

func (a *GetOrderAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"order_id": shared.GetOrderInput("Order ID", "Enter the ID of the order you want to retrieve.", true),
	}
}

func (a *GetOrderAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[getOrderActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	url := "/orders/" + input.OrderID

	order, err := shared.GetSendOwlClient(shared.AltBaseURL, ctx.Auth.Extra["api_key"], ctx.Auth.Extra["api_secret"], url)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (a *GetOrderAction) Auth() *sdk.Auth {
	return nil
}

func (a *GetOrderAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"id":             123456,
		"buyer_name":     "John Doe",
		"buyer_email":    "john.doe@example.com",
		"total":          29.99,
		"currency":       "USD",
		"status":         "completed",
		"payment_method": "credit_card",
		"transaction_id": "tx_abc123def456",
		"ip_address":     "192.168.1.1",
		"country":        "US",
		"state":          "CA",
		"city":           "San Francisco",
		"postal_code":    "94105",
		"address_line1":  "123 Main St",
		"address_line2":  "Apt 4B",
		"created_at":     "2023-08-15T14:30:22Z",
		"completed_at":   "2023-08-15T14:35:45Z",
		"products": []map[string]any{
			{
				"id":              789012,
				"name":            "Digital Marketing Guide",
				"price":           29.99,
				"product_type":    "pdf",
				"download_count":  1,
				"download_limit":  5,
				"download_expiry": "2023-09-15T14:35:45Z",
			},
		},
	}
}

func (a *GetOrderAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewGetOrderAction() sdk.Action {
	return &GetOrderAction{}
}

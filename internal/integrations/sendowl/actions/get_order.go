package actions

import (
	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/sendowl/shared"

	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type getOrderActionProps struct {
	OrderID string `json:"order_id"`
}

type GetOrderAction struct{}

func (a *GetOrderAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		DisplayName:   "Get Order",
		Description:   "Retrieve detailed information about a specific order from your SendOwl account using the order ID.",
		Type:          sdkcore.ActionTypeAction,
		Documentation: getOrderDocs,
		SampleOutput: map[string]any{
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
		},
		Settings: sdkcore.ActionSettings{},
	}
}

func (a *GetOrderAction) Properties() *smartform.FormSchema {

	form := smartform.NewForm("get_order", "Get Order")

	shared.GetOrderProp("order_id", "Order ID", "Enter the ID of the order you want to retrieve.", true, form)

	schema := form.Build()

	return schema

}

func (a *GetOrderAction) Perform(ctx sdkcontext.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[getOrderActionProps](ctx)
	if err != nil {
		return nil, err
	}

	// Get the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	url := "/orders/" + input.OrderID

	order, err := shared.GetSendOwlClient(shared.AltBaseURL, authCtx.Extra["api_key"], authCtx.Extra["api_secret"], url)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (a *GetOrderAction) Auth() *sdkcore.AuthMetadata {
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

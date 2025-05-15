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

type getTransactionActionProps struct {
	OrderID       uint64 `json:"orderId"`
	TransactionID uint64 `json:"transactionId"`
}

type GetTransactionAction struct{}

func (a *GetTransactionAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_transaction",
		DisplayName:   "Get Transaction",
		Description:   "Retrieves transaction details from a Shopify order, allowing you to access and utilize transactional data within your workflow.",
		Type:          core.ActionTypeAction,
		Documentation: getTransactionDocs,
		SampleOutput: map[string]any{
			"transaction": map[string]any{
				"id":                 123456789,
				"order_id":           987654321,
				"amount":             "100.00",
				"kind":               "sale",
				"gateway":            "shopify_payments",
				"status":             "success",
				"message":            "Approved",
				"created_at":         "2023-01-01T12:00:00Z",
				"processed_at":       "2023-01-01T12:01:00Z",
				"currency":           "USD",
				"authorization":      "auth_123456",
				"error_code":         nil,
				"source_name":        "web",
				"payment_details":    map[string]string{},
				"receipt":            map[string]string{},
				"administration_url": "https://example.myshopify.com/admin/orders/987654321/transactions/123456789",
			},
		},
		Settings: core.ActionSettings{},
	}
}

func (a *GetTransactionAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_transaction", "Get Transaction")

	form.NumberField("orderId", "Order ID").
		Required(true).
		HelpText("The ID of the Shopify order.")

	form.NumberField("transactionId", "Transaction ID").
		Required(true).
		HelpText("The ID of the transaction to retrieve.")

	schema := form.Build()

	return schema
}

func (a *GetTransactionAction) Auth() *core.AuthMetadata {
	return nil
}

func (a *GetTransactionAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[getTransactionActionProps](ctx)
	if err != nil {
		return nil, err
	}

	client, err := shared.CreateClient(ctx)
	if err != nil {
		return nil, err
	}

	transaction, err := client.Transaction.Get(context.Background(), input.OrderID, input.TransactionID, nil)
	if err != nil {
		return nil, err
	}

	if transaction == nil {
		return nil, fmt.Errorf("no transaction found with ID %d", input.TransactionID)
	}

	return transaction, nil
}

func NewGetTransactionAction() sdk.Action {
	return &GetTransactionAction{}
}

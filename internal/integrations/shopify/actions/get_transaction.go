package actions

import (
	"context"
	"errors"

	"github.com/wakflo/extensions/internal/integrations/shopify/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type getTransactionActionProps struct {
	OrderID       uint64 `json:"orderId"`
	TransactionID uint64 `json:"transactionId"`
}

type GetTransactionAction struct{}

func (a *GetTransactionAction) Name() string {
	return "Get Transaction"
}

func (a *GetTransactionAction) Description() string {
	return "Retrieves transaction details from a specified system or database, allowing you to access and utilize transactional data within your workflow."
}

func (a *GetTransactionAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *GetTransactionAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &getTransactionDocs,
	}
}

func (a *GetTransactionAction) Icon() *string {
	return nil
}

func (a *GetTransactionAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"orderId": autoform.NewNumberField().
			SetDisplayName("Order ID").
			SetDescription("The ID of the order.").
			SetRequired(true).
			Build(),
		"transactionId": autoform.NewNumberField().
			SetDisplayName("Transaction ID").
			SetDescription("The ID of the transaction.").
			SetRequired(true).
			Build(),
	}
}

func (a *GetTransactionAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[getTransactionActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	client, err := shared.CreateClient(ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	transaction, err := client.Transaction.Get(context.Background(), input.OrderID, input.TransactionID, nil)
	if err != nil {
		return nil, err
	}

	if transaction == nil {
		return nil, errors.New("no transaction found with ID ")
	}

	return sdk.JSON(map[string]interface{}{
		"transaction": transaction,
	}), nil
}

func (a *GetTransactionAction) Auth() *sdk.Auth {
	return nil
}

func (a *GetTransactionAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *GetTransactionAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewGetTransactionAction() sdk.Action {
	return &GetTransactionAction{}
}

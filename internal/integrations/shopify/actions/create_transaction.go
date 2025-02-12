package actions

import (
	"context"
	"fmt"

	goshopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/shopspring/decimal"
	"github.com/wakflo/extensions/internal/integrations/shopify/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type createTransactionActionProps struct {
	OrderID       uint64           `json:"orderId"`
	ParentID      *int64           `json:"parentId"`
	Kind          string           `json:"kind"`
	Currency      string           `json:"currency"`
	Amount        *decimal.Decimal `json:"amount"`
	Authorization string           `json:"authorization"`
	Source        string           `json:"source"`
}

type CreateTransactionAction struct{}

func (a *CreateTransactionAction) Name() string {
	return "Create Transaction"
}

func (a *CreateTransactionAction) Description() string {
	return "Create Transaction: Initiates a new transaction in your accounting or payment system, allowing you to automate the creation of financial records and streamline your business processes."
}

func (a *CreateTransactionAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *CreateTransactionAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &createTransactionDocs,
	}
}

func (a *CreateTransactionAction) Icon() *string {
	return nil
}

func (a *CreateTransactionAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"orderId": autoform.NewNumberField().
			SetDisplayName("Order ID").
			SetDescription("The ID of the order to create a transaction for.").
			SetRequired(true).
			Build(),
		"kind": autoform.NewSelectField().
			SetDisplayName("Type").
			SetOptions(shared.ShopifyTransactionKinds).
			SetRequired(true).
			Build(),
		"parentId": autoform.NewNumberField().
			SetDisplayName("Parent ID").
			SetDescription("The ID of an associated transaction.").
			SetRequired(false).
			Build(),
		"currency": autoform.NewShortTextField().
			SetDisplayName("Currency").
			SetRequired(false).
			Build(),
		"amount": autoform.NewNumberField().
			SetDisplayName("Amount").
			SetRequired(false).
			Build(),
		"authorization": autoform.NewShortTextField().
			SetDisplayName("Authorization Key.").
			SetRequired(false).
			Build(),
		"source": autoform.NewShortTextField().
			SetDisplayName("Source").
			SetDescription("An optional origin of the transaction. Set to external to import a cash transaction for the associated order.").
			SetRequired(false).
			Build(),
	}
}

func (a *CreateTransactionAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createTransactionActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	client, err := shared.CreateClient(ctx.BaseContext)
	if err != nil {
		return nil, err
	}
	newTransaction := goshopify.Transaction{
		ParentId:      input.ParentID,
		Kind:          input.Kind,
		OrderId:       input.OrderID,
		Amount:        input.Amount,
		Authorization: input.Authorization,
		Currency:      input.Currency,
		Source:        input.Source,
	}
	transaction, err := client.Transaction.Create(context.Background(), input.OrderID, newTransaction)
	if err != nil {
		return nil, err
	}

	if transaction == nil {
		return nil, fmt.Errorf("no order found with ID '%d'", input.OrderID)
	}

	return map[string]interface{}{
		"new_transaction": transaction,
	}, nil
}

func (a *CreateTransactionAction) Auth() *sdk.Auth {
	return nil
}

func (a *CreateTransactionAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *CreateTransactionAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewCreateTransactionAction() sdk.Action {
	return &CreateTransactionAction{}
}

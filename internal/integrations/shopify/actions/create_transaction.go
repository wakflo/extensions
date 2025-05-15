package actions

import (
	"context"
	"fmt"

	goshopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/juicycleff/smartform/v1"
	"github.com/shopspring/decimal"
	"github.com/wakflo/extensions/internal/integrations/shopify/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
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

func (a *CreateTransactionAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_transaction",
		DisplayName:   "Create Transaction",
		Description:   "Create Transaction: Initiates a new transaction in your accounting or payment system, allowing you to automate the creation of financial records and streamline your business processes.",
		Type:          core.ActionTypeAction,
		Documentation: createTransactionDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *CreateTransactionAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_transaction", "Create Transaction")

	form.NumberField("orderId", "Order ID").
		Placeholder("The ID of the order to create a transaction for.").
		Required(true).
		HelpText("The ID of the order to create a transaction for.")

	form.SelectField("kind", "Type").
		AddOption("sale", "Sale").
		AddOption("refund", "Refund").
		AddOption("void", "Void").
		AddOption("capture", "Capture").
		AddOption("authorization", "Authorization").
		Required(true).
		HelpText("Transaction type")

	form.NumberField("parentId", "Parent ID").
		Placeholder("The ID of an associated transaction.").
		HelpText("The ID of an associated transaction.")

	form.TextField("currency", "Currency").
		Placeholder("Currency").
		HelpText("Currency")

	form.NumberField("amount", "Amount").
		Placeholder("Amount").
		HelpText("Amount")

	form.TextField("authorization", "Authorization Key").
		Placeholder("Authorization Key.").
		HelpText("Authorization Key.")

	form.TextField("source", "Source").
		Placeholder("An optional origin of the transaction. Set to external to import a cash transaction for the associated order.").
		HelpText("An optional origin of the transaction. Set to external to import a cash transaction for the associated order.")

	schema := form.Build()
	return schema
}

func (a *CreateTransactionAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[createTransactionActionProps](ctx)
	if err != nil {
		return nil, err
	}

	client, err := shared.CreateClient(ctx)
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

func (a *CreateTransactionAction) Auth() *core.AuthMetadata {
	return nil
}

func NewCreateTransactionAction() sdk.Action {
	return &CreateTransactionAction{}
}

package actions

import (
	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/harvest/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type getInvoiceActionProps struct {
	InvoiceID string `json:"invoice-id"`
}

type GetInvoiceAction struct{}

// Metadata returns metadata about the action
func (a *GetInvoiceAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_invoice",
		DisplayName:   "Get Invoice",
		Description:   "Retrieves an invoice from the accounting system, allowing you to automate tasks that require access to invoice data.",
		Type:          core.ActionTypeAction,
		Documentation: getInvoiceDocs,
		SampleOutput: map[string]any{
			"id":             "13150453",
			"client_id":      "5735776",
			"number":         "1001",
			"purchase_order": "",
			"amount":         "88.23",
			"due_amount":     "88.23",
			"tax_amount":     "0",
			"tax2":           "0",
			"tax2_amount":    "0",
			"discount":       "0",
			"subject":        "Web Design",
			"notes":          "Thank you for your business!",
			"state":          "open",
			"issue_date":     "2023-05-01",
			"due_date":       "2023-05-31",
			"payment_term":   "net 30",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *GetInvoiceAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_invoice", "Get Invoice")

	form.TextField("invoice-id", "Invoice ID").
		Required(true).
		HelpText("the ID of the invoice.")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *GetInvoiceAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *GetInvoiceAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[getInvoiceActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	url := "/v2/invoices/" + input.InvoiceID

	invoice, err := shared.GetHarvestClient(authCtx.Token.AccessToken, url)
	if err != nil {
		return nil, err
	}
	return invoice, nil
}

func NewGetInvoiceAction() sdk.Action {
	return &GetInvoiceAction{}
}

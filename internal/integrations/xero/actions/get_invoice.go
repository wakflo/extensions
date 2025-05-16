package actions

import (
	"errors"
	"fmt"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/xero/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type getInvoiceActionProps struct {
	TenantID  string `json:"tenant_id"`
	InvoiceID string `json:"invoice_id"`
}

type GetInvoiceAction struct{}

func (a *GetInvoiceAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_invoice",
		DisplayName:   "Get Invoice",
		Description:   "Retrieves an invoice from the accounting system, allowing you to automate tasks that require access to invoice data.",
		Type:          sdkcore.ActionTypeAction,
		Documentation: getInvoiceDocs,
		SampleOutput: map[string]any{
			"message": "hello world",
		},
		Settings: sdkcore.ActionSettings{},
	}
}

func (a *GetInvoiceAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_invoice", "Get Invoice")

	shared.GetTenantProps("tenant_id", "Organization", "select organization", true, form)

	shared.GetInvoiceProp("invoice_id", "Invoices", "select invoice", false, form)

	schema := form.Build()

	return schema
}

func (a *GetInvoiceAction) Perform(ctx sdkcontext.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[getInvoiceActionProps](ctx)
	if err != nil {
		return nil, err
	}

	tokenSource := ctx.Auth().Token
	if tokenSource == nil {
		return nil, errors.New("missing authentication token")
	}
	token := tokenSource.AccessToken

	endpoint := "/Invoices/" + input.InvoiceID

	invoice, err := shared.GetXeroNewClient(token, endpoint, input.TenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch invoice: %v", err)
	}

	return invoice, nil
}

func (a *GetInvoiceAction) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewGetInvoiceAction() sdk.Action {
	return &GetInvoiceAction{}
}
